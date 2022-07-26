package v1

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
	autok8s_core "github.com/zeus-fyi/olympus/pkg/autok8s/core"
	"github.com/zeus-fyi/olympus/pkg/client"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodActionRequest struct {
	K8sRequest
	Action        string
	PodName       string
	ContainerName string

	FilterOpts *autok8s_core.FilterOpts
	ClientReq  *ClientRequest
	LogOpts    *v1.PodLogOptions
	DeleteOpts *metav1.DeleteOptions
}

type ClientRequest struct {
	MethodHTTP      string
	Endpoint        string
	Ports           []string
	Payload         *string
	PayloadBytes    *[]byte
	EndpointHeaders map[string]string
}

type ClientResp struct {
	ReplyBodies map[string]string
}

func HandlePodActionRequest(c echo.Context) error {
	request := new(PodActionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request.Action == "logs" {
		return podLogsActionRequest(c, request)
	}
	if request.Action == "describe" {
		return podsDescribeRequest(c, request)
	}
	if request.Action == "delete" {
		return podsDeleteRequest(c, request)
	}
	if request.Action == "delete-all" {
		return podsDeleteAllRequest(c, request)
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

	pods, err := K8util.GetPodsUsingCtxNs(ctx, request.Kns, nil, request.FilterOpts)
	if err != nil {
		return err
	}
	var respBody ClientResp
	respBody.ReplyBodies = make(map[string]string, len(pods.Items))

	podNameFilter := request.PodName
	var filterWords []string
	if request.FilterOpts != nil {
		filter := request.FilterOpts
		filterMatches := *filter
		filterWords = filterMatches.DoesNotInclude
	}
	for _, pod := range pods.Items {
		if string_utils.FilterString(pod.GetName(), podNameFilter, filterWords) {
			request.PodName = pod.GetName()
			bytesResp, reqErr := podsPortForwardRequest(request)
			if reqErr != nil {
				log.Err(reqErr).Msgf("port-forwarded request to pod %s failed", pod.GetName())
				return c.JSON(http.StatusBadRequest, "port-forwarded request failed")
			}
			respBody.ReplyBodies[pod.GetName()] = string(bytesResp)
		}
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
		err := K8util.PortForwardPod(ctx, request.Kns, request.PodName, address, clientReq.Ports, startChan, stopChan, request.FilterOpts)
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

func podsDeleteRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("podsDeleteRequest")
	err := K8util.DeleteFirstPodLike(ctx, request.Kns, request.PodName, request.DeleteOpts, request.FilterOpts)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("pod %s deleted", request.PodName))
}

func podsDeleteAllRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("podsDeleteAllRequest")
	err := K8util.DeleteAllPodsLike(ctx, request.Kns, request.PodName, request.DeleteOpts, request.FilterOpts)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("pods with name like %s deleted", request.PodName))
}

func podsDescribeRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	pods, err := K8util.GetPodsUsingCtxNs(ctx, request.Kns, nil, request.FilterOpts)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, pods)
}

func podLogsActionRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("podLogsActionRequest")
	pods, err := K8util.GetPodsUsingCtxNs(ctx, request.Kns, nil, request.FilterOpts)
	if err != nil {
		return err
	}

	podNameFilter := request.PodName
	var filterWords []string
	if request.FilterOpts != nil {
		filter := request.FilterOpts
		filterMatches := *filter
		filterWords = filterMatches.DoesNotInclude
	}

	p := v1.Pod{}
	for _, pod := range pods.Items {
		if string_utils.FilterString(pod.GetName(), podNameFilter, filterWords) {
			p = pod
		}
	}
	logs, err := K8util.GetPodLogs(ctx, p.GetName(), request.Kns.Namespace, request.LogOpts, request.FilterOpts)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, string(logs))
}
