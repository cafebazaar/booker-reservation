package main

import (
	"github.com/cafebazaar/booker-reservation/cmd"
)

func main() {
	// rand.Seed(time.Now().UTC().UnixNano())
	cmd.Execute()
}
