package main

/* Bootstrap the plugin for Pass */

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/rafaelfelix/terraform-provider-pass/pass"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pass.Provider,
	})
}
