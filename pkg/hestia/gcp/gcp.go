package hestia_gcp

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

const (
	ProjectID            = "zeusfyi"
	ComputeScope         = "https://www.googleapis.com/auth/compute"
	ComputeReadOnlyScope = "https://www.googleapis.com/auth/compute.readonly"
)

/*
https://www.googleapis.com/auth/cloud-billing
https://www.googleapis.com/auth/cloud-billing.readonly
*/

/*
	General-purpose—best price-performance ratio for a variety of workloads.
	Compute-optimized—highest performance per core on Compute Engine and optimized for compute-intensive workloads.
	Memory-optimized—ideal for memory-intensive workloads, offering more memory per core than other machine families, with up to 12 TB of memory.
	Accelerator-optimized—ideal for massively parallelized Compute Unified Device Architecture (CUDA) compute workloads, such as machine learning (ML) and high performance computing (HPC).
	This family is the best option for workloads that require GPUs.
*/

type GcpClusterInfo struct {
	ClusterName string
	ProjectID   string
	Zone        string
}

type GcpClient struct {
	*container.Service
	*resty.Client
}

func NonSupported(name string) float64 {
	if strings.HasPrefix(name, "f1-micro") {
		return 0
	}
	if strings.HasPrefix(name, "g1-small") {
		return 0
	}
	if strings.HasPrefix(name, "n1-ultramem") {
		return 0
	}
	if strings.HasPrefix(name, "n1-megamem") {
		return 0
	}
	if strings.HasPrefix(name, "n1-highcpu-96") {
		return 0
	}
	if strings.HasPrefix(name, "n1-standard-96") {
		return 0
	}
	if strings.HasPrefix(name, "m2") {
		return 0
	}
	if strings.HasPrefix(name, "m1") {
		return 0
	}
	if strings.Contains(name, "node") {
		return 0
	}
	if strings.Contains(name, "a2-ultragpu") {
		return 0
	}
	if strings.Contains(name, "a2-highgpu") {
		return 0
	}
	if strings.Contains(name, "a2-megagpu") {
		return 0
	}
	return 1
}

func InitGcpClient(ctx context.Context, authJsonBytes []byte) (GcpClient, error) {
	client, err := container.NewService(ctx, option.WithCredentialsJSON(authJsonBytes), option.WithScopes(container.CloudPlatformScope))
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to create GKE API client")
		return GcpClient{}, err
	}
	jwtConfig, jerr := google.JWTConfigFromJSON(authJsonBytes, container.CloudPlatformScope, ComputeScope, ComputeReadOnlyScope)
	if jerr != nil {
		log.Ctx(ctx).Err(jerr).Msgf("Error creating JWT config: %v\n", jerr)
		return GcpClient{}, jerr
	}
	httpClient := jwtConfig.Client(ctx)
	restyClient := resty.NewWithClient(httpClient)
	return GcpClient{client, restyClient}, nil
}

func (g *GcpClient) ListMachineTypes(ctx context.Context, ci GcpClusterInfo, authJsonBytes []byte) (MachineTypes, error) {
	mt := MachineTypes{}
	project := ci.ProjectID
	zone := ci.Zone
	maxResults := 500
	orderBy := "creationTimestamp desc"
	queryParams := url.Values{}
	queryParams.Set("maxResults", fmt.Sprintf("%d", maxResults))
	queryParams.Set("orderBy", orderBy)
	queryParams.Set("pageToken", mt.NextPageToken)
	for {
		// GET /compute/v1/projects/{project}/zones/{zone}/machineTypes
		requestURL := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/machineTypes", project, zone)
		requestURL = fmt.Sprintf("%s?%s", requestURL, queryParams.Encode())
		// Execute the request

		tmp := MachineTypes{}
		resp, err := g.R().SetResult(&tmp).Get(requestURL)
		if err != nil {
			fmt.Printf("Error executing request: %v\n", err)
			return mt, err
		}
		mt.Items = append(mt.Items, tmp.Items...)
		// Check for non-2xx status codes
		if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
			fmt.Printf("Error: API responded with status code %d\n", resp.StatusCode())
			return mt, err
		}
		if tmp.NextPageToken == "" {
			return mt, nil
		}
		queryParams.Set("pageToken", tmp.NextPageToken)
	}
}

