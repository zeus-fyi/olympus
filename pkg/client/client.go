package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	C                    http.Client
	EnableBytesStrDecode bool
}

type Reply struct {
	Ctx       context.Context
	Body      string
	Status    string
	Err       error
	BodyBytes []byte
}

func (c Client) AddHeader(req *http.Request, header, value string) *http.Request {
	req.Header.Add(header, value)
	return req
}

func (c Client) Get(ctx context.Context, url string) Reply {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return Reply{Err: err}
	}
	return c.Do(ctx, req)
}

func (c Client) GetWithPayload(ctx context.Context, url string, payload []byte) Reply {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("%v", err)
		return Reply{Err: err}
	}
	return c.Do(ctx, req)
}

func (c Client) Post(ctx context.Context, url string, payload []byte) Reply {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("%v", err)
		return Reply{Err: err}
	}
	return c.Do(ctx, req)
}

func (c Client) Do(ctx context.Context, req *http.Request) Reply {
	resp, err := c.C.Do(req)
	if err != nil {
		fmt.Printf("%v", err)
		return Reply{Err: err}
	}
	return c.respParser(resp)
}

func (c Client) respParser(resp *http.Response) Reply {
	r := Reply{}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("%v", err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.Err = err
		fmt.Printf("%v", err)
	}
	r.Status = resp.Status
	r.BodyBytes = body
	if c.EnableBytesStrDecode {
		r.Body = fmt.Sprintf("%s", body)
	}
	return r
}
