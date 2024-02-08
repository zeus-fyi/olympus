package pods

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/iris/resty_base"
	"github.com/zeus-fyi/olympus/pkg/utils/client"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	zeus_pods_reqs "github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types/pods"
	zeus_pods_resp "github.com/zeus-fyi/zeus/zeus/z_client/zeus_resp_types/pods"
)

func podsPortForwardRequestToAllPods(c echo.Context, request *zeus_pods_reqs.PodActionRequest) error {
	ctx := context.Background()
	log.Debug().Msg("start podsPortForwardRequestToAllPods")
	k := zeus.K8Util
	k8CfgInterface := c.Get("k8Cfg")
	if k8CfgInterface != nil {
		k8Cfg, ok := k8CfgInterface.(autok8s_core.K8Util) // Ensure the type assertion is correct
		if ok {
			k = k8Cfg
		}
	}

	pods, err := k.GetPodsUsingCtxNs(ctx, request.CloudCtxNs, nil, request.FilterOpts)
	if err != nil {
		return err
	}
	var respBody zeus_pods_resp.ClientResp
	respBody.ReplyBodies = make(map[string][]byte, len(pods.Items))
	for _, pod := range pods.Items {
		request.PodName = pod.GetName()
		bytesResp, reqErr := PodsPortForwardRequest(c, request)
		if reqErr != nil {
			log.Err(reqErr).Msgf("port-forwarded request to pod %s failed", pod.GetName())
			return c.JSON(http.StatusBadRequest, "port-forwarded request failed")
		}
		respBody.ReplyBodies[pod.GetName()] = bytesResp
	}
	return c.JSON(http.StatusOK, respBody)
}

func PodsPortForwardRequest(c echo.Context, request *zeus_pods_reqs.PodActionRequest) ([]byte, error) {
	k := zeus.K8Util
	k8CfgInterface := c.Get("k8Cfg")
	if k8CfgInterface != nil {
		k8Cfg, ok := k8CfgInterface.(autok8s_core.K8Util) // Ensure the type assertion is correct
		if ok {
			k = k8Cfg
		}
	}

	ctx := context.Background()
	log.Debug().Msg("start PodsPortForwardRequest")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	startChan := make(chan struct{}, 1)
	stopChan := make(chan struct{}, 1)

	var emptyBytes []byte
	if request.ClientReq == nil {
		return emptyBytes, errors.New("no client request info provided")
	}
	clientReq := *request.ClientReq
	go func() {
		log.Debug().Msg("start port-forward thread")
		address := "localhost"
		err := k.PortForwardPod(ctx, request.CloudCtxNs, request.PodName, address, clientReq.Ports, startChan, stopChan, request.FilterOpts)
		log.Err(err).Msg("error in port forwarding")
		log.Debug().Msg("done port-forward")
	}()

	log.Debug().Msg("awaiting signal")
	<-startChan
	log.Debug().Msg("port ready chan ok")
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		close(stopChan)
	}()

	log.Debug().Msg("do port-forwarded commands")
	port := ""
	for _, po := range clientReq.Ports {
		port, _, _ = strings.Cut(po, ":")
	}
	bearer := c.Get("bearer").(string)
	if len(clientReq.EndpointHeaders) > 0 {
		v, ok := clientReq.EndpointHeaders["Authorization"]
		if ok {
			bearer = v
			delete(clientReq.EndpointHeaders, "Authorization")
		}
	}

	restyC := resty_base.GetBaseRestyClient(fmt.Sprintf("http://localhost:%s", port), bearer)
	var r client.Reply
	payload := clientReq.Payload
	switch clientReq.MethodHTTP {
	case http.MethodPost:
		if payload == nil {
			return emptyBytes, errors.New("no payload supplied for POST request")
		}
		resp, err := restyC.R().
			SetHeaders(clientReq.EndpointHeaders).
			SetBody(payload).
			Post(clientReq.Endpoint)
		if err != nil {
			log.Err(err)
		}
		r.Err = err
		r.BodyBytes = resp.Body()
	default:
		restyC.SetAllowGetMethodPayload(true)
		resp, err := restyC.R().
			SetHeaders(clientReq.EndpointHeaders).
			SetBody(payload).
			Get(clientReq.Endpoint)
		if err != nil {
			log.Err(err)
		}
		r.Err = err
		r.BodyBytes = resp.Body()
	}
	close(stopChan)
	log.Debug().Msg("end port-forwarded commands")
	return r.BodyBytes, r.Err
}
