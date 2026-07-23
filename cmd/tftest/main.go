package main

import (
	"github.com/matthew-dresden/terraform-terratest-framework/cmd/tftest/cmd"
)

// Version will be set during build
var Version = "dev"

func main() {
	// Set version
	cmd.Version = Version

	// Execute the root command
	cmd.Execute()
}
