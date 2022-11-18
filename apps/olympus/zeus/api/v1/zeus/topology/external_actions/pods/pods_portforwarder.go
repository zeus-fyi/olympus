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
	"github.com/zeus-fyi/olympus/pkg/utils/client"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func podsPortForwardRequestToAllPods(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("start podsPortForwardRequestToAllPods")

	pods, err := zeus.K8Util.GetPodsUsingCtxNs(ctx, request.CloudCtxNs, nil, request.FilterOpts)
	if err != nil {
		return err
	}
	var respBody ClientResp
	respBody.ReplyBodies = make(map[string]string, len(pods.Items))

	for _, pod := range pods.Items {
		request.PodName = pod.GetName()
		bytesResp, reqErr := podsPortForwardRequest(request)
		if reqErr != nil {
			log.Err(reqErr).Msgf("port-forwarded request to pod %s failed", pod.GetName())
			return c.JSON(http.StatusBadRequest, "port-forwarded request failed")
		}
		respBody.ReplyBodies[pod.GetName()] = string(bytesResp)

	}
	return c.JSON(http.StatusOK, respBody)
}

func podsPortForwardRequest(request *PodActionRequest) ([]byte, error) {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("start podsPortForwardRequest")

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
		log.Ctx(ctx).Debug().Msg("start port-forward thread")
		address := "localhost"
		err := zeus.K8Util.PortForwardPod(ctx, request.CloudCtxNs, request.PodName, address, clientReq.Ports, startChan, stopChan, request.FilterOpts)
		log.Ctx(ctx).Err(err).Msg("error in port forwarding")
		log.Ctx(ctx).Debug().Msg("done port-forward")
	}()

	log.Ctx(ctx).Debug().Msg("awaiting signal")
	<-startChan
	log.Ctx(ctx).Debug().Msg("port ready chan ok")
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		close(stopChan)
	}()

	log.Ctx(ctx).Debug().Msg("do port-forwarded commands")
	cli := client.Client{}
	port := ""
	for _, po := range clientReq.Ports {
		port, _, _ = strings.Cut(po, ":")
	}
	cli.E = client.Endpoint(fmt.Sprintf("http://localhost:%s", port))
	cli.Headers = clientReq.EndpointHeaders

	var r client.Reply
	payloadBytes := clientReq.PayloadBytes
	payload := clientReq.Payload
	var finalPayload []byte

	// prefer bytes, but use string if exists
	if payloadBytes != nil {
		finalPayload = *payloadBytes
	} else if payload != nil {
		finalPayload = []byte(*payload)
	}

	switch clientReq.MethodHTTP {
	case http.MethodPost:
		if finalPayload == nil {
			return emptyBytes, errors.New("no payload supplied for POST request")
		}
		r = cli.Post(ctx, string(cli.E)+"/"+clientReq.Endpoint, finalPayload)
	default:
		if finalPayload != nil {
			r = cli.GetWithPayload(ctx, string(cli.E)+"/"+clientReq.Endpoint, finalPayload)
		} else {
			r = cli.Get(ctx, string(cli.E)+"/"+clientReq.Endpoint)
		}
	}
	close(stopChan)
	log.Ctx(ctx).Debug().Msg("end port-forwarded commands")
	return r.BodyBytes, nil
}
