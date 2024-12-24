// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package moduledeps contains types that can be used to describe the
// providers required for all of the modules in a module tree.
//
// It does not itself contain the functionality for populating such
// data structures; that's in Terracina core, since this package intentionally
// does not depend on terracina core to avoid package dependency cycles.
package moduledeps
