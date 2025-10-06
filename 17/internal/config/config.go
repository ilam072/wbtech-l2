package config

import (
	flag "github.com/spf13/pflag"
	"log"
	"time"
)

type Config struct {
	Client ClientCfg
}

type ClientCfg struct {
	Host    string
	Port    string
	Timeout *time.Duration
}

func New() Config {
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout for TCP-server connection")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatalln("usage: go run main.go [--timeout=10s] host port")
	}

	clientCfg := ClientCfg{
		Host:    args[0],
		Port:    args[1],
		Timeout: timeout,
	}

	return Config{clientCfg}
}
