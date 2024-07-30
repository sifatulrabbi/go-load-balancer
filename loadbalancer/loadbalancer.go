package loadbalancer

import (
	"fmt"
	"net/http"
)

type LoadBalancer struct {
	Name        string
	ServerURLs  map[int]string
	serverCount int
	currIdx     int
	nextIdx     int
}

func New(strategy string, serverURLs []string) *LoadBalancer {
	ld := LoadBalancer{
		Name:        "round_robin",
		serverCount: len(serverURLs),
		ServerURLs:  make(map[int]string),
	}
	for i, v := range serverURLs {
		ld.ServerURLs[i] = v
	}

	ld.currIdx = 0
	ld.nextIdx = ld.currIdx + 1

	return &ld
}

func (ld *LoadBalancer) ForwardHTTPReq(req *http.Request) (*http.Response, error) {
	forwardingUrl := fmt.Sprintf("%s%s", ld.chooseServer(), req.RequestURI)
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

func (ld *LoadBalancer) chooseServer() string {
	ld.currIdx = ld.nextIdx
	ld.nextIdx = ld.currIdx + 1
	if ld.nextIdx >= ld.serverCount {
		ld.nextIdx = 0
	}

	return ld.ServerURLs[ld.currIdx]
}
