package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Client struct {
	C       http.Client
	E       Endpoint
	Headers map[string]string

	EnableBytesStrDecode bool
}

type Reply struct {
	Ctx        context.Context
	Body       string
	StatusCode int
	Status     string
	Err        error
	BodyBytes  []byte
}

func (c Client) Get(ctx context.Context, url string) Reply {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Err(err).Msgf("had error when sending url request: %s", url)
		return Reply{Err: err}
	}
	return c.Do(ctx, req)
}

func (c Client) GetWithPayload(ctx context.Context, url string, payload []byte) Reply {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Err(err).Bytes("payload that failed", payload).Msgf("url: %s", url)
		return Reply{Err: err}
	}
	return c.Do(ctx, req)
}

func (c Client) Post(ctx context.Context, url string, payload []byte) Reply {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Err(err).Bytes("payload that failed", payload).Msgf("url: %s", url)
		return Reply{Err: err}
	}
	return c.Do(ctx, req)
}

func (c Client) Do(ctx context.Context, req *http.Request) Reply {
	if len(c.Headers) > 0 {
		req = c.AppendArbitraryHeaders(ctx, req)
	}
	resp, err := c.C.Do(req)
	if err != nil {
		log.Err(err).Msg("had error when sending url request")
		return Reply{Err: err}
	}
	return c.respParser(ctx, resp)
}

func (c Client) AppendArbitraryHeaders(ctx context.Context, req *http.Request) *http.Request {
	for header, value := range c.Headers {
		c.AddHeader(ctx, req, header, value)
	}
	return req
}

func (c Client) AddHeader(ctx context.Context, req *http.Request, headerKey, headerValue string) *http.Request {
	req.Header.Add(headerKey, headerValue)
	return req
}

func (c Client) respParser(ctx context.Context, resp *http.Response) Reply {
	r := Reply{}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Err(err).Msg("had error closing body reader")
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.Err = err
		log.Err(err).Bytes("body bytes in error", body)
	}
	r.Status = resp.Status
	r.BodyBytes = body
	r.StatusCode = resp.StatusCode
	if r.StatusCode <= 599 && r.StatusCode >= 400 {
		log.Ctx(ctx).Info().Msgf("had 4xx-5xx status code resp from request, status: %d", r.StatusCode)
	}
	if c.EnableBytesStrDecode {
		r.Body = fmt.Sprintf("%s", body)
	}
	return r
}
