// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"github.com/hashicorp/terracina/internal/grpcwrap"
	plugin "github.com/hashicorp/terracina/internal/plugin6"
	simple "github.com/hashicorp/terracina/internal/provider-simple-v6"
	"github.com/hashicorp/terracina/internal/tfplugin6"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		GRPCProviderFunc: func() tfplugin6.ProviderServer {
			return grpcwrap.Provider6(simple.Provider())
		},
	})
}
