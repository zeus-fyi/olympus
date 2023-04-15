package zeus_core

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeAudit struct {
	CloudProvider     string        `json:"cloudProvider"`
	KubernetesVersion string        `json:"kubernetesVersion"`
	NodeID            string        `json:"nodeID"`
	NodePoolID        string        `json:"nodePoolID"`
	Region            string        `json:"region"`
	Slug              string        `json:"slug"`
	Taints            []v1.Taint    `json:"taints"`
	Status            v1.NodeStatus `json:"status"`
}

func (k *K8Util) GetNodesAuditByLabel(ctx context.Context, kns zeus_common_types.CloudCtxNs, label string) ([]NodeAudit, error) {
	k.SetContext(kns.Context)
	nl, err := k.kc.CoreV1().Nodes().List(ctx, metav1.ListOptions{LabelSelector: label})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error getting nodes by label")
		return nil, err
	}
	nodesAudit := make([]NodeAudit, len(nl.Items))
	for i, n := range nl.Items {
		na := NodeAudit{}
		na.Status = n.Status
		na.Taints = n.Spec.Taints
		for key, v := range n.Labels {
			switch key {
			case "region":
				na.Region = v
			case "node.kubernetes.io/instance-type":
				na.Slug = v
			case "doks.digitalocean.com/node-pool-id":
				na.NodePoolID = v
			case "doks.digitalocean.com/node-id":
				na.NodeID = v
			case "doks.digitalocean.com/version":
				na.KubernetesVersion = v
				na.CloudProvider = "do"
			}
		}
		nodesAudit[i] = na
	}
	return nodesAudit, nil
}

func (k *K8Util) GetNodesByLabel(ctx context.Context, kns zeus_common_types.CloudCtxNs, label string) (*v1.NodeList, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().Nodes().List(ctx, metav1.ListOptions{LabelSelector: label})
}

func (k *K8Util) GetNodes(ctx context.Context, kns zeus_common_types.CloudCtxNs) (*v1.NodeList, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
}
