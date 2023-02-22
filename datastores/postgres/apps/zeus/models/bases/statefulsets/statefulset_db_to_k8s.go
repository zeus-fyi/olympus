package statefulsets

import (
	"strings"

	v1core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
)

func (s *StatefulSet) ParseDBConfigToK8s(pcSlice common_conversions.ParentChildDB) error {
	dbStatefulSet := NewStatefulSet()
	dbStatefulSet.Metadata.Metadata = common_conversions.CreateMetadataByFields(s.K8sStatefulSet.Name, s.K8sStatefulSet.Annotations, s.K8sStatefulSet.Labels)

	pvcMap := make(map[int][]common_conversions.PC)
	for pcGroupName, pc := range pcSlice.PCGroupMap {

		switch pcGroupName {
		case "Spec":
			err := s.ConvertDBSpecToK8s(pc)
			if err != nil {
				return err
			}
		case "PodTemplateSpec":
			db_to_k8s_conversions.ConvertPodSpecField(&s.K8sStatefulSet.Spec.Template, pc)
		case "StatefulSetParentMetadata":
			db_to_k8s_conversions.ConvertMetadata(&s.K8sStatefulSet.ObjectMeta, pc)
		case "PodTemplateSpecMetadata":
			db_to_k8s_conversions.ConvertMetadata(&s.K8sStatefulSet.Spec.Template.ObjectMeta, pc)
		case "VolumeClaimTemplate":
			for _, p := range pc {
				pcID := p.ParentClassTypesDB.ChartSubcomponentParentClassTypeID
				tmp := pvcMap[pcID]
				pvcMap[pcID] = append(tmp, p)
			}
		}
	}
	_ = s.ParseVolumeClaimTemplates(pvcMap)

	return nil
}

func (s *StatefulSet) ParseVolumeClaimTemplates(pvcMap map[int][]common_conversions.PC) error {

	pvc := v1core.PersistentVolumeClaim{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1core.PersistentVolumeClaimSpec{},
		Status:     v1core.PersistentVolumeClaimStatus{},
	}
	for _, v := range pvcMap {
		for _, p := range v {
			key := p.ChartSubcomponentKeyName
			val := p.ChartSubcomponentValue
			switch p.ChartSubcomponentChildClassTypeName {
			case "VolumeClaimTemplateMetadata":
				switch key {
				case "name":
					pvc.Name = val
				case "labels":
					if pvc.Labels == nil {
						pvc.Labels = make(map[string]string)
					}
					pvc.Labels[key] = val
				case "annotations":
					if pvc.Annotations == nil {
						pvc.Annotations = make(map[string]string)
					}
					pvc.Annotations[key] = val
				}
			case "VolumeClaimTemplateSpec":
				switch key {
				case "storageClassName":
					pvc.Spec.StorageClassName = &val
				case "requests":
					rl := v1core.ResourceList{}
					val = strings.Trim(val, `""`)
					qty := resource.MustParse(val)
					rl["storage"] = qty
					pvc.Spec.Resources.Requests = rl
				case "accessMode":
					if len(pvc.Spec.AccessModes) <= 0 {
						pvc.Spec.AccessModes = []v1core.PersistentVolumeAccessMode{}
					}
					switch val {
					case "ReadWriteOnce":
						pvc.Spec.AccessModes = append(pvc.Spec.AccessModes, v1core.ReadWriteOnce)
					case "ReadOnlyMany":
						pvc.Spec.AccessModes = append(pvc.Spec.AccessModes, v1core.ReadOnlyMany)
					case "ReadWriteMany":
						pvc.Spec.AccessModes = append(pvc.Spec.AccessModes, v1core.ReadWriteMany)
					case "ReadWriteOncePod":
						pvc.Spec.AccessModes = append(pvc.Spec.AccessModes, v1core.ReadWriteOncePod)
					}
				}
			}
		}
		s.K8sStatefulSet.Spec.VolumeClaimTemplates = append(s.K8sStatefulSet.Spec.VolumeClaimTemplates, pvc)
	}
	return nil
}
