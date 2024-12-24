// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package terracina

import "github.com/hashicorp/terracina/internal/dag"

// GraphDot returns the dot formatting of a visual representation of
// the given Terracina graph.
func GraphDot(g *Graph, opts *dag.DotOpts) (string, error) {
	return string(g.Dot(opts)), nil
}
