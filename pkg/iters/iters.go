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

func NewMergedSeq[A any](seqs ...iter.Seq[A]) iter.Seq[A] {
	return func(yield func(A) bool) {
		for _, seq := range seqs {
			for a := range seq {
				if !yield(a) {
					return
				}
			}
		}
	}
}

func NewFilteredSeq[A any](seq iter.Seq[A], f func(A) bool) iter.Seq[A] {
	return func(yield func(A) bool) {
		for a := range seq {
			if !f(a) {
				continue
			}
			if !yield(a) {
				return
			}
		}
	}
}

func IsEmpty[A any](seq iter.Seq[A]) bool {
	for range seq {
		return false
	}
	return true
}
