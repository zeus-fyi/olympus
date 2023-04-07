package zeus_core

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type ResourceSums struct {
	Replicas     string `json:"replicas"`
	MemRequests  string `json:"memRequests"`
	MemLimits    string `json:"memLimits"`
	CpuRequests  string `json:"cpuRequests"`
	CpuLimits    string `json:"cpuLimits"`
	DiskRequests string `json:"diskRequests"`
	DiskLimits   string `json:"diskLimits"`
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

func GetDiskRequirements(ctx context.Context, pvcs []v1.PersistentVolumeClaim, r *ResourceSums) {
	diskRequests := resource.NewQuantity(0, resource.BinarySI)
	diskLimits := resource.NewQuantity(0, resource.BinarySI)

	for _, pvc := range pvcs {
		dr := pvc.Spec.Resources.Requests.Storage()
		diskRequests.Add(*dr)
		dl := pvc.Spec.Resources.Limits.Storage()
		diskLimits.Add(*dl)
	}
	r.DiskRequests = diskRequests.String()
	r.DiskLimits = diskLimits.String()
}
