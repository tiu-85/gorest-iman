package clients

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"go.uber.org/multierr"

	"github.com/rs/dnscache"

	"tiu-85/gorest-iman/pkg/common/infra/adapters"
)

const (
	dialContextTimeout   = 60 * time.Second
	dialContextKeepAlive = 60 * time.Second
	httpClientTimeout    = 15 * time.Second
	tlsHandshakeTimeout  = 15 * time.Second
)

type DialContextFunc func(ctx context.Context, network string, addr string) (conn net.Conn, err error)

type HTTPClientFactory interface {
	CreateClient() HTTPClient
}

type httpClientFactory struct {
	ctx    context.Context
	logger adapters.Logger
}

func NewHTTPClientFactory(ctx context.Context, logger adapters.Logger) HTTPClientFactory {
	return &httpClientFactory{
		ctx:    ctx,
		logger: logger,
	}
}

func (f *httpClientFactory) CreateClient() HTTPClient {
	t := &http.Transport{
		DialContext: makeDialContextFunc(),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		TLSHandshakeTimeout: tlsHandshakeTimeout,
		MaxConnsPerHost:     1000,
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 100,
	}
	cli := &http.Client{
		Transport: t,
		Timeout:   httpClientTimeout,
	}
	return NewHTTPClient(f.ctx, f.logger.Named("http_client"), cli)
}

func makeDialContextFunc() func(ctx context.Context, network string, addr string) (conn net.Conn, err error) {
	resolver := dnscache.Resolver{}
	dialer := net.Dialer{
		Timeout:   dialContextTimeout,
		KeepAlive: dialContextKeepAlive,
	}

	return func(ctx context.Context, network string, addr string) (conn net.Conn, err error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}

		ips, err := resolver.LookupHost(ctx, host)
		if err != nil {
			return nil, err
		}

		var connErr error
		for _, ip := range ips {
			conn, connErr = dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
			if connErr != nil {
				err = multierr.Append(err, connErr)
				continue
			}
			return conn, nil
		}
		return nil, err
	}
}
