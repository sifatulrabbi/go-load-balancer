package loadblancer

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
	}
	for i, v := range serverURLs {
		ld.ServerURLs[i] = v
	}
	ld.nextIdx = 0
	return &ld
}

func (ld *LoadBalancer) fwdHTTPReq(req *http.Request) (*http.Response, error) {
	forwardingUrl := fmt.Sprintf("%s%s", ld.chooseServer(), req.RequestURI)
	newReq, err := http.NewRequest(req.Method, forwardingUrl, req.Body)
	if err != nil {
		return nil, err
	}

	for k, values := range req.Header {
		for _, v := range values {
			newReq.Header.Set(k, v)
		}
	}

	for _, c := range req.Cookies() {
		newReq.AddCookie(c)
	}

	res, err := http.DefaultClient.Do(newReq)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ld *LoadBalancer) chooseServer() string {
	if ld.nextIdx+1 >= ld.serverCount {
		ld.currIdx = 0
	}
	ld.nextIdx = ld.currIdx + 1

	return ld.ServerURLs[ld.currIdx]
}
