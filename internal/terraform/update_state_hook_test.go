// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package terracina

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terracina/internal/addrs"
	"github.com/hashicorp/terracina/internal/states"
)

func TestUpdateStateHook(t *testing.T) {
	mockHook := new(MockHook)

	state := states.NewState()
	state.SetOutputValue(
		addrs.OutputValue{Name: "foo"}.Absolute(addrs.RootModuleInstance),
		cty.StringVal("hello"), false,
	)

	ctx := new(MockEvalContext)
	ctx.HookHook = mockHook
	ctx.StateState = state.SyncWrapper()

	if err := updateStateHook(ctx); err != nil {
		t.Fatalf("err: %s", err)
	}

	if !mockHook.PostStateUpdateCalled {
		t.Fatal("should call PostStateUpdate")
	}
	if mockHook.PostStateUpdateState.RootOutputValues["foo"].Value != cty.StringVal("hello") {
		t.Fatalf("wrong state passed to hook: %s", spew.Sdump(mockHook.PostStateUpdateState))
	}
}
