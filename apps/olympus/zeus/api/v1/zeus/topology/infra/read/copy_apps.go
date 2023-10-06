package read_infra

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	sui_cookbooks "github.com/zeus-fyi/zeus/cookbooks/sui"
	zeus_client "github.com/zeus-fyi/zeus/zeus/z_client"
)

var (
	CookbooksDirIn = "/etc/cookbooks"
)

func CopySuiApp(ctx context.Context, appName, bearer string) error {
	err := os.Chdir(CookbooksDirIn)
	if err != nil {
		return err
	}
	cfg := sui_cookbooks.SuiConfigOpts{
		WithLocalNvme:        true,
		DownloadSnapshot:     true,
		WithIngress:          true,
		WithServiceMonitor:   true,
		WithArchivalFallback: true,
	}
	sui_cookbooks.SuiMasterChartPath.DirIn = path.Join(CookbooksDirIn, "/sui/node/infra")
	sui_cookbooks.SuiIngressChartPath.DirIn = path.Join(CookbooksDirIn, "/sui/node/ingress")
	sui_cookbooks.SuiServiceMonitorChartPath.DirIn = path.Join(CookbooksDirIn, "/sui/node/servicemonitor")
	nameSplit := strings.Split(appName, "-")
	if len(nameSplit) != 3 {
		return fmt.Errorf("invalid name: %s", nameSplit)
	}
	net := nameSplit[1]
	cp := nameSplit[2]
	cfg.Network = net
	cfg.CloudProvider = cp
	suiNodeDefinition := sui_cookbooks.GetSuiClientClusterDef(cfg)
	gcd := suiNodeDefinition.BuildClusterDefinitions()
	zc := zeus_client.NewDefaultZeusClient(bearer)
	err = gcd.CreateClusterClassDefinitions(ctx, zc)
	if err != nil {
		return err
	}
	_, err = suiNodeDefinition.UploadChartsFromClusterDefinition(ctx, zc, false)
	if err != nil {
		return err
	}
	return nil
}
