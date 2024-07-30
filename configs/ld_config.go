package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type LoadBalancerConfig struct {
	Network        string   `json:"network"`
	OutboundPort   string   `json:"outboundPort"`
	ServerURIs     []string `json:"serverURIs"`
	Strategy       string   `json:"strategy"`
	MaxServerCount int      `json:"maxServerCount"`
}

func ParseLoadBalancerConfig(filepath string) (*LoadBalancerConfig, error) {
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
