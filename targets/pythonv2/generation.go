// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"github.com/andreychh/tgen/targets"
)

// Generation is one run that writes a Python package: the specification the
// files are rendered from and the tgen that wrote them. It is the root every
// template renders against, and reaches the specification through Spec.
type Generation struct {
	spec     Specification
	snapshot targets.Snapshot
}

// NewGeneration creates a Generation rendering spec, stamped with snapshot.
func NewGeneration(spec Specification, snapshot targets.Snapshot) Generation {
	return Generation{spec: spec, snapshot: snapshot}
}

// Spec returns the specification the files are rendered from.
func (g Generation) Spec() Specification {
	return g.spec
}

// Snapshot returns the metadata of the run: when it happened and which tgen
// performed it.
func (g Generation) Snapshot() targets.Snapshot {
	return g.snapshot
}
