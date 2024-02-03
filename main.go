package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/winebarrel/terraform-provider-lambdazip/lambdazip"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name lambdazip

func main() {
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: lambdazip.Provider,
		ProviderAddr: "registry.terraform.io/winebarrel/lambdazip",
		Debug:        *debug,
	}

	plugin.Serve(opts)
}