type GkeNodePoolInfo struct {
	Name             string `json:"name"`
	MachineType      string `json:"machineType"`
	InitialNodeCount int64  `json:"initialNodeCount"`
	NvmeDisks        int64  `json:"nvmeDisks,omitempty"`
}

func (g *GcpClient) RemoveNodePool(ctx context.Context, ci GcpClusterInfo, ni GkeNodePoolInfo) (any, error) {
	resp, err := g.Projects.Zones.Clusters.NodePools.Delete(ci.ProjectID, ci.Zone, ci.ClusterName, ni.Name).Context(ctx).Do()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to delete node pool")
		return nil, err
	}
	return resp, err
}

func (g *GcpClient) AddNodePool(ctx context.Context, ci GcpClusterInfo, ni GkeNodePoolInfo, taints []*container.NodeTaint, labels map[string]string) (*container.Operation, error) {
	if ni.MachineType == "n2-highmem-16" {
		ni.NvmeDisks = 16
		log.Info().Msg("n2-highmem-16 has 16 nvme disks")
	}
	cnReq := &container.CreateNodePoolRequest{
		ClusterId: ci.ClusterName,
		NodePool: &container.NodePool{
			Name:             ni.Name,
			InitialNodeCount: ni.InitialNodeCount,
			Autoscaling: &container.NodePoolAutoscaling{
				Autoprovisioned: false,
				Enabled:         false,
			},
			Config: &container.NodeConfig{
				Labels:          labels,
				LinuxNodeConfig: nil,
				LocalNvmeSsdBlockConfig: &container.LocalNvmeSsdBlockConfig{
					LocalSsdCount: ni.NvmeDisks,
				},
				MachineType: ni.MachineType,
				Taints:      taints,
			},
		},
		ProjectId: ci.ProjectID,
		Zone:      ci.Zone,
	}
	resp, err := g.Projects.Zones.Clusters.NodePools.Create(ci.ProjectID, ci.Zone, ci.ClusterName, cnReq).Context(ctx).Do()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to create node pool")
		return nil, err
	}
	return resp, err
}

func (g *GcpClient) ListNodes(ctx context.Context, ci GcpClusterInfo) ([]*container.NodePool, error) {
	nodePools, err := g.Projects.Zones.Clusters.NodePools.List(ci.ProjectID, ci.Zone, ci.ClusterName).Context(ctx).Do()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to retrieve node pools")
		return nil, err
	}
	return nodePools.NodePools, err
}

type MachineTypes struct {
	Kind          string               `json:"kind"`
	ID            string               `json:"id"`
	Items         []ComputeEngineItem  `json:"items"`
	NextPageToken string               `json:"nextPageToken"`
	SelfLink      string               `json:"selfLink"`
	Warning       ComputeEngineWarning `json:"warning"`
}

type ComputeEngineItem struct {
	Kind                         string                     `json:"kind"`
	ID                           string                     `json:"id"`
	CreationTimestamp            string                     `json:"creationTimestamp"`
	Name                         string                     `json:"name"`
	Description                  string                     `json:"description"`
	GuestCpus                    int                        `json:"guestCpus"`
	MemoryMb                     int                        `json:"memoryMb"`
	ImageSpaceGb                 int                        `json:"imageSpaceGb"`
	ScratchDisks                 []ComputeEngineDisk        `json:"scratchDisks"`
	MaximumPersistentDisks       int                        `json:"maximumPersistentDisks"`
	MaximumPersistentDisksSizeGb string                     `json:"maximumPersistentDisksSizeGb"`
	Deprecated                   ComputeEngineDeprecated    `json:"deprecated"`
	Zone                         string                     `json:"zone"`
	SelfLink                     string                     `json:"selfLink"`
	IsSharedCpu                  bool                       `json:"isSharedCpu"`
	Accelerators                 []ComputeEngineAccelerator `json:"accelerators"`
}

