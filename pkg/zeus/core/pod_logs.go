package zeus_core

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
)

func (k *K8Util) GetPodLogs(ctx context.Context, name string, kns zeus_common_types.CloudCtxNs, logOpts *v1.PodLogOptions, filter *strings_filter.FilterOpts) ([]byte, error) {
	k.SetContext(kns.Context)
	log.Debug().Msg("GetPodLogs")
	if logOpts == nil {
		logOpts = &v1.PodLogOptions{}
	}
	req := k.kc.CoreV1().Pods(kns.Namespace).GetLogs(name, logOpts)
	if req == nil {
		return nil, fmt.Errorf("GetPodLogs: req is nil")
	}
	buf := new(bytes.Buffer)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return buf.Bytes(), err // Return here if an error occurs, no need for defer since podLogs is nil or undefined
	}
	defer func() {
		closeErr := podLogs.Close()
		if closeErr != nil {
			fmt.Printf("%s", closeErr.Error())
		}
	}()
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}
