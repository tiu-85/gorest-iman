package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"tiu-85/gorest-iman/pkg/common/infra/adapters"
	"tiu-85/gorest-iman/pkg/common/infra/clients"
	"tiu-85/gorest-iman/pkg/common/infra/values"
	pdc "tiu-85/gorest-iman/pkg/post-fetch/domains/clients"
)

type PostApiClient struct {
	logger adapters.Logger
	cfg    *values.Config
	client clients.HTTPClient
}

func NewPostApiClient(
	logger adapters.Logger,
	cfg *values.Config,
	httpFactory clients.HTTPClientFactory,
) (pdc.PostApiClient, error) {
	return &PostApiClient{
		logger: logger.Named("post_api_client"),
		cfg:    cfg,
		client: httpFactory.CreateClient(),
	}, nil
}

func (c *PostApiClient) doRequest(ctx context.Context, req *pdc.PostApiClientRequest, result interface{}) error {
	logger := c.logger.WithCtx(ctx, "method", "do_request", "url", req.Url, "path", req.Path)

	u, err := url.Parse(fmt.Sprintf("%s%s", req.Url, req.Path))
	if err != nil {
		logger.Error(err)
		return err
	}

	request, err := http.NewRequest(req.Method, u.String(), nil)
	if err != nil {
		logger.Error(err)
		return err
	}

	for key, value := range req.Headers {
		request.Header.Add(key, value)
	}

	response, err := c.client.Do(ctx, request)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = response.Body.Close()
	if err != nil {
		logger.Error(err)
		return err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *PostApiClient) Fetch(ctx context.Context, page uint32) (*pdc.PostFetchResponse, error) {
	logger := c.logger.WithCtx(ctx, "method", "fetch", "page", page)

	request := pdc.PostApiClientRequest{
		Method: http.MethodGet,
		Url:    c.cfg.ExternalApiUrl,
		Path:   fmt.Sprintf("/public/v1/posts?page=%d", page),
	}
	request.Headers = map[string]string{
		"Content-Type": "application/json",
	}

	response := new(pdc.PostFetchResponse)

	err := c.doRequest(ctx, &request, &response)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return response, nil
}
