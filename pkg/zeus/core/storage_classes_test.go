package zeus_core

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

type StorageClassTestSuite struct {
	K8TestSuite
}

func (s *StatefulSetsTestSuite) TestStorageClass() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "aws", Region: "us-east-1", Context: "zeus-eks-us-east-1"}
	sc := &v1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "aws-ebs-gp3-max-performance", // Provide a meaningful name for the StorageClass
		},
		Provisioner: "ebs.csi.aws.com", // AWS EBS CSI driver
		Parameters: map[string]string{
			"type":       "gp3",   // Specify gp3 type for the EBS volume
			"iops":       "16000", // Maximum IOPS for gp3
			"throughput": "1000",  // Maximum throughput in MB/s for gp3
			//"encrypted":  "true",  // Optionally, ensure encryption is enabled
			// "fsType":      "ext4",            // Specify filesystem type if needed, e.g., ext4 or xfs
		},
		ReclaimPolicy:        nil,                   // You can specify a ReclaimPolicy if needed
		AllowVolumeExpansion: pointer.BoolPtr(true), // Optionally allow volume expansion
	}
	_, err := s.K.CreateStorageClass(ctx, kns, sc)
	s.Require().Nil(err)
}

func TestStorageClass(t *testing.T) {
	suite.Run(t, new(StatefulSetsTestSuite))
}
