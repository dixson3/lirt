package main

import (
	"fmt"
	"os"
)

var version = "dev"

func main() {
	fmt.Printf("lirt version %s\n", version)
	os.Exit(0)
}
