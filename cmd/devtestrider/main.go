package main

import (
	"log"

	"github.com/ismailtsdln/DevTestrider/cmd/devtestrider/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
