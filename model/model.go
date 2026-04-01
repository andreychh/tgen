// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package model defines the foundational value interfaces shared across all
// layers of the tgen pipeline: Name, Type, Description, Key, and their
// companions. Both explicit and implicit layers implement these interfaces;
// rendering consumes them.
package model

import (
	"time"

	"github.com/andreychh/tgen/model/types"
)

type Name interface { //nolint:iface // identical signature to Reference/Key/DiscriminatorValue/ReleaseVersion but semantically distinct domain concepts
	AsString() (string, error)
}

type Reference interface { //nolint:iface // identical signature to Name/Key/DiscriminatorValue/ReleaseVersion but semantically distinct domain concepts
	AsString() (string, error)
}

type Type interface {
	AsExpression() (types.TypeExpression, error)
}

type Optionality interface {
	AsBool() (bool, error)
}

type Description interface {
	AsString() (string, error)
	Links() ([]string, error)
}

type Key interface { //nolint:iface // identical signature to Name/Reference/DiscriminatorValue/ReleaseVersion but semantically distinct domain concepts
	AsString() (string, error)
}

type DiscriminatorValue interface { //nolint:iface // identical signature to Name/Reference/Key/ReleaseVersion but semantically distinct domain concepts
	AsString() (string, error)
}

type ReleaseVersion interface { //nolint:iface // identical signature to Name/Reference/Key/DiscriminatorValue but semantically distinct domain concepts
	AsString() (string, error)
}

type ReleaseDate interface {
	AsTime() (time.Time, error)
}
