package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/turt2live/terraform-provider-matrix/matrix"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: matrix.Provider,
	})
}
