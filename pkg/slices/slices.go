// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package slices

func NewMapped[A, B any](s []A, f func(A) B) []B {
	res := make([]B, len(s))
	for i, a := range s {
		res[i] = f(a)
	}
	return res
}
