package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
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
		log.Panicln("'-server' or '-config' is required to run the load balancer.")
	}

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
		go handleNewConn(conn)
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

func handleNewConn(conn *net.TCPConn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("tcp read error:", err)
			break
		}
	}

	fmt.Println(string(buf))
}
