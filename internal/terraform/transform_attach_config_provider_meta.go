// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package terracina

import (
	"github.com/hashicorp/terracina/internal/addrs"
	"github.com/hashicorp/terracina/internal/configs"
)

// GraphNodeAttachProviderMetaConfigs is an interface that must be implemented
// by nodes that want provider meta configurations attached.
type GraphNodeAttachProviderMetaConfigs interface {
	GraphNodeConfigResource

	// Sets the configuration
	AttachProviderMetaConfigs(map[addrs.Provider]*configs.ProviderMeta)
}
