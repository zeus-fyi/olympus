package zeus_core

import (
	"context"

	"github.com/rs/zerolog/log"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) DeletePod(ctx context.Context, name string, kubeCtxNs zeus_common_types.CloudCtxNs, deletePodOpts *metav1.DeleteOptions) error {
	log.Ctx(ctx).Debug().Msg("DeletePod")
	k.SetContext(kubeCtxNs.Context)
	opts := metav1.DeleteOptions{}
	if deletePodOpts == nil {
		deletePodOpts = &opts
	}
	err := k.kc.CoreV1().Pods(kubeCtxNs.Namespace).Delete(ctx, name, *deletePodOpts)
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (k *K8Util) DeleteFirstPodLike(ctx context.Context, kubeCtxNs zeus_common_types.CloudCtxNs, podName string, deletePodOpts *metav1.DeleteOptions, filter *strings_filter.FilterOpts) error {
	log.Ctx(ctx).Debug().Msg("DeleteFirstPodLike")
	k.SetContext(kubeCtxNs.Context)

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

func (k *K8Util) DeleteAllPodsLike(ctx context.Context, kubeCtxNs zeus_common_types.CloudCtxNs, podName string, deletePodOpts *metav1.DeleteOptions, filter *strings_filter.FilterOpts) error {
	log.Ctx(ctx).Debug().Msg("DeleteAllPodsLike")
	k.SetContext(kubeCtxNs.Context)
	if filter == nil {
		filter = &strings_filter.FilterOpts{
			DoesNotStartWithThese: nil,
			StartsWithThese:       nil,
			StartsWith:            "",
			Contains:              podName,
			DoesNotInclude:        nil,
		}
	}
	pods, err := k.GetPodsUsingCtxNs(ctx, kubeCtxNs, nil, filter)
	log.Ctx(ctx).Err(err).Msg("DeleteAllPodsLike")
	if err != nil && errors.IsNotFound(err) {
		log.Ctx(ctx).Err(err).Msg("DeleteAllPodsLike, Pods Like Not Found")
		return err
	}
	opts := metav1.DeleteOptions{}
	if deletePodOpts == nil {
		deletePodOpts = &opts
	}
	p := v1.Pod{}
	for _, pod := range pods.Items {
		name := pod.ObjectMeta.Name
		if strings_filter.FilterStringWithOpts(name, filter) {
			p = pod
			err = k.kc.CoreV1().Pods(kubeCtxNs.Namespace).Delete(ctx, p.GetName(), *deletePodOpts)
			if err != nil && errors.IsNotFound(err) {
				log.Ctx(ctx).Err(err).Msg("DeleteAllPodsLike, Pods Like Not Found")
				return err
			}
		}
	}
	return err
}
