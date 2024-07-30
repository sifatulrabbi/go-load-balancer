package configs

import (
	"fmt"
	"testing"
)

func TestConfigParsing(t *testing.T) {
	cfg, err := ParseLoadBalancerConfig("../tests/sample-config.json")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cfg)
}
