package hestia_gcp

import (
	"context"
	"fmt"
	"net/url"

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
}

func InitGcpClient(ctx context.Context, authJsonBytes []byte) (GcpClient, error) {
	client, err := container.NewService(ctx, option.WithCredentialsJSON(authJsonBytes), option.WithScopes(container.CloudPlatformScope))
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to create GKE API client")
		return GcpClient{}, err
	}
	return GcpClient{client}, nil
}

func (g *GcpClient) ListMachineTypes(ctx context.Context, ci GcpClusterInfo, authJsonBytes []byte) (MachineTypes, error) {
	mt := MachineTypes{}
	jwtConfig, err := google.JWTConfigFromJSON(authJsonBytes, container.CloudPlatformScope, ComputeScope, ComputeReadOnlyScope)
	if err != nil {
		fmt.Printf("Error creating JWT config: %v\n", err)
		return mt, err
	}
	httpClient := jwtConfig.Client(ctx)
	restyClient := resty.NewWithClient(httpClient)
	project := ci.ProjectID
	zone := ci.Zone
	maxResults := 500
	orderBy := "creationTimestamp desc"
	returnPartialSuccess := false
	queryParams := url.Values{}
	queryParams.Set("maxResults", fmt.Sprintf("%d", maxResults))
	queryParams.Set("orderBy", orderBy)
	queryParams.Set("returnPartialSuccess", fmt.Sprintf("%t", returnPartialSuccess))

	// GET /compute/v1/projects/{project}/zones/{zone}/machineTypes
	requestURL := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/machineTypes", project, zone)
	requestURL = fmt.Sprintf("%s?%s", requestURL, queryParams.Encode())
	// Execute the request
	resp, err := restyClient.R().SetResult(&mt).Get(requestURL)
	if err != nil {
		fmt.Printf("Error executing request: %v\n", err)
		return mt, err
	}
	// Check for non-2xx status codes
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		fmt.Printf("Error: API responded with status code %d\n", resp.StatusCode())
		return mt, err
	}
	fmt.Println(resp.String())
	return mt, err
}

type GkeNodePoolInfo struct {
	NodePoolID       string `json:"nodePoolID,omitempty"`
	Name             string `json:"name"`
	MachineType      string `json:"machineType"`
	InitialNodeCount int64  `json:"initialNodeCount"`
}

func (g *GcpClient) RemoveNodePool(ctx context.Context, ci GcpClusterInfo, ni GkeNodePoolInfo) (any, error) {
	resp, err := g.Projects.Zones.Clusters.NodePools.Delete(ci.ProjectID, ci.Zone, ci.ClusterName, ni.NodePoolID).Context(ctx).Do()
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to delete node pool")
		return nil, err
	}
	return resp, err
}

func (g *GcpClient) AddNodePool(ctx context.Context, ci GcpClusterInfo, ni GkeNodePoolInfo) (any, error) {

	// TODO: add taints
	//t := container.NodeTaint{
	//	Effect:          "",
	//	Key:             "",
	//	Value:           "",
	//	ForceSendFields: nil,
	//	NullFields:      nil,
	//}

	cnReq := &container.CreateNodePoolRequest{
		ClusterId: ci.ClusterName,
		NodePool: &container.NodePool{
			Name:             "node-pool-1",
			InitialNodeCount: ni.InitialNodeCount,
			Autoscaling: &container.NodePoolAutoscaling{
				Autoprovisioned: false,
				Enabled:         false,
			},
			Config: &container.NodeConfig{
				Labels:         nil,
				MachineType:    ni.MachineType,
				Metadata:       nil,
				NodeGroup:      "",
				ResourceLabels: nil,
				ServiceAccount: "",
				Spot:           false,
				Tags:           nil,
				Taints:         nil,
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
