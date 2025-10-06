package client

import (
	"bufio"
	"github.com/ilam072/wbtech-l2/17/internal/config"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

type Client struct {
	Host   string
	Port   string
	Dialer net.Dialer
}

func New(cfg config.ClientCfg) Client {
	return Client{
		Host:   cfg.Host,
		Port:   cfg.Port,
		Dialer: net.Dialer{Timeout: *cfg.Timeout},
	}
}

func (c *Client) Read(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("[INFO] STDIN closed (Ctrl+D). Closing connection...")
				err := conn.Close()
				if err != nil {
					log.Fatal("[ERROR] could not close connection: ", err)
				}
				return
			}
			log.Printf("[WARN] Error reading from STDIN: %v\n", err)
			continue
		}

		_, err = conn.Write(line)
		if err != nil {
			log.Printf("[ERROR] Failed to send data to server: %v\n", err)
			_ = conn.Close()
			return
		}
	}
}

func (c *Client) Write(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("[INFO] Connection closed by remote host")
				return
			}
			log.Fatalf("[ERROR] Failed to read from server: %v\n", err)
		}

		_, err = os.Stdout.Write(buf[:n])
		if err != nil {
			log.Printf("[WARN] Failed to write to STDOUT: %v\n", err)
		}
	}
}
