package grpc_proxyserver

import "time"

type option struct {
	Address             string
	WriteTimeout        time.Duration
	ReadTimeout         time.Duration
	IdleTimeout         time.Duration
	ReadHeaderTimeout   time.Duration
	HealthCheckEndpoint string
}

func (o *option) clone() *option {
	n := *o
	return &n
}

var defaultOption = &option{
	Address:             ":8080",
	WriteTimeout:        time.Second * 5,
	ReadTimeout:         time.Second * 5,
	IdleTimeout:         time.Second * 75,
	ReadHeaderTimeout:   time.Second * 5,
	HealthCheckEndpoint: "/healthz",
}

type Option func(*option)

func Port(port string) Option {
	return func(o *option) {
		o.Address = port
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(o *option) {
		o.WriteTimeout = timeout
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(o *option) {
		o.ReadTimeout = timeout
	}
}

func IdleTimeout(timeout time.Duration) Option {
	return func(o *option) {
		o.IdleTimeout = timeout
	}
}

func ReadHeaderTimeout(timeout time.Duration) Option {
	return func(o *option) {
		o.IdleTimeout = timeout
	}
}

func HealthCheckEndpoint(path string) Option {
	return func(o *option) {
		o.HealthCheckEndpoint = path
	}
}
