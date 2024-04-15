package webapi

import (
	"fmt"
	"net/http"
	"single-window/config"
	"single-window/internal/entity"
	"single-window/pkg/httpclient"
	"single-window/pkg/httpclient/fasthttpclient"

	jsoniter "github.com/json-iterator/go"
)

type swAuthClient struct {
	client httpclient.IHttpClient
	url    string
}

func NewWbAuthClient(cfg config.AuthClient) (*swAuthClient, error) {
	client, err := fasthttpclient.NewFastHttpClient(&cfg.FastHttpConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot init fasthttpclient: %w", err)
	}
	return &swAuthClient{client: client, url: cfg.WBURL}, nil
}

func (c *swAuthClient) CheckAuthCode(
	body entity.AuthUserJSONBody,
	headers entity.AuthUserParams,
) (*entity.WbxAuthCodeCheckResponse, error) {
	rawBody, err := jsoniter.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal WBAuthV2RequestBody, err: %w", err)
	}

	resp, err := c.client.Post(c.url, httpclient.NewRequest(
		httpclient.WithJsonRequestBody(rawBody),
		httpclient.WithRequestHeaders(headers.ConvertHeadersToStringMap()),
	))
	if err != nil {
		return nil, fmt.Errorf("cannot send auth data to Authv3, err: %w", err)
	}

	if resp.Status != http.StatusOK {
		return nil, fmt.Errorf("wrong response status code")
	}

	var res entity.WbxAuthCodeCheckResponse
	if err := jsoniter.Unmarshal(resp.Payload, &res); err != nil {
		return nil, fmt.Errorf("cannot unmarshal response to WbxAuthCodeCheckResponse, err: %w", err)
	}
	return &res, nil
}
