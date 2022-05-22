package main

import (
	"github.com/vothanhdien/go-graceful-shutdown/cmd"
)

func main() {
	service := cmd.NewService()
	service.Start()
}
