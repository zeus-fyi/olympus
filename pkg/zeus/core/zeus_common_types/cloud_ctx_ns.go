package zeus_common_types

import "fmt"

type CloudCtxNs struct {
	CloudProvider string `json:"cloudProvider"`
	Region        string `json:"region"`
	Context       string `json:"context"`
	Namespace     string `json:"namespace"`
	Env           string `json:"env"`
}

func NewCloudCtxNs() CloudCtxNs {
	return CloudCtxNs{
		CloudProvider: "",
		Region:        "",
		Context:       "",
		Namespace:     "",
		Env:           "",
	}
}

func (kCtx *CloudCtxNs) GetCtxName(env string) string {
	return fmt.Sprintf("%s-%s-%s", kCtx.CloudProvider, kCtx.Region, kCtx.Context)
}
