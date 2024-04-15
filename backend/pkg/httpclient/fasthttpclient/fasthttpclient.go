package fasthttpclient

import (
	"fmt"
	"single-window/pkg/httpclient"

	"github.com/valyala/fasthttp"
)

type fastHttpClient struct {
	cli *fasthttp.Client
}

func NewFastHttpClient(cfgs *FastHttpConfig, opts ...FastHttpOption) (httpclient.IHttpClient, error) {
	var err error
	cfgs, err = initFastHttpConfig(cfgs)
	cfgs = initFastHttpOptions(cfgs, opts...)
	if err != nil {
		return nil, err
	}

	return &fastHttpClient{
		cli: cfgs.Client,
	}, nil
}

func WrapFastHttpClient(cli *fasthttp.Client) httpclient.IHttpClient {
	return &fastHttpClient{
		cli: cli,
	}
}

func (fc *fastHttpClient) Get(
	url string,
	body *httpclient.SwClientRequest,
) (
	*httpclient.SwClientResponse, error,
) {
	return fc.doRequest(
		fasthttp.MethodGet, url, body.Body,
		body.Headers,
		body.Params,
	)
}

func (fc *fastHttpClient) Post(
	url string,
	body *httpclient.SwClientRequest,
) (
	*httpclient.SwClientResponse, error,
) {
	return fc.doRequest(
		fasthttp.MethodPost, url, body.Body,
		body.Headers,
		body.Params,
	)
}

func (fc *fastHttpClient) doRequest(
	method, uri string,
	body *httpclient.SwClientRequestBody,
	headers, params map[string]string,
) (
	*httpclient.SwClientResponse, error,
) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	fc.initRequest(req, method, uri, body, headers, params)

	if err := fc.cli.Do(req, resp); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return &httpclient.SwClientResponse{
		Payload: resp.Body(),
		Status:  resp.StatusCode(),
	}, nil
}

func (fc *fastHttpClient) initRequest(
	req *fasthttp.Request,
	method, uri string,
	body *httpclient.SwClientRequestBody,
	headers, params map[string]string,
) {
	req.SetRequestURI(uri)
	req.Header.SetMethod(method)
	if body != nil {
		req.SetBody(body.Payload)
		req.Header.SetContentType(string(body.ContentType))
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	for key, value := range params {
		req.URI().QueryArgs().Add(key, value)
	}
}
