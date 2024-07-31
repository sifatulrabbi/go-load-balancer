package loadbalancer

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	ROUNDED_ROBIN          = "rounded_robin"
	WEIGHTED_ROUNDED_ROBIN = "weighted_rounded_robin"
)

type LoadBalancer struct {
	Name        string
	ServerList  map[int]ServerEntry
	serverCount int
	currIdx     int
	nextIdx     int
}

type ServerEntry struct {
	ID         int
	Url        string
	Healthy    bool
	reqHandled int
	reqFailed  int
}

type ServerList map[int]ServerEntry

type strategyFn func(ld *LoadBalancer) *ServerEntry

func New(strategyName string, serverURLs []string) *LoadBalancer {
	ld := LoadBalancer{
		Name:        "round_robin",
		serverCount: len(serverURLs),
		ServerList:  map[int]ServerEntry{},
	}
	for i, v := range serverURLs {
		ld.ServerList[i] = ServerEntry{i, v, true, 0, 0}
	}

	ld.currIdx = 0
	ld.nextIdx = 0

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

	res, err := http.DefaultClient.Do(fwdReq)
	if err != nil {
		ld.updateServerStatus(s.ID, false)
		return nil, err
	}
	ld.updateServerStatus(s.ID, true)
	return res, nil
}

func (ld *LoadBalancer) chooseServer() *ServerEntry {
	switch ld.Name {
	case WEIGHTED_ROUNDED_ROBIN:
		return weightedRoundRobinStrategy(ld)
	default:
		return roundRobinStrategy(ld)
	}
}

func (ld *LoadBalancer) periodicHealthCheck() {
	for {
		for i, s := range ld.ServerList {
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/health", s.Url), http.NoBody)
			res, err := http.DefaultClient.Do(req)
			if err != nil || res.StatusCode != http.StatusOK {
				ld.ServerList[i] = ServerEntry{
					s.ID,
					s.Url,
					false,
					s.reqHandled,
					s.reqFailed,
				}
				fmt.Printf("server unhealthy: %q\n", s.Url)
			} else {
				ld.ServerList[i] = ServerEntry{
					s.ID,
					s.Url,
					true,
					s.reqHandled,
					s.reqFailed,
				}
			}
		}

		time.Sleep(time.Minute * 1)
	}
}

func (ld *LoadBalancer) updateServerStatus(id int, success bool) {
	s := ld.ServerList[id]
	newEntry := ServerEntry{
		ID:         id,
		Url:        s.Url,
		Healthy:    s.Healthy,
		reqHandled: s.reqHandled,
		reqFailed:  s.reqFailed,
	}
	if success {
		newEntry.reqHandled++
	} else {
		newEntry.reqFailed++
	}
	ld.ServerList[id] = newEntry
}
