// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package iters

import "iter"

func NewMappedSeq[A, B any](seq iter.Seq[A], f func(A) B) iter.Seq[B] {
	return func(yield func(B) bool) {
		for a := range seq {
			if !yield(f(a)) {
				return
			}
		}
	}
}
