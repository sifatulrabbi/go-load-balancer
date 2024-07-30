package loadbalancer

import (
	"fmt"
	"net/http"
	"time"
)

type ServerEntry struct {
	Url        string
	Healthy    bool
	reqHandled int
	reqFailed  int
}

type ServerList map[int]ServerEntry

type Strategy interface {
	ChooseServer() *ServerEntry
}

func NewStrategy(name string, serverList []string) (Strategy, error) {
	sl := make(map[int]ServerEntry)
	for i, s := range serverList {
		sl[i] = ServerEntry{s, false, 0, 0}
	}

	var strategy Strategy = nil
	switch name {
	default:
		strategy = newRoundRobinStrategy(sl)
		break
	}

	return strategy, nil
}

type RoundRobinStrategy struct {
	Name        string
	ServerList  ServerList
	currIdx     int
	nextIdx     int
	serverCount int
}

func (rs *RoundRobinStrategy) ChooseServer() *ServerEntry {
	rs.currIdx = rs.nextIdx
	rs.nextIdx = rs.currIdx + 1
	if rs.nextIdx >= rs.serverCount {
		rs.nextIdx = 0
	}

	currIdx := rs.currIdx
	server := rs.ServerList[rs.currIdx]
	for !server.Healthy {
		server = *rs.ChooseServer()
		if rs.currIdx == currIdx {
			return nil
		}
	}

	return &server
}

func newRoundRobinStrategy(serverList ServerList) *RoundRobinStrategy {
	s := &RoundRobinStrategy{
		Name:        "round_robin",
		ServerList:  serverList,
		currIdx:     0,
		nextIdx:     1,
		serverCount: len(serverList),
	}

	periodicHealthCheck(&s.ServerList)

	return s
}

type WeightedRoundRobin struct {
	Name string
}

func (s *WeightedRoundRobin) ChooseServer() *ServerEntry {
	return nil
}

// TODO:
func newWeightedRoundRobin() *WeightedRoundRobin {
	s := &WeightedRoundRobin{}
	return s
}

func periodicHealthCheck(sl *ServerList) {
	go func() {
		for {
			for i, s := range *sl {
				req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/health", s.Url), http.NoBody)
				res, err := http.DefaultClient.Do(req)
				if err != nil || res.StatusCode != http.StatusOK {
					(*sl)[i] = ServerEntry{s.Url, false, 0, 0}
					fmt.Printf("server unhealthy: %q\n", s.Url)
				} else {
					(*sl)[i] = ServerEntry{s.Url, true, 0, 0}
				}
			}

			time.Sleep(time.Minute * 1)
		}
	}()
}
