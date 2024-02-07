package zeus_core

import (
	"fmt"
	"os"
	"path/filepath"

	monitoringclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	GcpContext          = "gke_zeusfyi_us-central1-a_zeus-gcp-pilot-0"
	DigitalOceanContext = "do-nyc1-do-nyc1-zeus-demo"
	zeusfyi             = "zeusfyi"
	zeusfyiShared       = "zeusfyi-shared"
)

type K8Util struct {
	kc        *kubernetes.Clientset
	cfgAccess clientcmd.ConfigAccess
	mc        *monitoringclient.Clientset
	kcCfg     clientcmd.ClientConfig
	clientCfg *rest.Config

	CfgPath string
	env     string

	FileName  string
	PrintPath string
	PrintOn   bool
}

type FilterOpts struct {
	DoesNotInclude []string
}

func (k *K8Util) GetRawConfigs() (clientcmdapi.Config, error) {
	cfgs, err := k.kcCfg.RawConfig()
	if err != nil {
		log.Err(err).Msg("Zeus: GetRawConfigs, failed to get raw config")
		return cfgs, err
	}
	return cfgs, err
}

func (k *K8Util) GetContexts() (map[string]*clientcmdapi.Context, error) {
	cfg, err := k.kcCfg.RawConfig()
	if err != nil {
		log.Err(err)
		return nil, err
	}
	return cfg.Contexts, err
}

func (k *K8Util) SetContext(context string) {
	switch context {
	case "zeus-us-west-1":
		context = "arn:aws:eks:us-west-1:480391564655:cluster/zeus-us-west-1"
	case zeusfyi:
		context = fmt.Sprintf("kubernetes-admin@%s", zeusfyi)
	case zeusfyiShared:
		context = fmt.Sprintf("kubernetes-admin@%s", zeusfyiShared)
	}
	var err error
	rc, err := k.kcCfg.RawConfig()
	if err != nil {
		log.Err(err)
	}
	cc := clientcmd.NewNonInteractiveClientConfig(rc, context, nil, k.cfgAccess)
	k.cfgAccess = cc.ConfigAccess()
	k.clientCfg, err = cc.ClientConfig()
	if err != nil {
		log.Err(err)
	}
	k.SetClient(k.clientCfg)
	mclient, err := monitoringclient.NewForConfig(k.clientCfg)
	if err != nil {
		log.Panic().Msg("Zeus: NewForConfig, failed to set client config")
		misc.DelayedPanic(err)
	}
	k.mc = mclient
}

func (k *K8Util) SetClient(config *rest.Config) {
	var err error
	k.kc, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic().Msg("Zeus: SetClient, failed to set client")
		misc.DelayedPanic(err)
	}
}

func (k *K8Util) ConnectToK8s() {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}
	b, err := os.ReadFile(filepath.Join(home, ".kube", "config"))
	if err != nil {
		log.Panic().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath, failed to read inmemfs kube config")
		misc.DelayedPanic(err)
	}
	cc, err := clientcmd.NewClientConfigFromBytes(b)
	if err != nil {
		log.Panic().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath, failed to set context")
		misc.DelayedPanic(err)
	}
	k.kcCfg = cc
	k.cfgAccess = cc.ConfigAccess()
	k.clientCfg, err = cc.ClientConfig()
	if err != nil {
		log.Panic().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath, failed to set client config")
		misc.DelayedPanic(err)
	}
	k.SetClient(k.clientCfg)
	mclient, err := monitoringclient.NewForConfig(k.clientCfg)
	if err != nil {
		log.Panic().Msg("Zeus: NewForConfig, failed to set client config")
		misc.DelayedPanic(err)
	}
	k.mc = mclient
	log.Info().Msg("Zeus: DefaultK8sCfgPath complete")
	k.CfgPath = filepath.Join(home, ".kube", "config")
}

func (k *K8Util) ConnectToK8sFromConfig(dir string) {
	k.CfgPath = dir
}

