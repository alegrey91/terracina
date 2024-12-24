// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package differ

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terracina/internal/command/jsonformat/computed"
	"github.com/hashicorp/terracina/internal/command/jsonformat/computed/renderers"
	"github.com/hashicorp/terracina/internal/command/jsonformat/structured"
)

func computeAttributeDiffAsPrimitive(change structured.Change, ctype cty.Type) computed.Diff {
	return asDiff(change, renderers.Primitive(change.Before, change.After, ctype))
}
