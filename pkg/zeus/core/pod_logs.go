package zeus_core

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
)

func (k *K8Util) GetPodLogs(ctx context.Context, name string, kns zeus_common_types.CloudCtxNs, logOpts *v1.PodLogOptions, filter *string_utils.FilterOpts) ([]byte, error) {
	k.SetContext(kns.Context)
	log.Ctx(ctx).Debug().Msg("GetPodLogs")
	if logOpts == nil {
		logOpts = &v1.PodLogOptions{}
	}
	req := k.kc.CoreV1().Pods(kns.Namespace).GetLogs(name, logOpts)
	if req == nil {
		return nil, fmt.Errorf("GetPodLogs: req is nil")
	}
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
