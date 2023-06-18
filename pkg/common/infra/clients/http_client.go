package clients

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"tiu-85/gorest-iman/pkg/common/infra/adapters"
)

type HTTPClient interface {
	Get(ctx context.Context, rawURL string) (*http.Response, error)
	Post(ctx context.Context, rawURL, contentType string, body io.Reader) (*http.Response, error)
	PostForm(ctx context.Context, rawURL string, data url.Values) (*http.Response, error)
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

type httpClient struct {
	ctx    context.Context
	logger adapters.Logger
	cli    *http.Client
}

func NewHTTPClient(ctx context.Context, logger adapters.Logger, cli *http.Client) HTTPClient {
	return &httpClient{
		ctx:    ctx,
		logger: logger,
		cli:    cli,
	}
}

func (h *httpClient) Get(ctx context.Context, rawURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	return h.Do(ctx, req)
}

func (h *httpClient) PostForm(ctx context.Context, rawURL string, data url.Values) (*http.Response, error) {
	return h.Post(ctx, rawURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (h *httpClient) Post(ctx context.Context, rawURL, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", rawURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return h.Do(ctx, req)
}

func (h *httpClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	logger := h.logger.WithCtx(ctx, "method", "do", "url", req.URL)

	if req.Header == nil {
		req.Header = http.Header{}
	}

	requestLogger := logger.WithCtx(
		ctx,
		"dump_type",
		"request",
	)

	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		requestLogger.Warn(err)
	}
	requestLogger.Infof("dump request: %s", string(dump))

	debugDuration := time.Now()

	var resp *http.Response
	var originalBody []byte

	if req.ContentLength != 0 {
		originalBody, err = ioutil.ReadAll(req.Body)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	req.Body = io.NopCloser(bytes.NewBuffer(originalBody))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(originalBody)), nil
	}
	resp, err = h.cli.Do(req.Clone(ctx))
	if resp != nil {
		responseLogger := logger.WithCtx(
			ctx,
			"dump_type",
			"response",
			"status_code",
			resp.StatusCode,
		)

		dumpResp, err := httputil.DumpResponse(resp, true)
		if err != nil {
			responseLogger.Error(err)
		}
		responseLogger.Infof("dump response: %s with duration %s", string(dumpResp), time.Since(debugDuration))
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}
