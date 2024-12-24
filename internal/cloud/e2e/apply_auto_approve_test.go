// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"context"
	"testing"

	tfe "github.com/hashicorp/go-tfe"
	tfversion "github.com/hashicorp/terracina/version"
)

func Test_terracina_apply_autoApprove(t *testing.T) {
	t.Parallel()
	skipIfMissingEnvVar(t)
	skipWithoutRemoteTerracinaVersion(t)

	ctx := context.Background()

	cases := testCases{
		"workspace manual apply, terracina apply without auto-approve, expect prompt": {
			operations: []operationSets{
				{
					prep: func(t *testing.T, orgName, dir string) {
						wsName := "app"
						_ = createWorkspace(t, orgName, tfe.WorkspaceCreateOptions{
							Name:             tfe.String(wsName),
							TerracinaVersion: tfe.String(tfversion.String()),
							AutoApply:        tfe.Bool(false),
						})
						tfBlock := terracinaConfigCloudBackendName(orgName, wsName)
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `HCP Terracina has been successfully initialized!`,
						},
						{
							command:           []string{"apply"},
							expectedCmdOutput: `Do you want to perform these actions in workspace "app"?`,
							userInput:         []string{"yes"},
							postInputOutput:   []string{`Apply complete!`},
						},
					},
				},
			},
			validations: func(t *testing.T, orgName string) {
				workspace, err := tfeClient.Workspaces.ReadWithOptions(ctx, orgName, "app", &tfe.WorkspaceReadOptions{Include: []tfe.WSIncludeOpt{tfe.WSCurrentRun}})
				if err != nil {
					t.Fatal(err)
				}
				if workspace.CurrentRun == nil {
					t.Fatal("Expected workspace to have run, but got nil")
				}
				if workspace.CurrentRun.Status != tfe.RunApplied {
					t.Fatalf("Expected run status to be `applied`, but is %s", workspace.CurrentRun.Status)
				}
			},
		},
		"workspace auto apply, terracina apply without auto-approve, expect prompt": {
			operations: []operationSets{
				{
					prep: func(t *testing.T, orgName, dir string) {
						wsName := "app"
						_ = createWorkspace(t, orgName, tfe.WorkspaceCreateOptions{
							Name:             tfe.String(wsName),
							TerracinaVersion: tfe.String(tfversion.String()),
							AutoApply:        tfe.Bool(true),
						})
						tfBlock := terracinaConfigCloudBackendName(orgName, wsName)
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `HCP Terracina has been successfully initialized!`,
						},
						{
							command:           []string{"apply"},
							expectedCmdOutput: `Do you want to perform these actions in workspace "app"?`,
							userInput:         []string{"yes"},
							postInputOutput:   []string{`Apply complete!`},
						},
					},
				},
			},
			validations: func(t *testing.T, orgName string) {
				workspace, err := tfeClient.Workspaces.ReadWithOptions(ctx, orgName, "app", &tfe.WorkspaceReadOptions{Include: []tfe.WSIncludeOpt{tfe.WSCurrentRun}})
				if err != nil {
					t.Fatal(err)
				}
				if workspace.CurrentRun == nil {
					t.Fatal("Expected workspace to have run, but got nil")
				}
				if workspace.CurrentRun.Status != tfe.RunApplied {
					t.Fatalf("Expected run status to be `applied`, but is %s", workspace.CurrentRun.Status)
				}
			},
		},
		"workspace manual apply, terracina apply with auto-approve, no prompt": {
			operations: []operationSets{
				{
					prep: func(t *testing.T, orgName, dir string) {
						wsName := "app"
						_ = createWorkspace(t, orgName, tfe.WorkspaceCreateOptions{
							Name:             tfe.String(wsName),
							TerracinaVersion: tfe.String(tfversion.String()),
							AutoApply:        tfe.Bool(false),
						})
						tfBlock := terracinaConfigCloudBackendName(orgName, wsName)
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `HCP Terracina has been successfully initialized!`,
						},
						{
							command:           []string{"apply", "-auto-approve"},
							expectedCmdOutput: `Apply complete!`,
						},
					},
				},
			},
			validations: func(t *testing.T, orgName string) {
				workspace, err := tfeClient.Workspaces.ReadWithOptions(ctx, orgName, "app", &tfe.WorkspaceReadOptions{Include: []tfe.WSIncludeOpt{tfe.WSCurrentRun}})
				if err != nil {
					t.Fatal(err)
				}
				if workspace.CurrentRun == nil {
					t.Fatal("Expected workspace to have run, but got nil")
				}
				if workspace.CurrentRun.Status != tfe.RunApplied {
					t.Fatalf("Expected run status to be `applied`, but is %s", workspace.CurrentRun.Status)
				}
			},
		},
		"workspace auto apply, terracina apply with auto-approve, no prompt": {
			operations: []operationSets{
				{
					prep: func(t *testing.T, orgName, dir string) {
						wsName := "app"
						_ = createWorkspace(t, orgName, tfe.WorkspaceCreateOptions{
							Name:             tfe.String(wsName),
							TerracinaVersion: tfe.String(tfversion.String()),
							AutoApply:        tfe.Bool(true),
						})
						tfBlock := terracinaConfigCloudBackendName(orgName, wsName)
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `HCP Terracina has been successfully initialized!`,
						},
						{
							command:           []string{"apply", "-auto-approve"},
							expectedCmdOutput: `Apply complete!`,
						},
					},
				},
			},
			validations: func(t *testing.T, orgName string) {
				workspace, err := tfeClient.Workspaces.ReadWithOptions(ctx, orgName, "app", &tfe.WorkspaceReadOptions{Include: []tfe.WSIncludeOpt{tfe.WSCurrentRun}})
				if err != nil {
					t.Fatal(err)
				}
				if workspace.CurrentRun == nil {
					t.Fatal("Expected workspace to have run, but got nil")
				}
				if workspace.CurrentRun.Status != tfe.RunApplied {
					t.Fatalf("Expected run status to be `applied`, but is %s", workspace.CurrentRun.Status)
				}
			},
		},
	}

	testRunner(t, cases, 1)
}
