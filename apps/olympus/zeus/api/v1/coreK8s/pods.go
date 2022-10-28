package coreK8s

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/client"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v12 "github.com/zeus-fyi/olympus/zeus/api/v1"
	v1 "k8s.io/api/core/v1"
)

func HandlePodActionRequest(c echo.Context) error {
	request := new(PodActionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request.Action == "logs" {
		return PodLogsActionRequest(c, request)
	}
	if request.Action == "describe" {
		return PodsDescribeRequest(c, request)
	}
	if request.Action == "describe-audit" {
		return PodsAuditRequest(c, request)
	}
	if request.Action == "delete" {
		return PodsDeleteRequest(c, request)
	}
	if request.Action == "delete-all" {
		return PodsDeleteAllRequest(c, request)
	}
	if request.Action == "delete-all-delay" {
		time.Sleep(time.Second * 180)
		return PodsDeleteAllRequest(c, request)
	}
	if request.Action == "port-forward" {
		bytesResp, err := podsPortForwardRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, string(bytesResp))
		}
		return c.JSON(http.StatusOK, string(bytesResp))
	}
	if request.Action == "port-forward-all" {
		return podsPortForwardRequestToAllPods(c, request)
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func podsPortForwardRequestToAllPods(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("start podsPortForwardRequestToAllPods")

	pods, err := v12.K8util.GetPodsUsingCtxNs(ctx, request.Kns, nil, request.FilterOpts)
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
		err := v12.K8util.PortForwardPod(ctx, request.Kns, request.PodName, address, clientReq.Ports, startChan, stopChan, request.FilterOpts)
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

func PodsDeleteRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("PodsDeleteRequest")
	err := v12.K8util.DeleteFirstPodLike(ctx, request.Kns, request.PodName, request.DeleteOpts, request.FilterOpts)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("pod %s deleted", request.PodName))
}

func PodsDeleteAllRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("PodsDeleteAllRequest")
	err := v12.K8util.DeleteAllPodsLike(ctx, request.Kns, request.PodName, request.DeleteOpts, request.FilterOpts)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("pods with name like %s deleted", request.PodName))
}

func PodsDescribeRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	pods, err := v12.K8util.GetPodsUsingCtxNs(ctx, request.Kns, request.LogOpts, request.FilterOpts)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, pods)
}

func PodLogsActionRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("PodLogsActionRequest")
	pods, err := v12.K8util.GetPodsUsingCtxNs(ctx, request.Kns, nil, request.FilterOpts)
	if err != nil {
		return err
	}

	p := v1.Pod{}
	for _, pod := range pods.Items {
		if string_utils.FilterStringWithOpts(pod.GetName(), request.FilterOpts) {
			p = pod
		}
	}
	logs, err := v12.K8util.GetPodLogs(ctx, p.GetName(), request.Kns.Namespace, request.LogOpts, request.FilterOpts)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, string(logs))
}

func PodsAuditRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()

	pods, err := v12.K8util.GetPodsUsingCtxNs(ctx, request.Kns, request.LogOpts, request.FilterOpts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	parsedResp := parseResp(pods)
	return c.JSON(http.StatusOK, parsedResp)
}
