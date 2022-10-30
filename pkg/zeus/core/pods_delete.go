package zeus_core

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) DeletePod(ctx context.Context, name, ns string, deletePodOpts *metav1.DeleteOptions) error {
	log.Ctx(ctx).Debug().Msg("DeletePod")

	opts := metav1.DeleteOptions{}
	if deletePodOpts == nil {
		deletePodOpts = &opts
	}
	err := k.kc.CoreV1().Pods(ns).Delete(ctx, name, *deletePodOpts)
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (k *K8Util) DeleteFirstPodLike(ctx context.Context, kubeCtxNs KubeCtxNs, podName string, deletePodOpts *metav1.DeleteOptions, filter *string_utils.FilterOpts) error {
	log.Ctx(ctx).Debug().Msg("DeleteFirstPodLike")
	p, err := k.GetFirstPodLike(ctx, kubeCtxNs, podName, filter)
	if err != nil {
		return err
	}
	opts := metav1.DeleteOptions{}
	if deletePodOpts == nil {
		deletePodOpts = &opts
	}
	err = k.kc.CoreV1().Pods(kubeCtxNs.Namespace).Delete(ctx, p.GetName(), *deletePodOpts)
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (k *K8Util) DeleteAllPodsLike(ctx context.Context, kubeCtxNs KubeCtxNs, podName string, deletePodOpts *metav1.DeleteOptions, filter *string_utils.FilterOpts) error {
	log.Ctx(ctx).Debug().Msg("DeleteAllPodsLike")

	pods, err := k.GetPodsUsingCtxNs(ctx, kubeCtxNs, nil, filter)
	log.Ctx(ctx).Err(err).Msg("DeleteAllPodsLike")
	if !errors.IsNotFound(err) {
		log.Ctx(ctx).Err(err).Msg("DeleteAllPodsLike")
		return err
	}
	opts := metav1.DeleteOptions{}
	if deletePodOpts == nil {
		deletePodOpts = &opts
	}
	p := v1.Pod{}
	for _, pod := range pods.Items {
		name := pod.ObjectMeta.Name
		if strings.Contains(name, podName) {
			p = pod
			err = k.kc.CoreV1().Pods(kubeCtxNs.Namespace).Delete(ctx, p.GetName(), *deletePodOpts)
			if !errors.IsNotFound(err) {
				log.Ctx(ctx).Err(err).Msg("DeleteAllPodsLike")
				return err
			}
		}
	}
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}
