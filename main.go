package main

import (
	"github.com/cloudkarafka/terraform-provider-cloudkarafka/cloudkarafka"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return cloudkarafka.Provider()
		},
	})
}
