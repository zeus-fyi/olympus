package ingresses

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"

func (i *Ingress) ParseK8sConfigToDB() error {
	i.Metadata.ChartSubcomponentParentClassTypeName = "IngressParentMetadata"
	metadata := common_conversions.ConvertMetadata(i.K8sIngress.ObjectMeta)
	i.Metadata.Metadata = metadata
	err := i.ConvertIngressSpecConfigToDB()
	if err != nil {
		return err
	}
	return err
}

func (i *Ingress) ConvertIngressSpecConfigToDB() error {
	i.Spec = NewIngressSpec()
	err := i.ConvertK8sIngressRuleToDB()
	if err != nil {
		return err
	}
	i.ConvertK8sIngressTLSToDB()
	return err
}
