// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"testing"
)

func Test_init_with_empty_tags(t *testing.T) {
	t.Parallel()
	skipIfMissingEnvVar(t)
	skipWithoutRemoteTerracinaVersion(t)

	cases := testCases{
		"terracina init with cloud block - no tagged workspaces exist yet": {
			operations: []operationSets{
				{
					prep: func(t *testing.T, orgName, dir string) {
						wsTag := "emptytag"
						tfBlock := terracinaConfigCloudBackendTags(orgName, wsTag)
						writeMainTF(t, tfBlock, dir)
					},
					commands: []tfCommand{
						{
							command:           []string{"init"},
							expectedCmdOutput: `There are no workspaces with the configured tags`,
							userInput:         []string{"emptytag-prod"},
							postInputOutput:   []string{`HCP Terracina has been successfully initialized!`},
						},
					},
				},
			},
		},
	}

	testRunner(t, cases, 1)
}
