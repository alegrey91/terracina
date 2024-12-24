// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package terracina

import (
	"log"

	"github.com/hashicorp/terracina/internal/addrs"
	"github.com/hashicorp/terracina/internal/configs"
	"github.com/hashicorp/terracina/internal/states"
)

// OrphanOutputTransformer finds the outputs that aren't present
// in the given config that are in the state and adds them to the graph
// for deletion.
type OrphanOutputTransformer struct {
	Config   *configs.Config // Root of config tree
	State    *states.State   // State is the root state
	Planning bool
}

func (t *OrphanOutputTransformer) Transform(g *Graph) error {
	if t.State == nil {
		log.Printf("[DEBUG] No state, no orphan outputs")
		return nil
	}

	cfgs := t.Config.Module.Outputs
	for name := range t.State.RootOutputValues {
		if _, exists := cfgs[name]; exists {
			continue
		}
		g.Add(&NodeDestroyableOutput{
			Addr:     addrs.OutputValue{Name: name}.Absolute(addrs.RootModuleInstance),
			Planning: t.Planning,
		})
	}

	return nil
}
