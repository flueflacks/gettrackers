package main

import (
	"os"

	"github.com/flueflacks/gettrackers/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
