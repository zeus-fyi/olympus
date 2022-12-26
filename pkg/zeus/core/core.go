package zeus_core

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type K8Util struct {
	kc        *kubernetes.Clientset
	cfgAccess clientcmd.ConfigAccess

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

func (k *K8Util) GetContexts() (map[string]*clientcmdapi.Context, error) {
	cfg, err := k.kcCfg.RawConfig()
	if err != nil {
		log.Err(err)
		return nil, err
	}
	return cfg.Contexts, err
}

func (k *K8Util) SetContext(context string) {
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
	log.Info().Msg("Zeus: ConnectToK8sFromInMemFsCfgPath complete")
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
