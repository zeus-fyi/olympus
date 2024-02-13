package deploy_topology_activities_create_setup

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
)

var HelmRootPath string

func DeployHelmChart() {
	actionConfig := new(action.Configuration)

	// Set up the install action
	client := action.NewInstall(actionConfig)
	client.ReleaseName = "my-release"
	client.Namespace = "default"

	// Load the chart
	chartPath := "/path/to/chart"
	chart, err := loader.Load(chartPath)
	if err != nil {
		panic(err)
	}

	// Install the chart
	_, err = client.Run(chart, nil) // nil because we're not passing any specific values
	if err != nil {
		panic(err)
	}

}
