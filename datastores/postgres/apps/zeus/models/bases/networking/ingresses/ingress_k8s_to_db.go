package ingresses

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"

func (i *Ingress) ConvertK8sIngressToDB() error {
	i.Metadata.ChartSubcomponentParentClassTypeName = "IngressParentMetadata"
	metadata := common_conversions.ConvertMetadata(i.K8sIngress.ObjectMeta)
	i.Metadata.Metadata = metadata
	err := i.ConvertK8sIngressSpecConfigToDB()
	if err != nil {
		return err
	}
	return err
}

func (i *Ingress) ConvertK8sIngressSpecConfigToDB() error {
	i.Spec = NewIngressSpec()

	if i.K8sIngress.Spec.IngressClassName != nil {
		i.NewIngressClassName(*i.K8sIngress.Spec.IngressClassName)
	}

	err := i.ConvertK8sIngressRuleToDB()
	if err != nil {
		return err
	}
	i.ConvertK8sIngressTLSToDB()
	return err
}
