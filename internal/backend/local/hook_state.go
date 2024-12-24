// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package local

import (
	"log"
	"sync"
	"time"

	"github.com/hashicorp/terracina/internal/schemarepo"
	"github.com/hashicorp/terracina/internal/states"
	"github.com/hashicorp/terracina/internal/states/statemgr"
	"github.com/hashicorp/terracina/internal/terracina"
)

// StateHook is a hook that continuously updates the state by calling
// WriteState on a statemgr.Full.
type StateHook struct {
	terracina.NilHook
	sync.Mutex

	StateMgr statemgr.Writer

	// If PersistInterval is nonzero then for any new state update after
	// the duration has elapsed we'll try to persist a state snapshot
	// to the persistent backend too.
	// That's only possible if field Schemas is valid, because the
	// StateMgr.PersistState function for some backends needs schemas.
	PersistInterval time.Duration

	// Schemas are the schemas to use when persisting state due to
	// PersistInterval. This is ignored if PersistInterval is zero,
	// and PersistInterval is ignored if this is nil.
	Schemas *schemarepo.Schemas

	intermediatePersist statemgr.IntermediateStatePersistInfo
}

var _ terracina.Hook = (*StateHook)(nil)

func (h *StateHook) PostStateUpdate(new *states.State) (terracina.HookAction, error) {
	h.Lock()
	defer h.Unlock()

	h.intermediatePersist.RequestedPersistInterval = h.PersistInterval

	if h.intermediatePersist.LastPersist.IsZero() {
		// The first PostStateUpdate starts the clock for intermediate
		// calls to PersistState.
		h.intermediatePersist.LastPersist = time.Now()
	}

	if h.StateMgr != nil {
		if err := h.StateMgr.WriteState(new); err != nil {
			return terracina.HookActionHalt, err
		}
		if mgrPersist, ok := h.StateMgr.(statemgr.Persister); ok && h.PersistInterval != 0 && h.Schemas != nil {
			if h.shouldPersist() {
				err := mgrPersist.PersistState(h.Schemas)
				if err != nil {
					return terracina.HookActionHalt, err
				}
				h.intermediatePersist.LastPersist = time.Now()
			} else {
				log.Printf("[DEBUG] State storage %T declined to persist a state snapshot", h.StateMgr)
			}
		}
	}

	return terracina.HookActionContinue, nil
}

func (h *StateHook) Stopping() {
	h.Lock()
	defer h.Unlock()

	// If Terracina has been asked to stop then that might mean that a hard
	// kill signal will follow shortly in case Terracina doesn't stop
	// quickly enough, and so we'll try to persist the latest state
	// snapshot in the hope that it'll give the user less recovery work to
	// do if they _do_ subsequently hard-kill Terracina during an apply.

	if mgrPersist, ok := h.StateMgr.(statemgr.Persister); ok && h.Schemas != nil {
		// While we're in the stopping phase we'll try to persist every
		// new state update to maximize every opportunity we get to avoid
		// losing track of objects that have been created or updated.
		// Terracina Core won't start any new operations after it's been
		// stopped, so at most we should see one more PostStateUpdate
		// call per already-active request.
		h.intermediatePersist.ForcePersist = true

		if h.shouldPersist() {
			err := mgrPersist.PersistState(h.Schemas)
			if err != nil {
				// This hook can't affect Terracina Core's ongoing behavior,
				// but it's a best effort thing anyway so we'll just emit a
				// log to aid with debugging.
				log.Printf("[ERROR] Failed to persist state after interruption: %s", err)
			}
		} else {
			log.Printf("[DEBUG] State storage %T declined to persist a state snapshot", h.StateMgr)
		}
	}

}

func (h *StateHook) shouldPersist() bool {
	if m, ok := h.StateMgr.(statemgr.IntermediateStateConditionalPersister); ok {
		return m.ShouldPersistIntermediateState(&h.intermediatePersist)
	}
	return statemgr.DefaultIntermediateStatePersistRule(&h.intermediatePersist)
}
