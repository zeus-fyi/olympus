package hestia_quiknode

type ProvisionRequest struct {
	QuicknodeID       string        `json:"quicknode-id"`
	EndpointID        string        `json:"endpoint-id"`
	WssUrl            string        `json:"wss-url"`
	HttpUrl           string        `json:"http-url"`
	Referers          []string      `json:"referers"`
	ContractAddresses []interface{} `json:"contract_addresses"`
	Chain             string        `json:"chain"`
	Network           string        `json:"network"`
	Plan              string        `json:"plan"`
}

type CreateResponse struct {
	Status       string `json:"status"`
	DashboardUrl string `json:"dashboard-url"`
	AccessUrl    string `json:"access-url"`
}

type UpdateResponse struct {
	Status string `json:"status"`
}

type DeactivateRequest struct {
	QuicknodeID  string `json:"quicknode-id"`
	EndpointID   string `json:"endpoint-id"`
	DeactivateAt int64  `json:"deactivate-at"`
	Chain        string `json:"chain"`
	Network      string `json:"network"`
}

type DeprovisionRequest struct {
	QuicknodeId   string `json:"quicknode-id"`
	EndpointId    string `json:"endpoint-id"`
	DeprovisionAt int64  `json:"deprovision-at"`
}
