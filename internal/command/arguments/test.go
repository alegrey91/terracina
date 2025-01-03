// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package arguments

import (
	"github.com/hashicorp/terracina/internal/configs"
	"github.com/hashicorp/terracina/internal/tfdiags"
)

// Test represents the command-line arguments for the test command.
type Test struct {
	// CloudRunSource specifies the remote private module that this test run
	// should execute against in a remote HCP Terracina run.
	CloudRunSource string

	// Filter contains a list of test files to execute. If empty, all test files
	// will be executed.
	Filter []string

	// TestDirectory allows the user to override the directory that the test
	// command will use to discover test files, defaults to "tests". Regardless
	// of the value here, test files within the configuration directory will
	// always be discovered.
	TestDirectory string

	// ViewType specifies which output format to use: human or JSON.
	ViewType ViewType

	// JUnitXMLFile specifies an optional filename to write a JUnit XML test
	// result report to, in addition to the information written to the selected
	// view type.
	JUnitXMLFile string

	// You can specify common variables for all tests from the command line.
	Vars *Vars

	// Verbose tells the test command to print out the plan either in
	// human-readable format or JSON for each run step depending on the
	// ViewType.
	Verbose bool
}

func ParseTest(args []string) (*Test, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics

	test := Test{
		Vars: new(Vars),
	}

	var jsonOutput bool
	cmdFlags := extendedFlagSet("test", nil, nil, test.Vars)
	cmdFlags.Var((*FlagStringSlice)(&test.Filter), "filter", "filter")
	cmdFlags.StringVar(&test.TestDirectory, "test-directory", configs.DefaultTestDirectory, "test-directory")
	cmdFlags.BoolVar(&jsonOutput, "json", false, "json")
	cmdFlags.StringVar(&test.JUnitXMLFile, "junit-xml", "", "junit-xml")
	cmdFlags.BoolVar(&test.Verbose, "verbose", false, "verbose")

	// TODO: Finalise the name of this flag.
	cmdFlags.StringVar(&test.CloudRunSource, "cloud-run", "", "cloud-run")

	if err := cmdFlags.Parse(args); err != nil {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Failed to parse command-line flags",
			err.Error()))
	}

	switch {
	case jsonOutput:
		test.ViewType = ViewJSON
	default:
		test.ViewType = ViewHuman
	}

	return &test, diags
}
