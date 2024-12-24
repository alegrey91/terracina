// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package clistate

import (
	"testing"

	"github.com/hashicorp/terracina/internal/command/arguments"
	"github.com/hashicorp/terracina/internal/command/views"
	"github.com/hashicorp/terracina/internal/states/statemgr"
	"github.com/hashicorp/terracina/internal/terminal"
)

func TestUnlock(t *testing.T) {
	streams, _ := terminal.StreamsForTesting(t)
	view := views.NewView(streams)

	l := NewLocker(0, views.NewStateLocker(arguments.ViewHuman, view))
	l.Lock(statemgr.NewUnlockErrorFull(nil, nil), "test-lock")

	diags := l.Unlock()
	if diags.HasErrors() {
		t.Log(diags.Err().Error())
	} else {
		t.Error("expected error")
	}
}
