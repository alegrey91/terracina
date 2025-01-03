// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugin

import (
	"net/rpc"

	"github.com/hashicorp/terracina/internal/terracina"
)

// UIOutput is an implementatin of terracina.UIOutput that communicates
// over RPC.
type UIOutput struct {
	Client *rpc.Client
}

func (o *UIOutput) Output(v string) {
	o.Client.Call("Plugin.Output", v, new(interface{}))
}

// UIOutputServer is the RPC server for serving UIOutput.
type UIOutputServer struct {
	UIOutput terracina.UIOutput
}

func (s *UIOutputServer) Output(
	v string,
	reply *interface{}) error {
	s.UIOutput.Output(v)
	return nil
}
