// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugin

import (
	"context"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terracina/internal/terracina"
)

// UIInput is an implementation of terracina.UIInput that communicates
// over RPC.
type UIInput struct {
	Client *rpc.Client
}

func (i *UIInput) Input(ctx context.Context, opts *terracina.InputOpts) (string, error) {
	var resp UIInputInputResponse
	err := i.Client.Call("Plugin.Input", opts, &resp)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		err = resp.Error
		return "", err
	}

	return resp.Value, nil
}

type UIInputInputResponse struct {
	Value string
	Error *plugin.BasicError
}

// UIInputServer is a net/rpc compatible structure for serving
// a UIInputServer. This should not be used directly.
type UIInputServer struct {
	UIInput terracina.UIInput
}

func (s *UIInputServer) Input(
	opts *terracina.InputOpts,
	reply *UIInputInputResponse) error {
	value, err := s.UIInput.Input(context.Background(), opts)
	*reply = UIInputInputResponse{
		Value: value,
		Error: plugin.NewBasicError(err),
	}

	return nil
}
