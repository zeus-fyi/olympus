package do_types

type DigitalOceanNodePoolRequestStatus struct {
	ExtClusterCfgID int    `json:"extClusterCfgID,omitempty"`
	ClusterID       string `json:"clusterID"`
	NodePoolID      string `json:"nodePoolID"`
}
