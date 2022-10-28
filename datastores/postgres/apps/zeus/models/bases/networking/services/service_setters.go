package services

func (s *Service) SetChartPackageID(id int) {
	s.ServiceSpec.ChartPackageID = id
	s.Metadata.ChartPackageID = id
}
