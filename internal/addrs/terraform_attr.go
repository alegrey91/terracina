// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package addrs

// TerracinaAttr is the address of an attribute of the "terracina" object in
// the interpolation scope, like "terracina.workspace".
type TerracinaAttr struct {
	referenceable
	Name string
}

func (ta TerracinaAttr) String() string {
	return "terracina." + ta.Name
}

func (ta TerracinaAttr) UniqueKey() UniqueKey {
	return ta // A TerracinaAttr is its own UniqueKey
}

func (ta TerracinaAttr) uniqueKeySigil() {}
