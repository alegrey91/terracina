// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"github.com/hashicorp/terracina/internal/builtin/providers/terracina"
	"github.com/hashicorp/terracina/internal/grpcwrap"
	"github.com/hashicorp/terracina/internal/plugin"
	"github.com/hashicorp/terracina/internal/tfplugin5"
)

func main() {
	// Provide a binary version of the internal terracina provider for testing
	plugin.Serve(&plugin.ServeOpts{
		GRPCProviderFunc: func() tfplugin5.ProviderServer {
			return grpcwrap.Provider(terracina.NewProvider())
		},
	})
}
