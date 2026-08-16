// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"github.com/andreychh/tgen/model"
)

// Stub represents a declaration the target does not spell yet: it holds the
// place the page gave it and names the kind that would fill it. Nothing
// generated from a Stub is Python — it is scaffolding, and the last kind to
// land takes it out of the target with it.
type Stub struct {
	ref  model.Reference
	kind string
}

// NewStub creates a Stub standing in for the declaration at ref, of the named
// kind.
func NewStub(ref model.Reference, kind string) Stub {
	return Stub{ref: ref, kind: kind}
}

// Ref implements [Declaration].
func (s Stub) Ref() string {
	return string(s.ref)
}

// Template implements [Declaration].
func (s Stub) Template() string {
	return "stub"
}

// Kind returns the name of the kind the declaration would be rendered as.
func (s Stub) Kind() string {
	return s.kind
}
