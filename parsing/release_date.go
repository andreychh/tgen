// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"fmt"
	"time"

	"github.com/andreychh/tgen/parsing/gq"
)

type ReleaseDate struct {
	selection gq.Selection
}

func NewReleaseDate(a gq.Selection) ReleaseDate {
	return ReleaseDate{selection: a}
}

func (d ReleaseDate) Value() (time.Time, error) {
	val, _ := d.selection.Attr("href")
	parsed, err := time.Parse("#January-2-2006", val)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing release date: %w", err)
	}
	return parsed, nil
}
