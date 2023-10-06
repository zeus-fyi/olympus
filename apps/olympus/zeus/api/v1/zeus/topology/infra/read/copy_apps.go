package read_infra

import (
	"context"
	"os"
	"path"

	sui_cookbooks "github.com/zeus-fyi/zeus/cookbooks/sui"
	zeus_client "github.com/zeus-fyi/zeus/zeus/z_client"
)

var (
	CookbooksDirIn = "/etc/cookbooks"
)

func CopySuiApp(ctx context.Context, bearer string) error {
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

	cps := []string{"aws", "gcp", "do"}
	networks := []string{"mainnet", "testnet", "devnet"}
	for _, cp := range cps {
		for _, network := range networks {
			cfg.Network = network
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
		}
	}
	return nil
}
