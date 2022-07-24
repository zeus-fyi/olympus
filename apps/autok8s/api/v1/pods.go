package v1

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodActionRequest struct {
	K8sRequest
	Action        string
	PodName       string
	ContainerName string

	Ports           *[]string
	Endpoint        *string
	EndpointHeaders *map[string]string
	LogOpts         *v1.PodLogOptions
	DeleteOpts      *metav1.DeleteOptions
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
		return podsPortForwardRequest(c, request)
	}

	return c.JSON(http.StatusOK, nil)
}

func podsPortForwardRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("start podsPortForwardRequest")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	startChan := make(chan struct{}, 1)
	stopChan := make(chan struct{}, 1)

	var ports []string
	if request.Ports != nil {
		ports = *request.Ports
	}
	go func() {
		log.Ctx(ctx).Debug().Msg("start port-forward thread")
		address := "localhost"
		err := K8util.PortForwardPod(ctx, request.Kns, request.PodName, address, ports, startChan, stopChan)
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
	if request.Endpoint != nil && request.Ports != nil {
		endpoint := *request.Endpoint
		cli := client.Client{}
		port := ""
		for _, po := range ports {
			port, _, _ = strings.Cut(po, ":")
		}
		cli.E = client.Endpoint(fmt.Sprintf("http://localhost:%s", port))

		if request.EndpointHeaders != nil {
			cli.Headers = *request.EndpointHeaders
		}

		r := cli.Get(ctx, string(cli.E)+"/"+endpoint)
		close(stopChan)
		log.Ctx(ctx).Debug().Msg("end port-forwarded commands")
		return c.JSON(http.StatusOK, string(r.BodyBytes))
	}
	return c.JSON(http.StatusBadRequest, "")
}

func podsDeleteRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("podsDeleteRequest")
	err := K8util.DeleteFirstPodLike(ctx, request.Kns, request.PodName, request.DeleteOpts)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("pod %s deleted", request.PodName))
}

func podsDeleteAllRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("podsDeleteAllRequest")
	err := K8util.DeleteAllPodsLike(ctx, request.Kns, request.PodName, request.DeleteOpts)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("pods with name like %s deleted", request.PodName))
}

func podsDescribeRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()

	pods, err := K8util.GetPodsUsingCtxNs(ctx, request.Kns, nil)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pods)
}

func podLogsActionRequest(c echo.Context, request *PodActionRequest) error {
	ctx := context.Background()
	log.Ctx(ctx).Debug().Msg("podLogsActionRequest")
	pods, err := K8util.GetPodsUsingCtxNs(ctx, request.Kns, nil)
	if err != nil {
		return err
	}
	p := v1.Pod{}
	for _, pod := range pods.Items {
		name := pod.ObjectMeta.Name
		if strings.Contains(name, request.PodName) {
			p = pod
		}
	}
	logs, err := K8util.GetPodLogs(ctx, p.GetName(), request.Kns.Namespace, request.LogOpts)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, string(logs))
}
