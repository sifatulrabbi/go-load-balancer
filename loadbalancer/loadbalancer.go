package loadbalancer

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ServerEntry struct {
	Url     string
	Healthy bool
}

type LoadBalancer struct {
	Name        string
	ServerList  map[int]ServerEntry
	serverCount int
	currIdx     int
	nextIdx     int
}

func New(strategy string, serverURLs []string) *LoadBalancer {
	ld := LoadBalancer{
		Name:        "round_robin",
		serverCount: len(serverURLs),
		ServerList:  map[int]ServerEntry{},
	}
	for i, v := range serverURLs {
		ld.ServerList[i] = ServerEntry{v, true}
	}

	ld.currIdx = 0
	ld.nextIdx = ld.currIdx + 1

	go ld.periodicHealthCheck()

	return &ld
}

func (ld *LoadBalancer) ForwardHTTPReq(req *http.Request) (*http.Response, error) {
	s := ld.chooseServer()
	if s == nil {
		return nil, errors.New("no healthy server available to serve the client")
	}
	forwardingUrl := fmt.Sprintf("%s%s", s.Url, req.RequestURI)
	fwdReq, err := http.NewRequest(req.Method, forwardingUrl, req.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("forwarding %s to %s\n", req.URL.String(), forwardingUrl)

	// copy over all the headers
	for k, values := range req.Header {
		for _, v := range values {
			fwdReq.Header.Set(k, v)
		}
	}
	// copy over all the cookies
	for _, c := range req.Cookies() {
		fwdReq.AddCookie(c)
	}

	return http.DefaultClient.Do(fwdReq)
}

func (ld *LoadBalancer) chooseServer() *ServerEntry {
	ld.currIdx = ld.nextIdx
	ld.nextIdx = ld.currIdx + 1
	if ld.nextIdx >= ld.serverCount {
		ld.nextIdx = 0
	}

	currIdx := ld.currIdx
	s := ld.ServerList[ld.currIdx]
	for !s.Healthy {
		s = *ld.chooseServer()
		if ld.currIdx == currIdx {
			return nil
		}
	}

	return &s
}

func (ld *LoadBalancer) periodicHealthCheck() {
	for {
		for i, s := range ld.ServerList {
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/health", s.Url), http.NoBody)
			res, err := http.DefaultClient.Do(req)
			if err != nil || res.StatusCode != http.StatusOK {
				ld.ServerList[i] = ServerEntry{s.Url, false}
				fmt.Printf("server unhealthy: %q\n", s.Url)
			} else {
				ld.ServerList[i] = ServerEntry{s.Url, true}
			}
		}

		time.Sleep(time.Minute * 1)
	}
}
