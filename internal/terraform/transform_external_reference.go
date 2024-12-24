// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package terracina

import "github.com/hashicorp/terracina/internal/addrs"

// ExternalReferenceTransformer will add a GraphNodeReferencer into the graph
// that makes no changes to the graph itself but, by referencing the addresses
// within ExternalReferences, ensures that any temporary nodes that are required
// by an external caller, such as the terracina testing framework, are not
// skipped because they are not referenced from within the module.
type ExternalReferenceTransformer struct {
	ExternalReferences []*addrs.Reference
}

func (t *ExternalReferenceTransformer) Transform(g *Graph) error {
	if len(t.ExternalReferences) == 0 {
		return nil
	}

	g.Add(&nodeExternalReference{
		ExternalReferences: t.ExternalReferences,
	})
	return nil
}
