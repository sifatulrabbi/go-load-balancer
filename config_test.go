package main

import (
	"fmt"
	"testing"
)

func TestConfigParsing(t *testing.T) {
	cfg, err := parseConfigFile("./tests/sample-config.json")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cfg)
}
