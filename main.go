package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	loadblancer "github.com/sifatulrabbi/gld/loadbalancer"
)

type LoadBalancerConfig struct {
	Network        string   `json:"network"`
	OutboundPort   string   `json:"outboundPort"`
	ServerURIs     []string `json:"serverURIs"`
	Strategy       string   `json:"strategy"`
	MaxServerCount int      `json:"maxServerCount"`
}

func main() {
	configFile := flag.String("config", "", "Path to the config file")
	flag.Parse()

	var cfg *LoadBalancerConfig
	if *configFile != "" {
		c, err := parseConfigFile(*configFile)
		if err != nil {
			log.Panicln(err)
		}
		cfg = c
	} else {
		log.Panicln("'-config' is required to run the load balancer.")
	}

	ld := loadblancer.New("round_robin", cfg.ServerURIs)

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", cfg.Network, cfg.OutboundPort))
	if err != nil {
		log.Panicln(err)
	}
	tcpListener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Panicln(err)
	}

	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleNewConn(ld, conn)
	}
}

func parseConfigFile(filepath string) (*LoadBalancerConfig, error) {
	if filepath == "" {
		return nil, errors.New("Please enter the config file path")
	}

	b, err := os.ReadFile(filepath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("'%s' does not exists\n", filepath)
	} else if err != nil {
		return nil, fmt.Errorf("Unable to read %q\n", filepath)
	}

	cfg := &LoadBalancerConfig{}
	if err = json.Unmarshal(b, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func handleNewConn(ld *loadblancer.LoadBalancer, conn *net.TCPConn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	if _, err := conn.Read(buf); err != nil {
		fmt.Println("tcp read error:", err)
	}

	inReq, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buf)))
	if err != nil {
		fmt.Println("error while parsing incoming http request. error:", err)
	}

	if inReq.ProtoMajor != 1 && inReq.ProtoMinor != 1 {
		log.Fatalf("Incompatible HTTP version. Supported version=1.1, got=%d.%d\n",
			inReq.ProtoMajor, inReq.ProtoMinor)
	}
}

// var outRes bytes.Buffer
// // the first line that defines the HTTP response status
// if _, err = outRes.WriteString(fmt.Sprintf("HTTP/%d.%d %d %s\r\n",
// 	res.ProtoMajor, res.ProtoMinor, res.StatusCode, res.Status)); err != nil {
// 	log.Panicln(err)
// }
// // add all the response headers
// for k, values := range res.Header {
// 	valuesStr := ""
// 	for i := 0; i < len(values); i++ {
// 		valuesStr += values[i]
// 		if i < len(values) {
// 			valuesStr += ", "
// 		}
// 	}
// 	if _, err := outRes.WriteString(fmt.Sprintf("%s: %s\r\n", k, valuesStr)); err != nil {
// 		log.Panicln(err)
// 	}
// }
// // mandatory line break
// outRes.WriteString("\r\n")
// // add the body
// if res.Body != nil {
// 	resBody, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	if _, err = outRes.Write(resBody); err != nil {
// 		log.Panicln(err)
// 	}
// }
//
// fmt.Println("response:\n", outRes.String())
