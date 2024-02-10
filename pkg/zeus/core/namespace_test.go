package zeus_core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceTestSuite struct {
	K8TestSuite
}

func (s *NamespaceTestSuite) TestGetK8Namespace() {
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "ephemeral"
	nsl, err := s.K.GetNamespace(ctx, kns)
	s.Nil(err)
	s.Require().NotEmpty(nsl)
}

func (s *NamespaceTestSuite) TestCreateNamespaceIfDoesNotExist() {
	s.K.SetContext("gke_zeusfyi_us-central1-a_zeus-gcp-pilot-0")
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "demo"
	nsl, err := s.K.CreateNamespaceIfDoesNotExist(ctx, kns)
	s.Nil(err)
	s.Require().NotEmpty(nsl)
}

func (s *NamespaceTestSuite) TestListK8Namespaces() {
	s.K.SetContext("do-sfo3-dev-do-sfo3-zeus")
	nsl, err := s.K.GetNamespaces(ctx, zeus_common_types.CloudCtxNs{
		CloudProvider: "",
		Region:        "",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "",
		Env:           "",
	})
	s.Nil(err)
	s.Greater(len(nsl.Items), 0)
	for _, n := range nsl.Items {
		fmt.Println(n.Name)
	}

	fmt.Println("=========== new context ===========")
	s.K.SetContext("do-nyc1-do-nyc1-zeus-demo")
	nsl, err = s.K.GetNamespaces(ctx, zeus_common_types.CloudCtxNs{
		CloudProvider: "",
		Region:        "",
		Context:       "do-nyc1-do-nyc1-zeus-demo",
		Namespace:     "",
		Env:           "",
	})
	s.Nil(err)
	s.Greater(len(nsl.Items), 0)
	for _, n := range nsl.Items {
		fmt.Println(n.Name)
	}

	// gke from inmemfs uses a different src for gcloud binary
	//fmt.Println("=========== new context ===========")
	//s.K.SetContext("gke_zeusfyi_us-central1-a_zeus-gcp-pilot-0")
	//nsl, err = s.K.GetNamespaces(ctx, zeus_common_types.CloudCtxNs{
	//	CloudProvider: "",
	//	Region:        "",
	//	Context:       "gke_zeusfyi_us-central1-a_zeus-gcp-pilot-0",
	//	Namespace:     "",
	//	Env:           "",
	//})
	//s.Nil(err)
	//s.Greater(len(nsl.Items), 0)
	//for _, n := range nsl.Items {
	//	fmt.Println(n.Name)
	//}

	fmt.Println("=========== new context ===========")
	s.K.SetContext("zeus-us-west-1")
	nsl, err = s.K.GetNamespaces(ctx, zeus_common_types.CloudCtxNs{
		CloudProvider: "",
		Region:        "",
		Context:       "arn:aws:eks:us-west-1:480391564655:cluster/zeus-us-west-1",
		Namespace:     "",
		Env:           "",
	})
	s.Nil(err)
	s.Greater(len(nsl.Items), 0)
	for _, n := range nsl.Items {
		fmt.Println(n.Name)
	}
}

func (s *NamespaceTestSuite) TestCreateK8sNamespace() {
	ns := v1.Namespace{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1.NamespaceSpec{},
		Status:     v1.NamespaceStatus{},
	}
	ns.Name = "demo"
	newNamespace, err := s.K.CreateNamespace(ctx, zeus_common_types.CloudCtxNs{}, &ns)
	s.Require().Nil(err)
	s.NotEmpty(newNamespace)
}

func (s *NamespaceTestSuite) TestDeleteNamespace() {
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "demo"
	err := s.K.DeleteNamespace(ctx, kns)
	s.Require().Nil(err)
}

func TestNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(NamespaceTestSuite))
}
