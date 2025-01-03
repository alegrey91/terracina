// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/terracina/internal/moduleref"
)

func TestModules_noJsonFlag(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(dir, 0755)
	testCopyDir(t, testFixturePath("modules-nested-dependencies"), dir)
	ui := new(cli.MockUi)
	view, done := testView(t)
	defer testChdir(t, dir)()

	cmd := &ModulesCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{}
	code := cmd.Run(args)
	if code != 0 {
		t.Fatalf("Got a non-zero exit code: %d\n", code)
	}

	actual := done(t).All()

	for _, part := range expectedOutputHuman {
		if !strings.Contains(actual, part) {
			t.Fatalf("unexpected output: %s\n", part)
		}
	}
}

func TestModules_noJsonFlag_noModules(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(dir, 0755)
	ui := new(cli.MockUi)
	view, done := testView(t)
	defer testChdir(t, dir)()

	cmd := &ModulesCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{}
	code := cmd.Run(args)
	if code != 0 {
		t.Fatalf("Got a non-zero exit code: %d\n", code)
	}

	actual := done(t).All()

	if diff := cmp.Diff("No modules found in configuration.\n", actual); diff != "" {
		t.Fatalf("unexpected output (-want +got):\n%s", diff)
	}
}

func TestModules_fullCmd(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(dir, 0755)
	testCopyDir(t, testFixturePath("modules-nested-dependencies"), dir)

	ui := new(cli.MockUi)
	view, done := testView(t)
	defer testChdir(t, dir)()

	cmd := &ModulesCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{"-json"}
	code := cmd.Run(args)
	if code != 0 {
		t.Fatalf("Got a non-zero exit code: %d\n", code)
	}

	output := done(t).All()
	compareJSONOutput(t, output, expectedOutputJSON)
}

func TestModules_fullCmd_unreferencedEntries(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(dir, 0755)
	testCopyDir(t, testFixturePath("modules-unreferenced-entries"), dir)

	ui := new(cli.MockUi)
	view, done := testView(t)
	defer testChdir(t, dir)()

	cmd := &ModulesCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{"-json"}
	code := cmd.Run(args)
	if code != 0 {
		t.Fatalf("Got a non-zero exit code: %d\n", code)
	}
	output := done(t).All()
	compareJSONOutput(t, output, expectedOutputJSON)
}

func TestModules_uninstalledModules(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(dir, 0755)
	testCopyDir(t, testFixturePath("modules-uninstalled-entries"), dir)

	ui := new(cli.MockUi)
	view, done := testView(t)
	defer testChdir(t, dir)()

	cmd := &ModulesCommand{
		Meta: Meta{
			testingOverrides: metaOverridesForProvider(testProvider()),
			Ui:               ui,
			View:             view,
		},
	}

	args := []string{"-json"}
	code := cmd.Run(args)
	if code == 0 {
		t.Fatal("Expected a non-zero exit code\n")
	}
	output := done(t).All()
	if !strings.Contains(output, "Module not installed") {
		t.Fatalf("expected to see a `not installed` error message: %s\n", output)
	}

	if !strings.Contains(output, `Run "terracina init"`) {
		t.Fatalf("expected error message to ask user to run terracina init: %s\n", output)
	}
}

func compareJSONOutput(t *testing.T, got string, want string) {
	var expected, actual moduleref.Manifest

	if err := json.Unmarshal([]byte(got), &actual); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(want), &expected); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	sort.Slice(actual.Records, func(i, j int) bool {
		return actual.Records[i].Key < actual.Records[j].Key
	})
	sort.Slice(expected.Records, func(i, j int) bool {
		return expected.Records[i].Key < expected.Records[j].Key
	})

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("unexpected output, got: %s\n, want:%s\n", got, want)
	}
}

var expectedOutputJSON = `{"format_version":"1.0","modules":[{"key":"test","source":"./mods/test","version":""},{"key":"test2","source":"./test2","version":""},{"key":"test3","source":"./test3","version":""},{"key":"other","source":"./mods/other","version":""}]}`

var expectedOutputHuman = []string{"── \"other\"[./mods/other]", "── \"test\"[./mods/test]\n    └── \"test2\"[./test2]\n        └── \"test3\"[./test3]"}
