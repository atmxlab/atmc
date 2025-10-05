package main

import (
	"fmt"
	"log"

	"github.com/atmxlab/atmcfg"
)

func main() {
	scanner, err := atmcfg.NewATMC().Load("./cmd/config.atmc")
	if err != nil {
		log.Fatal(err)
	}

	cfg := make(map[string]any)
	err = scanner.Scan(cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cfg)
}
