// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package depsfile

// LockFilePath is the path, relative to a configuration's root module
// directory, where Terracina expects to find the dependency lock file for
// that configuration.
//
// This file is intended to be kept in version control, so it lives directly
// in the root module directory. The ".terracina" prefix is intended to
// suggest that it's metadata about several types of objects that ultimately
// end up in the .terracina directory after running "terracina init".
const LockFilePath = ".terracina.lock.hcl"

// DevOverrideFilePath is the path, relative to a configuration's root module
// directory, where Terracina will look to find a possible override file that
// represents a request to temporarily (within a single working directory only)
// use specific local directories in place of packages that would normally
// need to be installed from a remote location.
const DevOverrideFilePath = ".terracina/dev-overrides.hcl"
