// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package terracina

import (
	"github.com/hashicorp/terracina/internal/addrs"
	"github.com/hashicorp/terracina/internal/configs"
)

// GraphNodeAttachProvider is an interface that must be implemented by nodes
// that want provider configurations attached.
type GraphNodeAttachProvider interface {
	// ProviderName with no module prefix. Example: "aws".
	ProviderAddr() addrs.AbsProviderConfig

	// Sets the configuration
	AttachProvider(*configs.Provider)
}
