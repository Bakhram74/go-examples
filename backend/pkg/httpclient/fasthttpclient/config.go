package fasthttpclient

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	defaultMaxConcurrency            = 8192
	defaultMaxConns                  = 1024
	defaultMaxIdleConnDuration       = 30 * time.Second
	defaultMaxConnDuration           = 30 * time.Second
	defaultResponseTimeout           = 30 * time.Second
	defaultMaxIdemponentCallAttempts = 300
	defaultTimeoutDivider            = 2
	defaultMaxConnWaitTimeout        = 500 * time.Millisecond
)

type FastHttpConfig struct {
	Client              *fasthttp.Client
	MaxIdleConnDuration string
	MaxConnDuration     string
	ReadTimeout         string
	WriteTimeout        string
	MaxConnWaitTimeout  string
	ConnPoolStrategy    string
	MaxConcurrency      int
	TlsConfig           *TLSConfig
}

type TLSConfig struct {
	InsecureSkipVerify bool
}

type fastHttpOpts struct {
	RetryIf         fasthttp.RetryIfFunc
	ConfigureClient func(hc *fasthttp.HostClient) error
	TLSConfig       *tls.Config
	Dial            fasthttp.DialFunc
}

type FastHttpOption func(opt *fastHttpOpts)

func defaultDial(maxConcurrency int) fasthttp.DialFunc {
	if maxConcurrency == 0 {
		maxConcurrency = defaultMaxConcurrency
	}
	dial := &fasthttp.TCPDialer{
		Concurrency:      maxConcurrency,
		DNSCacheDuration: fasthttp.DefaultDNSCacheDuration, // 1 min
	}
	return dial.Dial
}

func defaultTLSConfig(tlsConf TLSConfig) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: tlsConf.InsecureSkipVerify, //nolint:gosec
	}
}

func initFastHttpOptions(cfg *FastHttpConfig, opts ...FastHttpOption) *FastHttpConfig {
	options := &fastHttpOpts{}
	for _, opt := range opts {
		opt(options)
	}
	if options.Dial == nil {
		options.Dial = defaultDial(cfg.MaxConcurrency)
	}

	if options.TLSConfig == nil && cfg.TlsConfig != nil {
		options.TLSConfig = defaultTLSConfig(*cfg.TlsConfig)
	}
	cfg.Client.RetryIf = options.RetryIf
	cfg.Client.Dial = options.Dial
	cfg.Client.TLSConfig = options.TLSConfig
	cfg.Client.ConfigureClient = options.ConfigureClient
	return cfg
}

func initFastHttpConfig(cfg *FastHttpConfig) (*FastHttpConfig, error) {
	var err error
	if cfg.Client == nil {
		cfg.Client = &fasthttp.Client{}
	}

	cfg.Client.ConnPoolStrategy = fasthttp.FIFO

	if len(cfg.ConnPoolStrategy) != 0 {
		switch cfg.ConnPoolStrategy {
		case "FIFO":
			cfg.Client.ConnPoolStrategy = fasthttp.FIFO
		case "LIFO":
			cfg.Client.ConnPoolStrategy = fasthttp.LIFO
		}
	}

	cfg.Client.MaxIdleConnDuration = defaultMaxIdleConnDuration
	if len(cfg.MaxIdleConnDuration) != 0 {
		cfg.Client.MaxIdleConnDuration, err = time.ParseDuration(cfg.MaxIdleConnDuration)
		if err != nil {
			return nil, fmt.Errorf("cannot init FastHttpConfig: %w", err)
		}
	}

	cfg.Client.MaxConnDuration = defaultMaxConnDuration
	if len(cfg.MaxIdleConnDuration) != 0 {
		cfg.Client.MaxConnDuration, err = time.ParseDuration(cfg.MaxConnDuration)
		if err != nil {
			return nil, fmt.Errorf("cannot init FastHttpConfig: %w", err)
		}
	}
	cfg.Client.ReadTimeout = defaultResponseTimeout / defaultTimeoutDivider
	if len(cfg.ReadTimeout) != 0 {
		cfg.Client.ReadTimeout, err = time.ParseDuration(cfg.ReadTimeout)
		if err != nil {
			return nil, fmt.Errorf("cannot init FastHttpConfig: %w", err)
		}
	}

	cfg.Client.WriteTimeout = defaultResponseTimeout / defaultTimeoutDivider
	if len(cfg.WriteTimeout) != 0 {
		cfg.Client.WriteTimeout, err = time.ParseDuration(cfg.WriteTimeout)
		if err != nil {
			return nil, fmt.Errorf("cannot init FastHttpConfig: %w", err)
		}
	}

	cfg.Client.MaxConnWaitTimeout = defaultMaxConnWaitTimeout
	if len(cfg.MaxConnWaitTimeout) != 0 {
		cfg.Client.MaxConnWaitTimeout, err = time.ParseDuration(cfg.MaxConnWaitTimeout)
		if err != nil {
			return nil, fmt.Errorf("cannot init FastHttpConfig: %w", err)
		}
	}

	return cfg, nil
}
