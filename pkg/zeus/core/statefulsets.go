package zeus_core

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetStatefulSetList(ctx context.Context, kubeCtxNs zeus_common_types.CloudCtxNs, filter *string_utils.FilterOpts) (*v1.StatefulSetList, error) {
	k.SetContext(kubeCtxNs.Context)
	opts := metav1.ListOptions{}
	ssl, err := k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).List(ctx, opts)
	return ssl, err
}

func (k *K8Util) GetStatefulSet(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) (*v1.StatefulSet, error) {
	k.SetContext(kns.Context)
	k.PrintPath = "stateful_sets"
	k.FileName = name
	opts := metav1.GetOptions{}
	ss, err := k.kc.AppsV1().StatefulSets(kns.Namespace).Get(ctx, name, opts)
	if err != nil {
		return ss, err
	}
	//_, err = k.K8Printer(ss, kns.Env)
	return ss, err
}

func (k *K8Util) DeleteStatefulSet(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) error {
	k.SetContext(kns.Context)
	opts := metav1.DeleteOptions{}
	err := k.kc.AppsV1().StatefulSets(kns.Namespace).Delete(ctx, name, opts)
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (k *K8Util) UpdateStatefulSet(ctx context.Context, kns zeus_common_types.CloudCtxNs, ss *v1.StatefulSet, filter *string_utils.FilterOpts) (*v1.StatefulSet, error) {
	k.SetContext(kns.Context)
	opts := metav1.UpdateOptions{}
	ss, err := k.kc.AppsV1().StatefulSets(kns.Namespace).Update(ctx, ss, opts)
	return ss, err
}

func (k *K8Util) CreateStatefulSet(ctx context.Context, kns zeus_common_types.CloudCtxNs, ss *v1.StatefulSet, filter *string_utils.FilterOpts) (*v1.StatefulSet, error) {
	k.SetContext(kns.Context)
	opts := metav1.CreateOptions{}
	ss, err := k.kc.AppsV1().StatefulSets(kns.Namespace).Create(ctx, ss, opts)
	alreadyExists := errors.IsAlreadyExists(err)
	if alreadyExists {
		log.Err(err).Interface("kns", kns).Msg("StatefulSet already exists, skipping creation")
		return ss, nil
	}
	return ss, err
}

func (k *K8Util) CreateStatefulSetIfVersionLabelChangesOrDoesNotExist(ctx context.Context, kns zeus_common_types.CloudCtxNs, nsts *v1.StatefulSet, filter *string_utils.FilterOpts) (*v1.StatefulSet, error) {
	k.SetContext(kns.Context)
	csts, err := k.GetStatefulSet(ctx, kns, nsts.Name, filter)
	switch {
	case csts != nil && len(csts.Name) > 0:
		switch IsVersionNew(csts.Labels, nsts.Labels) {
		case true:
			derr := k.DeleteStatefulSet(ctx, kns, csts.Name, filter)
			if derr != nil {
				return csts, derr
			}
		case false:
			return csts, nil
		}
	case errors.IsNotFound(err):
		newSts, newStsErr := k.CreateStatefulSet(ctx, kns, nsts, filter)
		return newSts, newStsErr
	}
	newSts, newStsErr := k.CreateStatefulSet(ctx, kns, nsts, filter)
	return newSts, newStsErr
}

func (k *K8Util) RolloutRestartStatefulSet(ctx context.Context, kubeCtxNs zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) error {
	k.SetContext(kubeCtxNs.Context)

	// Get the StatefulSet
	ss, err := k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Err(err).Interface("kubeCtxNs", kubeCtxNs).Msg("GetStatefulSet: error")
		return err
	}

	// Prepare for the restart
	if ss.Spec.Template.Annotations == nil {
		ss.Spec.Template.Annotations = make(map[string]string)
	}

	// Set a new annotation - this triggers a restart
	ss.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	// Update the StatefulSet
	_, err = k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).Update(ctx, ss, metav1.UpdateOptions{})
	if err != nil {
		log.Err(err).Interface("kubeCtxNs", kubeCtxNs).Msg("UpdateStatefulSet: error")
		return err
	}

	return nil
}
