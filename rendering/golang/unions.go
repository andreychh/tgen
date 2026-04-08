// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"iter"

	"github.com/andreychh/tgen/pkg/iters"
)

type Unions struct {
	inner iter.Seq[Field]
}

func (u Unions) All() iter.Seq[Field] {
	return iters.NewFilteredSeq(u.inner, func(f Field) bool {
		is, err := f.Type().IsUnion()
		if err != nil {
			panic(err)
		}
		return is
	})
}

func (u Unions) IsEmpty() (bool, error) {
	return iters.IsEmpty(u.All()), nil
}
