package zeus_core

import (
	"context"

	aws_nvme "github.com/zeus-fyi/zeus/zeus/cluster_resources/nvme/aws"
	do_nvme "github.com/zeus-fyi/zeus/zeus/cluster_resources/nvme/do"
	gcp_nvme "github.com/zeus-fyi/zeus/zeus/cluster_resources/nvme/gcp"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type ResourceMinMax struct {
	Max ResourceAggregate `json:"max"`
	Min ResourceAggregate `json:"min"`
}

type ResourceAggregate struct {
	MemRequests  string `json:"memRequests"`
	CpuRequests  string `json:"cpuRequests"`
	DiskRequests string `json:"diskRequests"`
}

type ResourceSums struct {
	Replicas          string `json:"replicas"`
	MemRequests       string `json:"memRequests"`
	MemLimits         string `json:"memLimits"`
	CpuRequests       string `json:"cpuRequests"`
	CpuLimits         string `json:"cpuLimits"`
	DiskRequests      string `json:"diskRequests"`
	LocalDiskRequests string `json:"localDiskRequests"`
	DiskLimits        string `json:"diskLimits"`
}

func ApplyMinMaxConstraints(sums ResourceSums, reMinMax ResourceMinMax) (ResourceMinMax, error) {
	if sums.MemRequests != "" && reMinMax.Min.MemRequests == "" && sums.MemRequests != "0" {
		reMinMax.Min.MemRequests = sums.MemRequests
	}
	if sums.MemRequests != "" && reMinMax.Max.MemRequests == "" && sums.MemRequests != "0" {
		reMinMax.Max.MemRequests = sums.MemRequests
	}
	if sums.CpuRequests != "" && reMinMax.Min.CpuRequests == "" && sums.CpuRequests != "0" {
		reMinMax.Min.CpuRequests = sums.CpuRequests
	}
	if sums.CpuRequests != "" && reMinMax.Max.CpuRequests == "" && sums.CpuRequests != "0" {
		reMinMax.Max.CpuRequests = sums.CpuRequests
	}
	if sums.DiskRequests != "" && reMinMax.Min.DiskRequests == "" && sums.DiskRequests != "0" {
		reMinMax.Min.DiskRequests = sums.DiskRequests
	}
	if sums.DiskRequests != "" && reMinMax.Max.DiskRequests == "" && sums.DiskRequests != "0" {
		reMinMax.Max.DiskRequests = sums.DiskRequests
	}

	if sums.MemRequests != "" && sums.MemRequests != "0" {
		sumsMemQuantity, err := resource.ParseQuantity(sums.MemRequests)
		if err != nil {
			return reMinMax, err
		}
		minMemQuantity, err := resource.ParseQuantity(reMinMax.Min.MemRequests)
		if err != nil {
			return reMinMax, err
		}
		if sumsMemQuantity.Cmp(minMemQuantity) < 0 {
			reMinMax.Min.MemRequests = sums.MemRequests
		}
		maxMemQuantity, err := resource.ParseQuantity(reMinMax.Max.MemRequests)
		if err != nil {
			return reMinMax, err
		}
		if sumsMemQuantity.Cmp(maxMemQuantity) > 0 {
			reMinMax.Max.MemRequests = sums.MemRequests
		}
	}

	if sums.CpuRequests != "" && sums.CpuRequests != "0" {
		sumsCpuQuantity, err := resource.ParseQuantity(sums.CpuRequests)
		if err != nil {
			return reMinMax, err
		}
		minCpuQuantity, err := resource.ParseQuantity(reMinMax.Min.CpuRequests)
		if err != nil {
			return reMinMax, err
		}
		if sumsCpuQuantity.Cmp(minCpuQuantity) < 0 {
			reMinMax.Min.CpuRequests = sums.CpuRequests
		}
		maxCpuQuantity, err := resource.ParseQuantity(reMinMax.Max.CpuRequests)
		if err != nil {
			return reMinMax, err
		}
		if sumsCpuQuantity.Cmp(maxCpuQuantity) > 0 {
			reMinMax.Max.CpuRequests = sums.CpuRequests
		}
	}

	if sums.DiskRequests != "" && sums.DiskRequests != "0" {
		sumsDiskQuantity, err := resource.ParseQuantity(sums.DiskRequests)
		if err != nil {
			return reMinMax, err
		}
		minDiskQuantity, err := resource.ParseQuantity(reMinMax.Min.DiskRequests)
		if err != nil {
			return reMinMax, err
		}
		if sumsDiskQuantity.Cmp(minDiskQuantity) < 0 {
			reMinMax.Min.DiskRequests = sums.DiskRequests
		}
		maxDiskQuantity, err := resource.ParseQuantity(reMinMax.Max.DiskRequests)
		if err != nil {
			return reMinMax, err
		}
		if sumsDiskQuantity.Cmp(maxDiskQuantity) > 0 {
			reMinMax.Max.DiskRequests = sums.DiskRequests
		}
	}
	return reMinMax, nil
}

func GetResourceRequirements(ctx context.Context, spec v1.PodSpec, r *ResourceSums) {
	memRequests := resource.NewQuantity(0, resource.BinarySI)
	memLimits := resource.NewQuantity(0, resource.BinarySI)
	cpuRequests := resource.NewQuantity(0, resource.DecimalSI)
	cpuLimits := resource.NewQuantity(0, resource.DecimalSI)
	for _, c := range spec.Containers {
		memReq := c.Resources.Requests.Memory()
		memRequests.Add(*memReq)
		memLim := c.Resources.Limits.Memory()
		memLimits.Add(*memLim)
		cpuReq := c.Resources.Requests.Cpu()
		cpuRequests.Add(*cpuReq)
		cpuLim := c.Resources.Limits.Cpu()
		cpuLimits.Add(*cpuLim)
	}
	r.MemRequests = memRequests.String()
	r.MemLimits = memLimits.String()
	r.CpuRequests = cpuRequests.String()
	r.CpuLimits = cpuLimits.String()
}

func GetBlockStorageDiskRequirements(ctx context.Context, pvcs []v1.PersistentVolumeClaim, r *ResourceSums) {
	diskRequests := resource.NewQuantity(0, resource.BinarySI)
	diskLimits := resource.NewQuantity(0, resource.BinarySI)
	localDiskRequests := resource.NewQuantity(0, resource.BinarySI)

	for _, pvc := range pvcs {
		if pvc.Spec.StorageClassName != nil {
			scName := *pvc.Spec.StorageClassName
			switch *pvc.Spec.StorageClassName {
			case aws_nvme.AwsStorageClass, gcp_nvme.GcpStorageClass:
				sr := pvc.Spec.Resources.Requests.Storage()
				if sr != nil {
					localDiskRequests.Add(*sr)
				}
				continue
			default:
				if scName == do_nvme.DoStorageClass {
					sr := pvc.Spec.Resources.Requests.Storage()
					if sr != nil {
						localDiskRequests.Add(*sr)
					}
					continue
				}
			}
		}
		dr := pvc.Spec.Resources.Requests.Storage()
		dl := pvc.Spec.Resources.Limits.Storage()
		if dr != nil {
			diskRequests.Add(*dr)
		}
		if dl != nil {
			diskLimits.Add(*dl)
		}
	}
	r.DiskRequests = diskRequests.String()
	r.DiskLimits = diskLimits.String()
}
