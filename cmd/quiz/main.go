package main

import (
	"fmt"
	"os"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

func mainErr() error {
	return nil
}