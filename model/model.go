// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package model defines the foundational value types shared across all layers
// of the tgen pipeline: Name, Type, Description, Key, and their companions.
package model

import "time"

type Name string

type Reference string

type Optionality bool

type Description interface {
	Value() (string, error)
	Links() ([]string, error)
}

type Key string

type DiscriminatorValue string

type ReleaseVersion string

type ReleaseDate time.Time
