package chart_workload

import (
	"encoding/json"
	"errors"

	v1sm "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/apps/v1"
	v1Batch "k8s.io/api/batch/v1"
	v1core "k8s.io/api/core/v1"
	v1networking "k8s.io/api/networking/v1"
)

func (nk *TopologyBaseInfraWorkload) DecodeBytes(jsonBytes []byte) error {
	metaType, err := nk.IdWorkloadFromBytes(jsonBytes)
	if err != nil {
		return err
	}
	switch metaType.Kind {
	case "Deployment":
		nk.Deployment = &v1.Deployment{}
		err = json.Unmarshal(jsonBytes, nk.Deployment)
	case "StatefulSet":
		nk.StatefulSet = &v1.StatefulSet{}
		err = json.Unmarshal(jsonBytes, nk.StatefulSet)
	case "Job":
		nk.Job = &v1Batch.Job{}
		err = json.Unmarshal(jsonBytes, nk.Job)
	case "CronJob":
		nk.CronJob = &v1Batch.CronJob{}
		err = json.Unmarshal(jsonBytes, nk.CronJob)
	case "ConfigMap":
		nk.ConfigMap = &v1core.ConfigMap{}
		err = json.Unmarshal(jsonBytes, nk.ConfigMap)
	case "Service":
		nk.Service = &v1core.Service{}
		err = json.Unmarshal(jsonBytes, nk.Service)
	case "Ingress":
		nk.Ingress = &v1networking.Ingress{}
		err = json.Unmarshal(jsonBytes, nk.Ingress)
	case "ServiceMonitor":
		nk.ServiceMonitor = &v1sm.ServiceMonitor{}
		err = json.Unmarshal(jsonBytes, nk.ServiceMonitor)
	default:
		err = errors.New("TopologyBaseInfraWorkload: DecodeBytes, no matching kind found")
		log.Err(err).Msg("TopologyBaseInfraWorkload: DecodeBytes, no matching kind found")
	}
	return err
}
