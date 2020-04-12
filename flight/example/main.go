package main

import (
	"time"

	"github.com/353solutions/carrow/flight"
)

func main() {
	flight.Start()
	time.Sleep(20 * time.Second)
}
