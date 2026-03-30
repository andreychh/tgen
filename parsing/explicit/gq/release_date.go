// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"fmt"
	"time"

	"github.com/andreychh/tgen/pkg/gq"
)

type ReleaseDate struct {
	a gq.Selection
}

func NewReleaseDate(a gq.Selection) ReleaseDate {
	return ReleaseDate{a: a}
}

func (d ReleaseDate) AsTime() (time.Time, error) {
	if d.a.IsEmpty() {
		return time.Time{}, errors.New("release date not found")
	}
	val, _ := d.a.Attr("href")
	parsed, err := time.Parse("#January-2-2006", val)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing release date: %w", err)
	}
	return parsed, nil
}
