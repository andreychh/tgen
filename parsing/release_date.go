// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"time"

	"github.com/andreychh/tgen/parsing/gq"
)

type GQReleaseDate struct {
	a gq.Selection
}

func NewGQReleaseDate(a gq.Selection) GQReleaseDate {
	return GQReleaseDate{a: a}
}

func (d GQReleaseDate) AsTime() (time.Time, error) {
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
