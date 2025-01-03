// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"context"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
)

func Test_migrate_single_to_tfc(t *testing.T) {
	t.Parallel()
	skipIfMissingEnvVar(t)
	skipWithoutRemoteTerracinaVersion(t)

	ctx := context.Background()

	cases := testCases{
		"migrate using cloud workspace name strategy": {
			operations: []operationSets{
				{
					prep: func(t *testing.T, orgName, dir string) {
						tfBlock := terracinaConfigLocalBackend()
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `Successfully configured the backend "local"!`,
						},
						{
							command:         []string{"apply", "-auto-approve"},
							postInputOutput: []string{`Apply complete!`},
						},
					},
				},
				{
					prep: func(t *testing.T, orgName, dir string) {
						wsName := "new-workspace"
						tfBlock := terracinaConfigCloudBackendName(orgName, wsName)
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `Migrating from backend "local" to HCP Terracina.`,
							userInput:         []string{"yes", "yes"},
							postInputOutput: []string{
								`Should Terracina migrate your existing state?`,
								`HCP Terracina has been successfully initialized!`},
						},
						{
							command:           []string{"workspace", "list"},
							expectedCmdOutput: `new-workspace`,
						},
					},
				},
			},
			validations: func(t *testing.T, orgName string) {
				wsList, err := tfeClient.Workspaces.List(ctx, orgName, nil)
				if err != nil {
					t.Fatal(err)
				}
				ws := wsList.Items[0]
				if ws.Name != "new-workspace" {
					t.Fatalf("Expected workspace to be `new-workspace`, but is %s", ws.Name)
				}
			},
		},
		"migrate using cloud workspace tags strategy": {
			operations: []operationSets{
				{
					prep: func(t *testing.T, orgName, dir string) {
						tfBlock := terracinaConfigLocalBackend()
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `Successfully configured the backend "local"!`,
						},
						{
							command:         []string{"apply", "-auto-approve"},
							postInputOutput: []string{`Apply complete!`},
						},
					},
				},
				{
					prep: func(t *testing.T, orgName, dir string) {
						tag := "app"
						tfBlock := terracinaConfigCloudBackendTags(orgName, tag)
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `Migrating from backend "local" to HCP Terracina.`,
							userInput:         []string{"yes", "new-workspace", "yes"},
							postInputOutput: []string{
								`Should Terracina migrate your existing state?`,
								`HCP Terracina requires all workspaces to be given an explicit name.`,
								`HCP Terracina has been successfully initialized!`},
						},
						{
							command:           []string{"workspace", "list"},
							expectedCmdOutput: `new-workspace`,
						},
					},
				},
			},
			validations: func(t *testing.T, orgName string) {
				wsList, err := tfeClient.Workspaces.List(ctx, orgName, &tfe.WorkspaceListOptions{
					Tags: "app",
				})
				if err != nil {
					t.Fatal(err)
				}
				ws := wsList.Items[0]
				if ws.Name != "new-workspace" {
					t.Fatalf("Expected workspace to be `new-workspace`, but is %s", ws.Name)
				}
			},
		},
	}

	testRunner(t, cases, 1)
}
