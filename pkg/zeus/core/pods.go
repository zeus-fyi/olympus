package zeus_core

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetPod(ctx context.Context, name string, kns zeus_common_types.CloudCtxNs) (*v1.Pod, error) {
	k.SetContext(kns.Context)
	log.Debug().Msg("GetPod")
	p, err := k.kc.CoreV1().Pods(kns.Namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Err(err).Interface("kns", kns).Msg("GetPod")
		return nil, err
	}
	return p, err
}

func (k *K8Util) GetPods(ctx context.Context, kns zeus_common_types.CloudCtxNs, opts *metav1.ListOptions) (*v1.PodList, error) {
	if opts == nil {
		opts = &metav1.ListOptions{}
	}
	k.SetContext(kns.Context)
	pods, err := k.kc.CoreV1().Pods(kns.Namespace).List(ctx, *opts)
	if err != nil {
		log.Err(err).Interface("kns", kns).Msg("GetPods")
		return nil, err
	}
	return pods, err
}

func (k *K8Util) GetPodsUsingCtxNs(ctx context.Context, kubeCtxNs zeus_common_types.CloudCtxNs, logOpts *v1.PodLogOptions, filter *strings_filter.FilterOpts) (*v1.PodList, error) {
	log.Debug().Msg("GetPodsUsingCtxNs")
	k.SetContext(kubeCtxNs.Context)

	if logOpts == nil {
		logOpts = &v1.PodLogOptions{}
	}
	pods, err := k.GetPods(ctx, kubeCtxNs, &metav1.ListOptions{})
	if err != nil {
		log.Err(err).Interface("kns", kubeCtxNs).Msg("GetPodsUsingCtxNs: GetPods")
		return pods, err
	}
	if filter != nil {
		filteredPods := v1.PodList{}
		for _, pod := range pods.Items {
			if strings_filter.FilterStringWithOpts(pod.GetName(), filter) {
				filteredPods.Items = append(filteredPods.Items, pod)
			}
		}
		return &filteredPods, nil
	}
	return pods, err
}

func (k *K8Util) GetFirstPodLike(ctx context.Context, kubeCtxNs zeus_common_types.CloudCtxNs, podName string, filter *strings_filter.FilterOpts) (*v1.Pod, error) {
	k.SetContext(kubeCtxNs.Context)
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
