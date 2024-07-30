package tcp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/sifatulrabbi/gld/configs"
	"github.com/sifatulrabbi/gld/loadbalancer"
)

type TCPServer struct {
	Port         string
	LoadBalancer *loadbalancer.LoadBalancer
	listener     *net.TCPListener
	addr         *net.TCPAddr
	errors       []error
	mut          *sync.Mutex
}

func New(cfg *configs.LoadBalancerConfig, ld *loadbalancer.LoadBalancer) *TCPServer {
	t := &TCPServer{
		Port:         cfg.OutboundPort,
		LoadBalancer: ld,
		mut:          &sync.Mutex{},
		errors:       []error{},
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", cfg.Network, cfg.OutboundPort))
	if err != nil {
		log.Panicf("Unable to parse TCP address. Error: %s\n", err.Error())
	}
	t.addr = addr

	l, err := net.ListenTCP("tcp", t.addr)
	if err != nil {
		log.Panicf("Unable to stablish TCP connection. Error: %s\n", err.Error())
	}
	t.listener = l

	return t
}

func (t *TCPServer) Start() {
	fmt.Printf("starting tcp server at %q\nload balancer strategy %q\nlog type \"dev\"\ninitial servers list %v\n",
		t.addr.String(), t.LoadBalancer.Name, t.LoadBalancer.ServerURLs)

	for {
		conn, err := t.listener.AcceptTCP()
		if err != nil {
			t.recordErr(err)
			continue
		}
		go t.handleConn(conn)
	}
}

func (t *TCPServer) handleConn(conn *net.TCPConn) {
	fmt.Println("New connection at")
	defer conn.Close()

	req, err := t.buildHTTPReqFromConn(conn)
	if err != nil {
		t.recordErr(err)
		t.writeResp(conn, nil)
		return
	}

	res, err := t.LoadBalancer.ForwardHTTPReq(req)
	if err != nil {
		t.recordErr(err)
		t.writeResp(conn, nil)
		return
	}

	t.writeResp(conn, res)
}

func (t *TCPServer) writeResp(conn *net.TCPConn, res *http.Response) {
	bcontent := []byte(`HTTP/1.1 500 Internal Server Error
Content-Type: application/json

{"message":"Internal Server Error","statusCode":500,"success":false}`)

	if res != nil {
		if b, err := t.encodeHTTPResp(res); err != nil {
			t.recordErr(err)
		} else {
			bcontent = b
		}
	}

	if _, err := conn.Write(bcontent); err != nil {
		fmt.Println("Error while writing to the TCP connection:", err)
	}
}

func (t *TCPServer) buildHTTPReqFromConn(conn *net.TCPConn) (*http.Request, error) {
	buf := make([]byte, 2048)
	if _, err := conn.Read(buf); err != nil {
		return nil, err
	}
	bufReader := bufio.NewReader(bytes.NewReader(buf))
	return http.ReadRequest(bufReader)
}

func (t *TCPServer) encodeHTTPResp(res *http.Response) ([]byte, error) {
	var outRes bytes.Buffer

	// the first line that defines the HTTP response status
	if _, err := outRes.WriteString(fmt.Sprintf("%s %d %s\r\n",
		res.Proto, res.StatusCode, res.Status)); err != nil {
		return nil, err
	}

	// add all the response headers
	for k, values := range res.Header {
		valuesStr := ""
		for i := 0; i < len(values); i++ {
			valuesStr += values[i]
			if i < len(values) {
				valuesStr += ", "
			}
		}
		if _, err := outRes.WriteString(fmt.Sprintf("%s: %s\r\n", k, valuesStr)); err != nil {
			return nil, err
		}
	}

	// mandatory line break
	outRes.WriteString("\r\n")

	// add the body
	if res.Body != nil {
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if _, err = outRes.Write(resBody); err != nil {
			return nil, err
		}
	}

	return outRes.Bytes(), nil
}

func (t *TCPServer) recordErr(e error) {
	t.mut.Lock()
	t.errors = append(t.errors, e)
	fmt.Printf("error in tcp server %s\n", e.Error())
	t.mut.Unlock()
}
