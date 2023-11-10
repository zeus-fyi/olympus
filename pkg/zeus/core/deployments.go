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

func (k *K8Util) GetDeploymentList(ctx context.Context, kns zeus_common_types.CloudCtxNs, filter *string_utils.FilterOpts) (*v1.DeploymentList, error) {
	k.SetContext(kns.Context)
	d, err := k.kc.AppsV1().Deployments(kns.Namespace).List(ctx, metav1.ListOptions{})
	return d, err
}

func (k *K8Util) GetDeployment(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) (*v1.Deployment, error) {
	k.SetContext(kns.Context)
	d, err := k.kc.AppsV1().Deployments(kns.Namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Err(err).Interface("kns", kns).Msg("GetDeployment: error")
		return nil, err
	}
	return d, err
}

func (k *K8Util) CreateDeployment(ctx context.Context, kns zeus_common_types.CloudCtxNs, d *v1.Deployment, filter *string_utils.FilterOpts) (*v1.Deployment, error) {
	k.SetContext(kns.Context)
	opts := metav1.CreateOptions{}
	d, err := k.kc.AppsV1().Deployments(kns.Namespace).Create(ctx, d, opts)
	alreadyExists := errors.IsAlreadyExists(err)
	if alreadyExists {
		log.Err(err).Interface("kns", kns).Msg("Deployment already exists, skipping creation")
		return d, nil
	}
	return d, err
}

func (k *K8Util) DeleteDeployment(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) error {
	k.SetContext(kns.Context)
	opts := metav1.DeleteOptions{}
	err := k.kc.AppsV1().Deployments(kns.Namespace).Delete(ctx, name, opts)
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (k *K8Util) CreateDeploymentIfVersionLabelChangesOrDoesNotExist(ctx context.Context, kns zeus_common_types.CloudCtxNs, nd *v1.Deployment, filter *string_utils.FilterOpts) (*v1.Deployment, error) {
	k.SetContext(kns.Context)

	cd, err := k.GetDeployment(ctx, kns, nd.Name, filter)
	switch {
	case cd != nil && len(cd.Name) > 0:
		switch IsVersionNew(cd.Labels, nd.Labels) {
		case true:
			derr := k.DeleteDeployment(ctx, kns, cd.Name, filter)
			if derr != nil {
				return cd, derr
			}
		case false:
			return cd, nil
		}
	case errors.IsNotFound(err):
		newD, newDErr := k.CreateDeployment(ctx, kns, nd, filter)
		return newD, newDErr
	}
	newD, newDErr := k.CreateDeployment(ctx, kns, nd, filter)
	return newD, newDErr
}

func (k *K8Util) RolloutRestartDeployment(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) (*v1.Deployment, error) {
	k.SetContext(kns.Context)
	d, err := k.kc.AppsV1().Deployments(kns.Namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Err(err).Interface("kns", kns).Msg("GetDeployment: error")
		return nil, err
	}
	// Prepare for the restart
	if d.Spec.Template.Annotations == nil {
		d.Spec.Template.Annotations = make(map[string]string)
	}

	// Set a new annotation - this triggers a restart
	d.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	// Update the deployment
	_, err = k.kc.AppsV1().Deployments(kns.Namespace).Update(ctx, d, metav1.UpdateOptions{})
	if err != nil {
		log.Err(err).Interface("kns", kns).Msg("UpdateDeployment: error")
		return nil, err
	}

	return d, err
}
