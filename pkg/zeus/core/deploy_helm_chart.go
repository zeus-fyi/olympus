package zeus_core

import (
	"context"
	"fmt"
	"os"

	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

var HelmRootPath string

func (k *K8Util) DeployHelm(ctx context.Context, kns zeus_common_types.CloudCtxNs) error {
	k.SetContext(kns.Context)

	// Initialize Helm settings
	settings := cli.New()

	// Set up action configurations
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), kns.Namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Printf(format, v...)
	}); err != nil {
		return err
	}

	client := action.NewInstall(actionConfig)
	client.ReleaseName = "my-release"
	client.Namespace = kns.Namespace
	//
	//valsDefault := PromGrafanaDefaultValues()
	//
	//valueOpts := &values.Options{}
	//vals, err := valueOpts.MergeValues([]map[string]interface{}{valsDefault})
	//if err != nil {
	//	return err
	//}
	//// Load the chart from a repository URL or local path
	//chartPath := "prometheus-community/kube-prometheus-stack" // This should be adjusted based on actual chart location
	//chart, err := loader.Load(chartPath)
	//if err != nil {
	//	return err
	//} // Set up the install action
	return nil
}

func PromGrafanaDefaultValues() map[string]interface{} {

	// Define the custom values
	vals := map[string]interface{}{
		"grafana": map[string]interface{}{
			"persistence": map[string]interface{}{
				"accessModes": []string{"ReadWriteOnce"},
				"enabled":     true,
				"finalizers":  []string{"kubernetes.io/pvc-protection"},
				"size":        "10Gi",
				"type":        "pvc",
			},
			"resources": map[string]interface{}{
				"limits": map[string]interface{}{
					"cpu":    "1",
					"memory": "5Gi",
				},
				"requests": map[string]interface{}{
					"cpu":    "1",
					"memory": "5Gi",
				},
			},
		},
		"prometheus": map[string]interface{}{
			"prometheusSpec": map[string]interface{}{
				"resources": map[string]interface{}{
					"limits": map[string]interface{}{
						"cpu":    "2",
						"memory": "10Gi",
					},
					"requests": map[string]interface{}{
						"cpu":    "2",
						"memory": "10Gi",
					},
				},
				"storageSpec": map[string]interface{}{
					"volumeClaimTemplate": map[string]interface{}{
						"spec": map[string]interface{}{
							"accessModes": []string{"ReadWriteOnce"},
							"resources": map[string]interface{}{
								"requests": map[string]interface{}{
									"storage": "2Ti",
								},
							},
						},
					},
				},
			},
		},
	}
	return vals
}
