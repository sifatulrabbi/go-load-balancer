package main

import (
	"flag"
	"log"

	"github.com/sifatulrabbi/gld/configs"
	"github.com/sifatulrabbi/gld/loadbalancer"
	"github.com/sifatulrabbi/gld/tcp"
)

func main() {
	configFile := flag.String("config", "", "Path to the config file")
	flag.Parse()

	var cfg *configs.LoadBalancerConfig
	if *configFile != "" {
		c, err := configs.ParseLoadBalancerConfig(*configFile)
		if err != nil {
			log.Panicln(err)
		}
		cfg = c
	} else {
		log.Panicln("'-config' is required to run the load balancer.")
	}

	ld := loadbalancer.New("round_robin", cfg.ServerURIs)
	tcpServer := tcp.New(cfg, ld)

	tcpServer.Start()
}
