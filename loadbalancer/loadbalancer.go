package loadbalancer

import (
	"github.com/sifatulrabbi/gld/configs"
	"github.com/sifatulrabbi/gld/tcp"
)

type LoadBalancer struct {
	TCP *tcp.TCPInterceptor
}

func New(cfg *configs.LoadBalancerConfig, ti *tcp.TCPInterceptor) *LoadBalancer {
	return &LoadBalancer{ti}
}

func (l *LoadBalancer) Start() {
	go l.TCP.Start()
}
