package chart_workload

import (
	"encoding/json"
)

func (nk *TopologyBaseInfraWorkload) MarshallStatefulSet() ([]byte, error) {
	var b []byte
	if nk.StatefulSet != nil {
		bsts, err := json.Marshal(nk.StatefulSet)
		return bsts, err
	}
	return b, nil
}

func (nk *TopologyBaseInfraWorkload) MarshallDeployment() ([]byte, error) {
	var b []byte
	if nk.Deployment != nil {
		bsts, err := json.Marshal(nk.Deployment)
		return bsts, err
	}
	return b, nil
}

func (nk *TopologyBaseInfraWorkload) MarshallService() ([]byte, error) {
	var b []byte
	if nk.Service != nil {
		bsts, err := json.Marshal(nk.StatefulSet)
		return bsts, err
	}
	return b, nil
}

func (nk *TopologyBaseInfraWorkload) MarshallIngress() ([]byte, error) {
	var b []byte
	if nk.Ingress != nil {
		bsts, err := json.Marshal(nk.Ingress)
		return bsts, err
	}
	return b, nil
}

func (nk *TopologyBaseInfraWorkload) MarshallConfigMap() ([]byte, error) {
	var b []byte
	if nk.ConfigMap != nil {
		bsts, err := json.Marshal(nk.ConfigMap)
		return bsts, err
	}
	return b, nil
}
