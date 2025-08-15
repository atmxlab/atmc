package main

import (
	"fmt"
	"log"

	"github.com/atmxlab/atmc"
)

type Config struct {
	Logging struct {
		Enabled        bool     `json:"enabled"`
		EnabledTracing bool     `json:"enabled_tracing"`
		Level          []string `json:"level"`
	} `json:"logging"`
	Outbox struct {
		Enabled     bool `json:"enabled"`
		WorkerCount int  `json:"worker_count"`
	} `json:"outbox"`
}

func main() {
	path := "./cmd/config.atmc"
	scanner, err := atmc.New(atmc.WithFieldTag("json")).Load(path)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err = scanner.Scan(&cfg); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
}