func (k *K8Util) DefaultK8sCfgPath() string {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}
	b, err := os.ReadFile(filepath.Join(home, ".kube", "config"))
	if err != nil {
		log.Panic().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath, failed to read inmemfs kube config")
		misc.DelayedPanic(err)
	}
	cc, err := clientcmd.NewClientConfigFromBytes(b)
	if err != nil {
		log.Panic().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath, failed to set context")
		misc.DelayedPanic(err)
	}
	k.kcCfg = cc
	k.cfgAccess = cc.ConfigAccess()
	k.clientCfg, err = cc.ClientConfig()
	if err != nil {
		log.Panic().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath, failed to set client config")
		misc.DelayedPanic(err)
	}
	k.SetClient(k.clientCfg)
	mclient, err := monitoringclient.NewForConfig(k.clientCfg)
	if err != nil {
		log.Panic().Msg("Zeus: NewForConfig, failed to set client config")
		misc.DelayedPanic(err)
	}
	k.mc = mclient
	log.Info().Msg("Zeus: DefaultK8sCfgPath complete")

	return filepath.Join(home, ".kube", "config")
}

func (k *K8Util) ConnectToK8sFromInMemFsCfgPath(fs memfs.MemFS) {
	log.Info().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath starting")
	var err error
	b, err := fs.ReadFile("/.kube/config")
	if err != nil {
		log.Panic().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath, failed to read inmemfs kube config")
		misc.DelayedPanic(err)
	}
	cc, err := clientcmd.NewClientConfigFromBytes(b)
	if err != nil {
		log.Panic().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath, failed to set context")
		misc.DelayedPanic(err)
	}
	k.kcCfg = cc
	k.cfgAccess = cc.ConfigAccess()
	k.clientCfg, err = cc.ClientConfig()
	if err != nil {
		log.Panic().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath, failed to set client config")
		misc.DelayedPanic(err)
	}
	k.SetClient(k.clientCfg)
	mclient, err := monitoringclient.NewForConfig(k.clientCfg)
	if err != nil {
		log.Panic().Msg("Zeus: NewForConfig, failed to set client config")
		misc.DelayedPanic(err)
	}
	k.mc = mclient
	log.Info().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath complete")
}

func (k *K8Util) ConnectToK8sFromInMemFsCfgPathOrErr(fs memfs.MemFS) error {
	log.Info().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath starting")
	var err error
	b, err := fs.ReadFile("/.kube/config")
	if err != nil {
		log.Err(err).Msg("Zeus: ConnectToK8sFromInMemFsCfgPath, failed to read inmemfs kube config")
		return err
	}
	cc, err := clientcmd.NewClientConfigFromBytes(b)
	if err != nil {
		log.Err(err).Msg("Zeus: ConnectToK8sFromInMemFsCfgPath, failed to set context")
		misc.DelayedPanic(err)
	}
	k.kcCfg = cc
	k.cfgAccess = cc.ConfigAccess()
	k.clientCfg, err = cc.ClientConfig()
	if err != nil {
		log.Err(err).Msg("Zeus: ConnectToK8sFromInMemFsCfgPath, failed to set client config")
		misc.DelayedPanic(err)
	}
	k.SetClient(k.clientCfg)
	mclient, err := monitoringclient.NewForConfig(k.clientCfg)
	if err != nil {
		log.Err(err).Msg("Zeus: NewForConfig, failed to set client config")
		misc.DelayedPanic(err)
	}
	k.mc = mclient
	log.Info().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath complete")
	return nil
}

func (k *K8Util) K8Printer(v interface{}, env string) (interface{}, error) {
	//if k.PrintOn && k.FileName != "" {
	//	if k.PrintPath == "" && env != "" {
	//		var printPath file_io.PrintPath
	//		switch env {
	//		case "dev", "development":
	//			k.PrintPath = printPath.Dev()
	//		case "staging":
	//			k.PrintPath = printPath.Staging()
	//		case "production":
	//			k.PrintPath = printPath.Production()
	//		}
	//	}
	//	return file_io.InterfacePrinter(k.PrintPath, k.FileName, env, v)
	//}
	return v, nil
}
