package main

import (
	c "github.com/Mongey/terraform-provider-confluentcloud/ccloud"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: c.Provider})
}
