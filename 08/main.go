package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"time"
)

func main() {
	currentTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get current time from ntp server: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Current time: %v", currentTime.Format(time.RFC1123))
}
