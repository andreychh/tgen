// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"github.com/andreychh/tgen/targets"
)

// Generation is one run that writes a Go package: the specification the files
// are rendered from, the package they declare, and the tgen that wrote them. It
// is the root every template renders against, and reaches the specification
// through Spec.
type Generation struct {
	spec     Specification
	pkg      string
	snapshot targets.Snapshot
}

// NewGeneration creates a Generation rendering spec into the package named
// pkg, stamped with snapshot.
func NewGeneration(spec Specification, pkg string, snapshot targets.Snapshot) Generation {
	return Generation{spec: spec, pkg: pkg, snapshot: snapshot}
}

// Spec returns the specification the files are rendered from.
func (g Generation) Spec() Specification {
	return g.spec
}

// Package returns the name of the package the generated files declare.
func (g Generation) Package() string {
	return g.pkg
}

// Snapshot returns the metadata of the run: when it happened and which tgen
// performed it.
func (g Generation) Snapshot() targets.Snapshot {
	return g.snapshot
}
