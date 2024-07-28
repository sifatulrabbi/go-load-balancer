package tcp

import (
	"net"

	"github.com/sifatulrabbi/gld/configs"
)

type TCPInterceptor struct {
	listener *net.TCPListener
	errors   []TCPError
	handlers map[string]func() error
}

func New(cfg *configs.NetworkConfig) *TCPInterceptor {
	if cfg == nil {
		cfg = configs.DefaultNeworkConfig()
	}

	addr, err := net.ResolveTCPAddr(cfg.Ip, cfg.Port)
	l, err := net.ListenTCP(cfg.Ip, addr)

	ti := &TCPInterceptor{listener: l}
	return ti
}

func (ti *TCPInterceptor) Start() {
	for {
		conn, err := ti.listener.AcceptTCP()
	}
}

func (ti *TCPInterceptor) Errors() []TCPError {
	return ti.errors
}

func (ti *TCPInterceptor) RegisterHandler(k string, fn func() error) {
	ti.handlers[k] = fn
}
