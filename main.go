package main

import (
	"context"
	"terraform-provider-cloudkarafka/cloudkarafka"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name cloudkarafka

func main() {
	providerserver.Serve(context.Background(), cloudkarafka.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/cloudkarafka/cloudkarafka",
	})
}
