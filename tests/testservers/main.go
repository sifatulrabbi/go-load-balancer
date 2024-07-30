package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	serverGrp := NewServerGrp([]string{"5001", "5002", "5003"})
	serverGrp.StartAll()
}

type ServerGrp struct {
	List map[string]*http.ServeMux
}

func NewServerGrp(ports []string) *ServerGrp {
	sg := &ServerGrp{List: make(map[string]*http.ServeMux)}
	for _, p := range ports {
		sg.List[p] = http.NewServeMux()
	}
	sg.attachDefaultHandler()
	sg.attachHealthCheckRoutes()
	return sg
}

func (s *ServerGrp) AddNewhHTTPServer(port string) {
	if _, ok := s.List[port]; ok {
		return
	}
	s.List[port] = http.NewServeMux()
}

func (sg *ServerGrp) attachDefaultHandler() {
	for k := range sg.List {
		sg.List[k].HandleFunc("/api/*", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			b, _ := json.Marshal(map[string]any{
				"message":    fmt.Sprintf("API %s is up and running", k),
				"statusCode": 200,
				"success":    true,
			})

			w.WriteHeader(200)
			w.Write(b)
		})
	}
}

func (sg *ServerGrp) attachHealthCheckRoutes() {
	for k := range sg.List {
		sg.List[k].HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			b, _ := json.Marshal(map[string]any{
				"message":    fmt.Sprintf("API %s is up and running", k),
				"statusCode": 200,
				"success":    true,
			})

			w.WriteHeader(200)
			w.Write(b)
		})
	}
}

func (sg *ServerGrp) StartAll() {
	wg := &sync.WaitGroup{}
	for port, mux := range sg.List {
		wg.Add(1)

		go func() {
			addr := fmt.Sprintf("127.0.0.1:%s", port)
			fmt.Printf("Starting server at: %s\n", addr)
			if err := http.ListenAndServe(addr, mux); err != nil {
				log.Panicln(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
