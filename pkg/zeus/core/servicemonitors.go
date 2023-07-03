package zeus_core

import (
	"context"

	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetServiceMonitor(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) (*v1.ServiceMonitor, error) {
	k.SetContext(kns.Context)
	sm, err := k.mc.MonitoringV1().ServiceMonitors(kns.Namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("kns", kns).Str("name", name).Msg("GetServiceMonitor")
		return nil, err
	}
	return sm, err
}

func (k *K8Util) CreateServiceMonitor(ctx context.Context, kns zeus_common_types.CloudCtxNs, sm *v1.ServiceMonitor, filter *string_utils.FilterOpts) (*v1.ServiceMonitor, error) {
	k.SetContext(kns.Context)
	opts := metav1.CreateOptions{}
	sm, err := k.mc.MonitoringV1().ServiceMonitors(kns.Namespace).Create(ctx, sm, opts)
	if errors.IsAlreadyExists(err) {
		log.Err(err).Interface("kns", kns).Msg("ServiceMonitor already exists, skipping creation")
		return sm, nil
	}
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("kns", kns).Msg("CreateServiceMonitor")
		return nil, err
	}
	return sm, err
}

func (k *K8Util) DeleteServiceMonitor(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) error {
	k.SetContext(kns.Context)
	opts := metav1.DeleteOptions{}
	err := k.mc.MonitoringV1().ServiceMonitors(kns.Namespace).Delete(ctx, name, opts)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("kns", kns).Str("name", name).Msg("DeleteServiceMonitor")
		return err
	}
	return nil
}

func (k *K8Util) CreateServiceMonitorIfVersionLabelChangesOrDoesNotExist(ctx context.Context, kns zeus_common_types.CloudCtxNs, sm *v1.ServiceMonitor, filter *string_utils.FilterOpts) (*v1.ServiceMonitor, error) {
	k.SetContext(kns.Context)
	csm, err := k.GetServiceMonitor(ctx, kns, sm.Name, filter)
	switch {
	case csm != nil && len(csm.Name) > 0:
		switch IsVersionNew(csm.Labels, sm.Labels) {
		case true:
			smerr := k.DeleteServiceMonitor(ctx, kns, csm.Name, filter)
			if smerr != nil {
				return csm, smerr
			}
		case false:
			return csm, nil
		}
	case errors.IsNotFound(err):
		newSm, newSmErr := k.CreateServiceMonitor(ctx, kns, sm, filter)
		return newSm, newSmErr
	}
	newSm, newDErr := k.CreateServiceMonitor(ctx, kns, sm, filter)
	return newSm, newDErr
}
