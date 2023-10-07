package zeus_core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	v1Apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceSumsTestSuite struct {
	K8TestSuite
}

func (s *ResourceSumsTestSuite) TestGetDiskRequirements() {
	diskSizeOne := "20Gi"
	sts := v1Apps.StatefulSet{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: v1Apps.StatefulSetSpec{
			Template: v1.PodTemplateSpec{},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{{
				TypeMeta:   metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{},
				Spec: v1.PersistentVolumeClaimSpec{
					Resources: v1.ResourceRequirements{
						Limits:   nil,
						Requests: v1.ResourceList{"storage": resource.MustParse(diskSizeOne)},
						Claims:   nil,
					},
				},
				Status: v1.PersistentVolumeClaimStatus{},
			}},
		},
		Status: v1Apps.StatefulSetStatus{},
	}
	rs := ResourceSums{}
	GetBlockStorageDiskRequirements(ctx, sts.Spec.VolumeClaimTemplates, &rs)
	s.Assert().NotEmpty(rs.DiskRequests)
	s.Assert().NotEmpty(rs.DiskLimits)
}

func (s *ResourceSumsTestSuite) TestGetResourceRequirements() {
	requestRAM := "12Gi"
	requestLimitRAM := "12Gi"

	requestCPU := "7"
	requestLimitCPU := "7"

	requestRAM2 := "1Gi"
	requestLimitRAM2 := "1Gi"

	requestCPU2 := "500m"
	requestLimitCPU2 := "500m"
	rr := v1.ResourceRequirements{
		Limits: v1.ResourceList{
			"cpu":    resource.MustParse(requestLimitCPU),
			"memory": resource.MustParse(requestLimitRAM),
		},
		Requests: v1.ResourceList{
			"cpu":    resource.MustParse(requestCPU),
			"memory": resource.MustParse(requestRAM),
		},
	}
	rr2 := v1.ResourceRequirements{
		Limits: v1.ResourceList{
			"cpu":    resource.MustParse(requestLimitCPU2),
			"memory": resource.MustParse(requestLimitRAM2),
		},
		Requests: v1.ResourceList{
			"cpu":    resource.MustParse(requestCPU2),
			"memory": resource.MustParse(requestRAM2),
		},
	}
	ps := v1.PodSpec{
		Containers: []v1.Container{{
			Resources: rr,
		}, {Resources: rr2}},
	}
	rs := ResourceSums{}

	GetResourceRequirements(ctx, ps, &rs)
	s.Assert().NotEmpty(rs)
	fmt.Println(rs)
}

func TestResourceSumsTestSuite(t *testing.T) {
	suite.Run(t, new(ResourceSumsTestSuite))
}
