package db_to_k8s_conversions

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertMetadata(k8sMetadata *metav1.ObjectMeta, pcSlice []common_conversions.PC) {
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		switch subClassName {
		case "labels":
			if k8sMetadata.Labels == nil {
				k8sMetadata.Labels = make(map[string]string)
			}
			k8sMetadata.Labels[pc.ChartSubcomponentKeyName] = pc.ChartSubcomponentValue
		case "annotations":
			if k8sMetadata.Annotations == nil {
				k8sMetadata.Annotations = make(map[string]string)
			}
			k8sMetadata.Annotations[pc.ChartSubcomponentKeyName] = pc.ChartSubcomponentValue
		case "name":
			k8sMetadata.Name = pc.ChartSubcomponentValue
		}
	}
}

func ConvertPodSpecField(k8sPodTempSpec *v1.PodTemplateSpec, pcSlice []common_conversions.PC) {
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		switch subClassName {
		case "shareProcessNamespace":
			sharedBool := false
			if pc.ChartSubcomponentValue == "true" {
				sharedBool = true
			}
			k8sPodTempSpec.Spec.ShareProcessNamespace = &sharedBool
		}
	}
	return
}
