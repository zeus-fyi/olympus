package zeus_core

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeAudit struct {
	KubernetesVersion string        `json:"kubernetesVersion"`
	NodeID            string        `json:"nodeID"`
	NodePoolID        string        `json:"nodePoolID"`
	Slug              string        `json:"slug"`
	Taints            []v1.Taint    `json:"taints"`
	Status            v1.NodeStatus `json:"status"`
}

type ClusterNodesAudit struct {
	CloudCtxNs zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`
	Nodes      []NodeAudit                  `json:"nodes"`
}

func (k *K8Util) GetNodesAuditByLabel(ctx context.Context, kns zeus_common_types.CloudCtxNs, label string) (*ClusterNodesAudit, error) {
	k.SetContext(kns.Context)
	nl, err := k.kc.CoreV1().Nodes().List(ctx, metav1.ListOptions{LabelSelector: label})
	if err != nil {
		log.Error().Err(err).Msg("error getting nodes by label")
		return nil, err
	}
	nodesAudit := make([]NodeAudit, len(nl.Items))
	for i, n := range nl.Items {
		na := NodeAudit{}
		na.Status = n.Status
		na.Taints = n.Spec.Taints
		for key, v := range n.Labels {
			switch key {
			case "region", "topology.kubernetes.io/region":
			case "node.kubernetes.io/instance-type":
				na.Slug = v
			case "doks.digitalocean.com/node-pool-id", "nodepool":
				na.NodePoolID = v
			case "doks.digitalocean.com/node-id":
				na.NodeID = v
			case "doks.digitalocean.com/version":
				na.KubernetesVersion = v
			}
		}

		switch kns.CloudProvider {
		case "ovh":
			na.NodePoolID = n.ObjectMeta.Name

		}
		nodesAudit[i] = na
	}

	cp := &ClusterNodesAudit{
		CloudCtxNs: kns,
		Nodes:      nodesAudit,
	}
	return cp, nil
}

func (k *K8Util) GetNodesByLabel(ctx context.Context, kns zeus_common_types.CloudCtxNs, label string) (*v1.NodeList, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().Nodes().List(ctx, metav1.ListOptions{LabelSelector: label})
}

func (k *K8Util) GetNodes(ctx context.Context, kns zeus_common_types.CloudCtxNs) (*v1.NodeList, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
}
