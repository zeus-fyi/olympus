package statefulset

import "github.com/zeus-fyi/olympus/pkg/utils/chronos"

func (s *StatefulSet) SetSpecParentIDs() {
	ts := chronos.Chronos{}
	parentID := ts.UnixTimeStampNow()
	s.Spec.SpecWorkload.SetParentClassTypeIDs(parentID)
	s.Spec.StatefulSetUpdateStrategy.ChartSubcomponentParentClassTypeID = parentID
	s.Spec.PodManagementPolicy.ChartSubcomponentParentClassTypeID = parentID
	s.Spec.ServiceName.ChartSubcomponentParentClassTypeID = parentID
}
