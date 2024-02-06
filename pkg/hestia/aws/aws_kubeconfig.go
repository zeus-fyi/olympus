package hestia_eks_aws

type KubeConfig struct {
	APIVersion     string         `json:"apiVersion"`
	Kind           string         `json:"kind"`
	Clusters       []ClusterEntry `json:"clusters"`
	Contexts       []ContextEntry `json:"contexts"`
	CurrentContext string         `json:"current-context"`
	Users          []UserEntry    `json:"users"`
}

type ClusterEntry struct {
	Name    string      `json:"name"`
	Cluster ClusterInfo `json:"cluster"`
}

type ClusterInfo struct {
	Server                   string `json:"server"`
	CertificateAuthorityData string `json:"certificate-authority-data"`
}

type ContextEntry struct {
	Name    string      `json:"name"`
	Context ContextInfo `json:"context"`
}

type ContextInfo struct {
	Cluster string `json:"cluster"`
	User    string `json:"user"`
}

type UserEntry struct {
	Name string   `json:"name"`
	User UserInfo `json:"user"`
}

type UserInfo struct {
	Exec ExecConfig `json:"exec"`
}

type ExecConfig struct {
	APIVersion string   `json:"apiVersion"`
	Command    string   `json:"command"`
	Args       []string `json:"args"`
}
