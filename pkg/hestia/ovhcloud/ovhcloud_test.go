package hestia_ovhcloud

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type OvhCloudTestSuite struct {
	test_suites_base.TestSuite
	o OvhCloud
}

func (s *OvhCloudTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	//apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	creds := OvhCloudCreds{
		Region:      OvhUS,
		AppKey:      s.Tc.OvhAppKey,
		AppSecret:   s.Tc.OvhSecretKey,
		ConsumerKey: s.Tc.OvhConsumerKey,
	}
	s.o = InitOvhClient(ctx, creds)
}

// zeusfyi-shared
// clusterID: a7ea8ded-fa8f-48f3-83d7-ce01410552bc
// 2vdfsl.c1.or1.k8s.ovh.us

// 	Flavor == instance type
// service name: zeusfyi
// clusterID: 750cf38b-0965-4b2b-b6ba-9728ca3f239e
// nodePoolID: 0f696c39-46cb-442c-bdc7-981887a5f08c
// nodePoolName: nodepool-7452815c-94ca-40d8-8334-cda7bbcc8f73
// unx0wx.c1.or1.k8s.ovh.us

func (s *OvhCloudTestSuite) TestListSizes() {

}

func (s *OvhCloudTestSuite) TestCreateNodePool() {
	poolName := "testpool"
	flavorName := ""
	req := OvhNodePoolCreationRequest{
		ServiceName: "zeusfyi",
		KubeId:      "750cf38b-0965-4b2b-b6ba-9728ca3f239e",
		ProjectKubeNodePoolCreation: ProjectKubeNodePoolCreation{
			AntiAffinity:  false,
			Autoscale:     false,
			DesiredNodes:  1,
			FlavorName:    flavorName,
			MaxNodes:      1,
			MinNodes:      1,
			MonthlyBilled: false,
			Name:          poolName,
			Template: struct {
				Metadata struct {
					Annotations struct{} `json:"annotations"`
					Finalizers  []string `json:"finalizers"`
					Labels      struct{} `json:"labels"`
				} `json:"metadata"`
				Spec struct {
					Taints []struct {
						Effect string `json:"effect"`
						Key    string `json:"key"`
						Value  string `json:"value"`
					} `json:"taints"`
					Unschedulable bool `json:"unschedulable"`
				} `json:"spec"`
			}{},
		},
	}
	resp, err := s.o.CreateNodePool(ctx, req)
	s.Require().NoError(err)
	s.Assert().NotNil(resp)
}

func (s *OvhCloudTestSuite) TestRemoveNodePool() {
	err := s.o.RemoveNodePool(ctx, "", "")
	s.Require().Nil(err)
}

func TestOvhCloudTestSuite(t *testing.T) {
	suite.Run(t, new(OvhCloudTestSuite))
}
