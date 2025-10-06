package main

import (
	"github.com/ilam072/wbtech-l2/17/internal/client"
	"github.com/ilam072/wbtech-l2/17/internal/config"
	"log"
	"net"
	"sync"
)

func main() {
	cfg := config.New()
	cl := client.New(cfg.Client)
	address := net.JoinHostPort(cl.Host, cl.Port)
	log.Printf("[INFO] Connecting to %s (timeout=%s)\n", address, cfg.Client.Timeout.String())

	conn, err := cl.Dialer.Dial("tcp", address)
	if err != nil {
		log.Fatalf("[ERROR] Failed to connect: %v\n", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal("[ERROR] could not close connection: ", err)
		}
	}()
	log.Printf("[INFO] Successfully connected to %s\n", address)

	wg := new(sync.WaitGroup)

	wg.Add(2)
	go cl.Read(wg, conn)
	go cl.Write(wg, conn)

	wg.Wait()
}
