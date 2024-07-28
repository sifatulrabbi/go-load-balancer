package main

import (
	"flag"
	"fmt"
)

func main() {
	uris := flag.String("server", "", "Enter forwarding server URIs")
	flag.Parse()

	fmt.Println(*uris)
}
