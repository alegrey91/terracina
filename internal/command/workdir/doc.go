// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package workdir models the various local artifacts and state we keep inside
// a Terracina "working directory".
//
// The working directory artifacts and settings are typically initialized or
// modified by "terracina init", after which they persist for use by other
// commands in the same directory, but are not visible to commands run in
// other working directories or on other computers.
//
// Although "terracina init" is the main command which modifies a workdir,
// other commands do sometimes make more focused modifications for settings
// which can typically change multiple times during a session, such as the
// currently-selected workspace name. Any command which modifies the working
// directory settings must discard and reload any objects which derived from
// those settings, because otherwise the existing objects will often continue
// to follow the settings that were present when they were created.
package workdir
