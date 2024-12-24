// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package views

import (
	"bufio"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terracina/internal/addrs"
	"github.com/hashicorp/terracina/internal/command/format"
	"github.com/hashicorp/terracina/internal/command/views/json"
	"github.com/hashicorp/terracina/internal/plans"
	"github.com/hashicorp/terracina/internal/terracina"
)

func newJSONHook(view *JSONView) *jsonHook {
	return &jsonHook{
		view:             view,
		resourceProgress: make(map[string]resourceProgress),
		timeNow:          time.Now,
		timeAfter:        time.After,
		periodicUiTimer:  defaultPeriodicUiTimer,
	}
}

type jsonHook struct {
	terracina.NilHook

	view *JSONView

	// Concurrent map of resource addresses to allow tracking
	// progress, and post-action messages to share data about the resource
	resourceProgress   map[string]resourceProgress
	resourceProgressMu sync.Mutex

	// Mockable functions for testing the progress timer goroutine
	timeNow   func() time.Time
	timeAfter func(time.Duration) <-chan time.Time

	periodicUiTimer time.Duration
}

var _ terracina.Hook = (*jsonHook)(nil)

type resourceProgress struct {
	addr   addrs.AbsResourceInstance
	action plans.Action
	start  time.Time

	// done is used for post-action to stop the progress goroutine
	done chan struct{}

	// heartbeatDone is used to allow tests to safely wait for the progress
	// goroutine to finish
	heartbeatDone chan struct{}
}

func (h *jsonHook) PreApply(id terracina.HookResourceIdentity, dk addrs.DeposedKey, action plans.Action, priorState, plannedNewState cty.Value) (terracina.HookAction, error) {
	if action != plans.NoOp {
		idKey, idValue := format.ObjectValueIDOrName(priorState)
		h.view.Hook(json.NewApplyStart(id.Addr, action, idKey, idValue))
	}

	progress := resourceProgress{
		addr:          id.Addr,
		action:        action,
		start:         h.timeNow().Round(time.Second),
		done:          make(chan struct{}),
		heartbeatDone: make(chan struct{}),
	}
	h.resourceProgressMu.Lock()
	h.resourceProgress[id.Addr.String()] = progress
	h.resourceProgressMu.Unlock()

	if action != plans.NoOp {
		go h.applyingHeartbeat(progress)
	}
	return terracina.HookActionContinue, nil
}

func (h *jsonHook) applyingHeartbeat(progress resourceProgress) {
	defer close(progress.heartbeatDone)
	for {
		select {
		case <-progress.done:
			return
		case <-h.timeAfter(h.periodicUiTimer):
		}

		elapsed := h.timeNow().Round(time.Second).Sub(progress.start)
		h.view.Hook(json.NewApplyProgress(progress.addr, progress.action, elapsed))
	}
}

func (h *jsonHook) PostApply(id terracina.HookResourceIdentity, dk addrs.DeposedKey, newState cty.Value, err error) (terracina.HookAction, error) {
	key := id.Addr.String()
	h.resourceProgressMu.Lock()
	progress := h.resourceProgress[key]
	if progress.done != nil {
		close(progress.done)
	}
	delete(h.resourceProgress, key)
	h.resourceProgressMu.Unlock()

	if progress.action == plans.NoOp {
		return terracina.HookActionContinue, nil
	}

	elapsed := h.timeNow().Round(time.Second).Sub(progress.start)

	if err != nil {
		// Errors are collected and displayed post-apply, so no need to
		// re-render them here. Instead just signal that this resource failed
		// to apply.
		h.view.Hook(json.NewApplyErrored(id.Addr, progress.action, elapsed))
	} else {
		idKey, idValue := format.ObjectValueID(newState)
		h.view.Hook(json.NewApplyComplete(id.Addr, progress.action, idKey, idValue, elapsed))
	}
	return terracina.HookActionContinue, nil
}

func (h *jsonHook) PreProvisionInstanceStep(id terracina.HookResourceIdentity, typeName string) (terracina.HookAction, error) {
	h.view.Hook(json.NewProvisionStart(id.Addr, typeName))
	return terracina.HookActionContinue, nil
}

