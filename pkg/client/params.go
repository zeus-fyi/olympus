package client

import "github.com/zeus-fyi/olympus/pkg/utils/env"

type Endpoint string

func (e Endpoint) Endpoint() string {
	if len(e) > 0 {
		return string(e)
	}
	return env.SetEnvParam(e)
}

func (e Endpoint) Local() string {
	return "http://localhost"
}
func (e Endpoint) Dev() string {
	return ""
}
func (e Endpoint) Staging() string {
	return ""
}
func (e Endpoint) Production() string {
	return ""
}
