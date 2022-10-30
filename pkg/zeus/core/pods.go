package zeus_core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/cmd/portforward"
)

func (k *K8Util) PortForwardPod(ctx context.Context, kubeCtxNs KubeCtxNs, podName, address string, ports []string, readyChan, stopChan chan struct{}, filter *string_utils.FilterOpts) error {
	log.Ctx(ctx).Debug().Msg("PortForwardPod")

	p, err := k.GetPodsUsingCtxNs(ctx, kubeCtxNs, nil, filter)
	if err != nil {
		return err
	}
	pod, err := k.getFirstPodByPrefix(podName, p)
	if err != nil {
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

func (k *K8Util) GetPod(ctx context.Context, name, ns string) (*v1.Pod, error) {
	log.Ctx(ctx).Debug().Msg("GetPod")
	p, err := k.kc.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	return p, err
}

func (k *K8Util) GetPodLogs(ctx context.Context, name, ns string, logOpts *v1.PodLogOptions, filter *string_utils.FilterOpts) ([]byte, error) {
	log.Ctx(ctx).Debug().Msg("GetPodLogs")
	if logOpts == nil {
		logOpts = &v1.PodLogOptions{}
	}
	req := k.kc.CoreV1().Pods(ns).GetLogs(name, logOpts)
	buf := new(bytes.Buffer)
	podLogs, err := req.Stream(ctx)
	defer func(podLogs io.ReadCloser) {
		closeErr := podLogs.Close()
		if closeErr != nil {
			fmt.Printf("%s", closeErr.Error())
		}
	}(podLogs)
	if err != nil {
		return buf.Bytes(), err
	}
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), err
}

func (k *K8Util) GetPods(ctx context.Context, ns string, opts metav1.ListOptions) (*v1.PodList, error) {
	return k.kc.CoreV1().Pods(ns).List(context.Background(), opts)
}

func (k *K8Util) GetPodsUsingCtxNs(ctx context.Context, kubeCtxNs KubeCtxNs, logOpts *v1.PodLogOptions, filter *string_utils.FilterOpts) (*v1.PodList, error) {
	log.Ctx(ctx).Debug().Msg("GetPodsUsingCtxNs")
	if logOpts == nil {
		logOpts = &v1.PodLogOptions{}
	}
	pods, err := k.GetPods(ctx, kubeCtxNs.Namespace, metav1.ListOptions{})
	if err != nil {
		return pods, err
	}

	if filter != nil {
		filteredPods := v1.PodList{}
		for _, pod := range pods.Items {
			if string_utils.FilterStringWithOpts(pod.GetName(), filter) {
				filteredPods.Items = append(filteredPods.Items, pod)
			}
		}
		_, err = k.K8Printer(filteredPods, kubeCtxNs.Env)
		return &filteredPods, nil
	}

	_, err = k.K8Printer(pods, kubeCtxNs.Env)
	return pods, err
}

func (k *K8Util) GetFirstPodLike(ctx context.Context, kubeCtxNs KubeCtxNs, podName string, filter *string_utils.FilterOpts) (*v1.Pod, error) {
	pods, err := k.GetPodsUsingCtxNs(ctx, kubeCtxNs, nil, filter)
	if err != nil {
		return nil, err
	}
	return k.getFirstPodLike(ctx, podName, pods)
}

func (k *K8Util) getFirstPodLike(ctx context.Context, podName string, pl *v1.PodList) (*v1.Pod, error) {
	p := v1.Pod{}
	for _, pod := range pl.Items {
		name := pod.ObjectMeta.Name
		if strings.Contains(name, podName) {
			p = pod
		}
	}
	return &p, errors.New("pod not found")
}

func (k *K8Util) getFirstPodByPrefix(podName string, pl *v1.PodList) (*v1.Pod, error) {
	for _, pod := range pl.Items {
		name := pod.ObjectMeta.Name
		if strings.HasPrefix(name, podName) {
			return &pod, nil
		}
	}
	return nil, errors.New("pod not found")
}
