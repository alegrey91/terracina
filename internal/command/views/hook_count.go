// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package views

import (
	"sync"

	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terracina/internal/addrs"
	"github.com/hashicorp/terracina/internal/plans"
	"github.com/hashicorp/terracina/internal/terracina"
)

// countHook is a hook that counts the number of resources
// added, removed, changed during the course of an apply.
type countHook struct {
	Added    int
	Changed  int
	Removed  int
	Imported int

	ToAdd          int
	ToChange       int
	ToRemove       int
	ToRemoveAndAdd int

	pending map[string]plans.Action

	sync.Mutex
	terracina.NilHook
}

var _ terracina.Hook = (*countHook)(nil)

func (h *countHook) Reset() {
	h.Lock()
	defer h.Unlock()

	h.pending = nil
	h.Added = 0
	h.Changed = 0
	h.Removed = 0
	h.Imported = 0
}

func (h *countHook) PreApply(id terracina.HookResourceIdentity, dk addrs.DeposedKey, action plans.Action, priorState, plannedNewState cty.Value) (terracina.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	if h.pending == nil {
		h.pending = make(map[string]plans.Action)
	}

	h.pending[id.Addr.String()] = action

	return terracina.HookActionContinue, nil
}

func (h *countHook) PostApply(id terracina.HookResourceIdentity, dk addrs.DeposedKey, newState cty.Value, err error) (terracina.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	if h.pending != nil {
		pendingKey := id.Addr.String()
		if action, ok := h.pending[pendingKey]; ok {
			delete(h.pending, pendingKey)

			if err == nil {
				switch action {
				case plans.CreateThenDelete, plans.DeleteThenCreate:
					h.Added++
					h.Removed++
				case plans.Create:
					h.Added++
				case plans.Delete:
					h.Removed++
				case plans.Update:
					h.Changed++
				}
			}
		}
	}

	return terracina.HookActionContinue, nil
}

func (h *countHook) PostDiff(id terracina.HookResourceIdentity, dk addrs.DeposedKey, action plans.Action, priorState, plannedNewState cty.Value) (terracina.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	// We don't count anything for data resources
	if id.Addr.Resource.Resource.Mode == addrs.DataResourceMode {
		return terracina.HookActionContinue, nil
	}

	switch action {
	case plans.CreateThenDelete, plans.DeleteThenCreate:
		h.ToRemoveAndAdd += 1
	case plans.Create:
		h.ToAdd += 1
	case plans.Delete:
		h.ToRemove += 1
	case plans.Update:
		h.ToChange += 1
	}

	return terracina.HookActionContinue, nil
}

func (h *countHook) PostApplyImport(id terracina.HookResourceIdentity, importing plans.ImportingSrc) (terracina.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	h.Imported++
	return terracina.HookActionContinue, nil
}
