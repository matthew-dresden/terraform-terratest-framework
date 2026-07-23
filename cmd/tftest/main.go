package main

import (
	"github.com/matthew-dresden/terraform-terratest-framework/cmd/tftest/cmd"
)

func main() {
	// Version is injected into the cmd package at build time via -ldflags; see Makefile.
	cmd.Execute()
}