func (h *jsonHook) PostProvisionInstanceStep(id terracina.HookResourceIdentity, typeName string, err error) (terracina.HookAction, error) {
	if err != nil {
		// Errors are collected and displayed post-apply, so no need to
		// re-render them here. Instead just signal that this provisioner step
		// failed.
		h.view.Hook(json.NewProvisionErrored(id.Addr, typeName))
	} else {
		h.view.Hook(json.NewProvisionComplete(id.Addr, typeName))
	}
	return terracina.HookActionContinue, nil
}

func (h *jsonHook) ProvisionOutput(id terracina.HookResourceIdentity, typeName string, msg string) {
	s := bufio.NewScanner(strings.NewReader(msg))
	s.Split(scanLines)
	for s.Scan() {
		line := strings.TrimRightFunc(s.Text(), unicode.IsSpace)
		if line != "" {
			h.view.Hook(json.NewProvisionProgress(id.Addr, typeName, line))
		}
	}
}

func (h *jsonHook) PreRefresh(id terracina.HookResourceIdentity, dk addrs.DeposedKey, priorState cty.Value) (terracina.HookAction, error) {
	idKey, idValue := format.ObjectValueID(priorState)
	h.view.Hook(json.NewRefreshStart(id.Addr, idKey, idValue))
	return terracina.HookActionContinue, nil
}

func (h *jsonHook) PostRefresh(id terracina.HookResourceIdentity, dk addrs.DeposedKey, priorState cty.Value, newState cty.Value) (terracina.HookAction, error) {
	idKey, idValue := format.ObjectValueID(newState)
	h.view.Hook(json.NewRefreshComplete(id.Addr, idKey, idValue))
	return terracina.HookActionContinue, nil
}

func (h *jsonHook) PreEphemeralOp(id terracina.HookResourceIdentity, action plans.Action) (terracina.HookAction, error) {
	// this uses the same plans.Read action as a data source to indicate that
	// the ephemeral resource can't be processed until apply, so there is no
	// progress hook
	if action == plans.Read {
		return terracina.HookActionContinue, nil
	}

	h.view.Hook(json.NewEphemeralOpStart(id.Addr, action))
	progress := resourceProgress{
		addr:          id.Addr,
		action:        action,
		start:         h.timeNow().Round(time.Second),
		done:          make(chan struct{}),
		heartbeatDone: make(chan struct{}),
	}
	h.resourceProgressMu.Lock()
	h.resourceProgress[id.Addr.String()] = progress
	h.resourceProgressMu.Unlock()

	go h.ephemeralOpHeartbeat(progress)

	return terracina.HookActionContinue, nil
}

func (h *jsonHook) ephemeralOpHeartbeat(progress resourceProgress) {
	defer close(progress.heartbeatDone)
	for {
		select {
		case <-progress.done:
			return
		case <-h.timeAfter(h.periodicUiTimer):
		}

		elapsed := h.timeNow().Round(time.Second).Sub(progress.start)
		h.view.Hook(json.NewEphemeralOpProgress(progress.addr, progress.action, elapsed))
	}
}

func (h *jsonHook) PostEphemeralOp(id terracina.HookResourceIdentity, action plans.Action, opErr error) (terracina.HookAction, error) {
	key := id.Addr.String()
	h.resourceProgressMu.Lock()
	progress := h.resourceProgress[key]
	if progress.done != nil {
		close(progress.done)
	}
	delete(h.resourceProgress, key)
	h.resourceProgressMu.Unlock()

	if progress.action == plans.NoOp {
		return terracina.HookActionContinue, nil
	}

	elapsed := h.timeNow().Round(time.Second).Sub(progress.start)

	if opErr != nil {
		// Errors are collected and displayed post-operation, so no need to
		// re-render them here. Instead just signal that this operation failed.
		h.view.Hook(json.NewEphemeralOpErrored(id.Addr, progress.action, elapsed))
	} else {
		h.view.Hook(json.NewEphemeralOpComplete(id.Addr, progress.action, elapsed))
	}

	return terracina.HookActionContinue, nil
}