type SkuPricesLookup struct {
	GPUType    string  `json:"gpuType"`
	GPUs       int     `json:"gpus"`
	CPUs       float64 `json:"cpus"`
	MemGB      float64 `json:"gb"`
	Name       string  `json:"name"`
	DiskSizeGB int     `json:"diskSizeGB"`
}

func (c *ComputeEngineItem) GetSkuLookup() SkuPricesLookup {
	name, err := c.GetSkuInstanceNamePrefix()
	if err != nil {
		return SkuPricesLookup{}
	}
	gpuType, gpuCount := c.CountGPUs()
	skuInfo := SkuPricesLookup{
		GPUType:    gpuType,
		GPUs:       gpuCount,
		CPUs:       c.CountCPUs(),
		MemGB:      c.CountGB(),
		Name:       name,
		DiskSizeGB: c.GetDiskSizeGB(),
	}
	return skuInfo
}

func (c *ComputeEngineItem) CountGPUs() (string, int) {
	gpus := 0
	if len(c.Accelerators) > 1 {
		panic("more than one accelerator, not supported")
	}
	name := ""
	for _, ac := range c.Accelerators {
		gpus += ac.GuestAcceleratorCount
		parts := strings.Split(ac.GuestAcceleratorType, "-")
		if len(parts) < 1 || parts[0] == "" {
			err := errors.New("invalid input")
			panic(err)
		}
		name = strings.ToUpper(parts[len(parts)-1]) + " GPU"
	}
	if gpus == 0 {
		name = "none"
	}
	return name, gpus
}

func (c *ComputeEngineItem) CountCPUs() float64 {
	vCPUs := float64(0)
	vCPUs += float64(c.GuestCpus)
	if c.Name == "e2-medium" {
		vCPUs /= 2
	}
	if c.Name == "e2-small" {
		vCPUs /= 4
	}
	if c.Name == "e2-micro" {
		vCPUs /= 8
	}
	return vCPUs
}
func (c *ComputeEngineItem) GetDiskSizeGB() int {
	size, err := strconv.Atoi(c.MaximumPersistentDisksSizeGb)
	if err != nil {
		panic(err)
	}
	return size
}

func (c *ComputeEngineItem) CountGB() float64 {
	return float64(c.MemoryMb) / 1024
}

func (c *ComputeEngineItem) GetSkuInstanceNamePrefix() (string, error) {
	parts := strings.Split(c.Name, "-")
	if len(parts) < 1 || parts[0] == "" {
		err := errors.New("invalid input")
		return "", err
	}
	prefixName := strings.ToUpper(parts[0])
	if prefixName == "M3" {
		prefixName = "M3 Memory-optimized"
	}
	if prefixName == "T2A" {
		prefixName = "T2A Arm"
	}
	if prefixName == "T2D" {
		prefixName = "T2D AMD"
	}
	if prefixName == "C2D" {
		prefixName = "C2D AMD"
	}
	if prefixName == "N2D" {
		prefixName = "N2D AMD"
	}
	if prefixName == "N1" {
		prefixName = "N1 Predefined"
	}
	firstPart := prefixName + " Instance"
	if prefixName == "C2" {
		firstPart = "Compute optimized"
	}
	return firstPart, nil
}

type ComputeEngineDisk struct {
	DiskGb int `json:"diskGb"`
}

type ComputeEngineDeprecated struct {
	State       string `json:"state"`
	Replacement string `json:"replacement"`
	Deprecated  string `json:"deprecated"`
	Obsolete    string `json:"obsolete"`
	Deleted     string `json:"deleted"`
}

type ComputeEngineAccelerator struct {
	GuestAcceleratorType  string `json:"guestAcceleratorType"`
	GuestAcceleratorCount int    `json:"guestAcceleratorCount"`
}

type ComputeEngineWarning struct {
	Code    string                     `json:"code"`
	Message string                     `json:"message"`
	Data    []ComputeEngineWarningData `json:"data"`
}

type ComputeEngineWarningData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
