package zeus_core

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceTestSuite struct {
	K8TestSuite
}

func (s *NamespaceTestSuite) TestGetK8Namespace() {
	ctx := context.Background()
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "ephemeral"
	nsl, err := s.K.GetNamespace(ctx, kns)
	s.Nil(err)
	s.Require().NotEmpty(nsl)
}

func (s *NamespaceTestSuite) TestCreateNamespaceIfDoesNotExist() {
	ctx := context.Background()
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "demo"
	nsl, err := s.K.CreateNamespaceIfDoesNotExist(ctx, kns)
	s.Nil(err)
	s.Require().NotEmpty(nsl)
}

func (s *NamespaceTestSuite) TestListK8Namespaces() {
	ctx := context.Background()
	s.K.SetContext("do-sfo3-dev-do-sfo3-zeus")
	nsl, err := s.K.GetNamespaces(ctx)
	s.Nil(err)
	s.Greater(len(nsl.Items), 0)
	for _, n := range nsl.Items {
		fmt.Println(n.Name)
	}

	fmt.Println("=========== new context ===========")
	s.K.SetContext("do-nyc1-do-nyc1-zeus-demo")
	nsl, err = s.K.GetNamespaces(ctx)
	s.Nil(err)
	s.Greater(len(nsl.Items), 0)
	for _, n := range nsl.Items {
		fmt.Println(n.Name)
	}
}

func (s *NamespaceTestSuite) TestCreateK8sNamespace() {
	ctx := context.Background()
	ns := v1.Namespace{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1.NamespaceSpec{},
		Status:     v1.NamespaceStatus{},
	}
	ns.Name = "demo"
	newNamespace, err := s.K.CreateNamespace(ctx, &ns)
	s.Require().Nil(err)
	s.NotEmpty(newNamespace)
}

func (s *NamespaceTestSuite) TestDeleteNamespace() {
	ctx := context.Background()
	var kns zeus_common_types.CloudCtxNs
	kns.Namespace = "demo"
	err := s.K.DeleteNamespace(ctx, kns)
	s.Require().Nil(err)
}

func TestNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(NamespaceTestSuite))
}
