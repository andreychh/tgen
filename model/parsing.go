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

type Name interface {
	AsString() (string, error)
}

type Reference interface {
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

type Key interface {
	AsString() (string, error)
}

type DiscriminatorValue interface {
	AsString() (string, error)
}

type ReleaseVersion interface {
	AsString() (string, error)
}

type ReleaseDate interface {
	AsTime() (time.Time, error)
}
