package autok8s_core

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/zeus-fyi/olympus/pkg/utils/printer"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type K8Util struct {
	kc        *kubernetes.Clientset
	cfgAccess clientcmd.ConfigAccess
	clientCfg *rest.Config

	CfgPath string
	env     string

	FileName  string
	PrintPath string
	PrintOn   bool
}

type KubeCtxNs struct {
	CloudProvider string
	Region        string
	CtxType       string
	Namespace     string
	Env           string
}

func (kCtx *KubeCtxNs) GetCtxName(env string) string {
	return fmt.Sprintf("%s-%s-%s", kCtx.CloudProvider, kCtx.Region, kCtx.CtxType)
}

func (k *K8Util) GetNamespaces() (*v1.NamespaceList, error) {
	return k.kc.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
}

func (k *K8Util) GetContexts() (map[string]*clientcmdapi.Context, error) {
	startingConfig, err := k.cfgAccess.GetStartingConfig()
	return startingConfig.Contexts, err
}

func (k *K8Util) SetContext(context string) {
	var err error

	cfgOveride := &clientcmd.ConfigOverrides{}
	if len(context) > 0 {
		cfgOveride = &clientcmd.ConfigOverrides{
			CurrentContext: context}
	}

	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: k.CfgPath},
		cfgOveride)
	k.cfgAccess = cc.ConfigAccess()
	k.clientCfg, err = cc.ClientConfig()
	if err != nil {
		log.Panicln("Failed to set context")
	}
	k.SetClient(k.clientCfg)
}

func (k *K8Util) SetClient(config *rest.Config) {
	var err error
	k.kc, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln("Failed to set client")
	}
}

func (k *K8Util) ConnectToK8s() {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}

	k.CfgPath = filepath.Join(home, ".kube", "config")
	k.SetContext("")
}

func (k *K8Util) ConnectToK8sFromConfig(dir string) {
	k.CfgPath = dir
	k.SetContext("")
}

func (k *K8Util) K8Printer(v interface{}, env string) (interface{}, error) {
	if k.PrintOn && k.FileName != "" {
		if k.PrintPath == "" && env != "" {
			var printPath printer.PrintPath
			switch env {
			case "dev", "development":
				k.PrintPath = printPath.Dev()
			case "staging":
				k.PrintPath = printPath.Staging()
			case "production":
				k.PrintPath = printPath.Production()
			}
		}
		return printer.InterfacePrinter(k.PrintPath, k.FileName, v)
	}
	return v, nil
}
