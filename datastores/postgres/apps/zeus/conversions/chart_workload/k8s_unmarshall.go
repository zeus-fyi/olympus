package chart_workload

import (
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/apps/v1"
	v1core "k8s.io/api/core/v1"
	v1networking "k8s.io/api/networking/v1"
)

func (nk *NativeK8s) DecodeBytes(jsonBytes []byte) error {
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
	case "ConfigMap":
		nk.ConfigMap = &v1core.ConfigMap{}
		err = json.Unmarshal(jsonBytes, nk.ConfigMap)
	case "Service":
		nk.Service = &v1core.Service{}
		err = json.Unmarshal(jsonBytes, nk.Service)
	case "Ingress":
		nk.Ingress = &v1networking.Ingress{}
		err = json.Unmarshal(jsonBytes, nk.Ingress)
	default:
		err = errors.New("NativeK8s: DecodeBytes, no matching kind found")
		log.Err(err).Msg("NativeK8s: DecodeBytes, no matching kind found")
	}
	return err
}
