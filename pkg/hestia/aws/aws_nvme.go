package hestia_eks_aws

func AddAwsEksNvmeLabels(labels map[string]string) map[string]string {
	labels["fast-disk-node"] = "pv-raid"
	return labels
}
