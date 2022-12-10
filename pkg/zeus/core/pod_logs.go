package zeus_core

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/core/v1"
)

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
