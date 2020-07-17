package main

import (
	c "github.com/Mongey/terraform-provider-confluentcloud/ccloud"

	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: c.Provider})
}
