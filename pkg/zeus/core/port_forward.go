package zeus_core

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	pf "k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/kubectl/pkg/cmd/portforward"
)

type portForwarder interface {
	ForwardPorts(method string, url *url.URL, opts portforward.PortForwardOptions) error
}

type defaultPortForwarder struct {
	genericclioptions.IOStreams
}

func (f *defaultPortForwarder) ForwardPorts(method string, url *url.URL, opts portforward.PortForwardOptions) error {
	transport, upgrader, err := spdy.RoundTripperFor(opts.Config)
	if err != nil {
		log.Err(err)
		return err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, method, url)
	fw, err := pf.NewOnAddresses(dialer, opts.Address, opts.Ports, opts.StopChannel, opts.ReadyChannel, f.Out, f.ErrOut)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}

func (k *K8Util) PortForwardPod(ctx context.Context, kubeCtxNs zeus_common_types.CloudCtxNs, podName, address string, ports []string, readyChan, stopChan chan struct{}, filter *strings_filter.FilterOpts) error {
	log.Ctx(ctx).Debug().Msg("PortForwardPod")
	k.SetContext(kubeCtxNs.Context)

	p, err := k.GetPodsUsingCtxNs(ctx, kubeCtxNs, nil, filter)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return err
	}
	pod, err := k.getFirstPodByPrefix(podName, p)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return err
	}
	if pod.Status.Phase != v1.PodRunning {
		return fmt.Errorf("unable to forward port because pod is not running. Current status=%v", pod.Status.Phase)
	}
	var localPF defaultPortForwarder
	podClient := k.kc.CoreV1()
	portFwd := portforward.PortForwardOptions{
		Namespace:     kubeCtxNs.Namespace,
		PodName:       pod.GetName(),
		Ports:         ports,
		PodClient:     podClient,
		Config:        k.clientCfg,
		PortForwarder: &localPF,
		ReadyChannel:  readyChan,
		StopChannel:   stopChan,
		Address:       []string{address},
	}

	req := podClient.RESTClient().Post().
		Resource("pods").
		Namespace(kubeCtxNs.Namespace).
		Name(pod.GetName()).
		SubResource("portforward")
	return portFwd.PortForwarder.ForwardPorts("POST", req.URL(), portFwd)
}
