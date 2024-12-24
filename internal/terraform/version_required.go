// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package terracina

import (
	"github.com/hashicorp/terracina/internal/tfdiags"

	"github.com/hashicorp/terracina/internal/configs"
)

// CheckCoreVersionRequirements visits each of the modules in the given
// configuration tree and verifies that any given Core version constraints
// match with the version of Terracina Core that is being used.
//
// The returned diagnostics will contain errors if any constraints do not match.
// The returned diagnostics might also return warnings, which should be
// displayed to the user.
func CheckCoreVersionRequirements(config *configs.Config) tfdiags.Diagnostics {
	if config == nil {
		return nil
	}

	var diags tfdiags.Diagnostics
	diags = diags.Append(config.CheckCoreVersionRequirements())

	return diags
}
