// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"fmt"
	"time"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/pkg/gq"
)

type ReleaseDate struct {
	a gq.Selection
}

func NewReleaseDate(a gq.Selection) ReleaseDate {
	return ReleaseDate{a: a}
}

// Value returns the release date parsed from the anchor href.
func (d ReleaseDate) Value() (model.ReleaseDate, error) {
	if d.a.IsEmpty() {
		return model.ReleaseDate{}, errors.New("release date not found")
	}
	val, _ := d.a.Attr("href")
	parsed, err := time.Parse("#January-2-2006", val)
	if err != nil {
		return model.ReleaseDate{}, fmt.Errorf("parsing release date: %w", err)
	}
	return model.ReleaseDate(parsed), nil
}
