package hestia_ovhcloud

// zeusfyi-shared
// clusterID: a7ea8ded-fa8f-48f3-83d7-ce01410552bc
// 2vdfsl.c1.or1.k8s.ovh.us

// 	Flavor == instance type
// service name: zeusfyi
// clusterID: 750cf38b-0965-4b2b-b6ba-9728ca3f239e
// nodePoolID: 0f696c39-46cb-442c-bdc7-981887a5f08c
// nodePoolName: nodepool-7452815c-94ca-40d8-8334-cda7bbcc8f73
// unx0wx.c1.or1.k8s.ovh.us

func (s *OvhCloudTestSuite) TestGetNodePool() {
	resp, err := s.o.GetNodePool(ctx, OvhNodePoolCreationRequest{
		ServiceName: OvhServiceName,
		KubeId:      OvhInternalKubeID,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
}

func (s *OvhCloudTestSuite) TestCreateNodePool() {
	poolName := "testpool2"
	flavorName := "c2-7"
	autoscaleEnabled := false
	req := OvhNodePoolCreationRequest{
		ServiceName: OvhServiceName,
		KubeId:      OvhInternalKubeID,
		ProjectKubeNodePoolCreation: ProjectKubeNodePoolCreation{
			AntiAffinity:  nil,
			Autoscale:     &autoscaleEnabled,
			Autoscaling:   nil,
			DesiredNodes:  1,
			FlavorName:    flavorName,
			MaxNodes:      1,
			MinNodes:      1,
			MonthlyBilled: nil,
			Name:          poolName,
			Template:      nil,
		},
	}

	resp, err := s.o.CreateNodePool(ctx, req)
	s.Require().NoError(err)
	s.Assert().NotNil(resp)
}

func (s *OvhCloudTestSuite) TestRemoveNodePool() {
	poolID := "82bf43bf-1809-461e-8ac1-591c2f45e7a8"
	err := s.o.RemoveNodePool(ctx, OvhInternalKubeID, poolID)
	s.Require().Nil(err)
}
