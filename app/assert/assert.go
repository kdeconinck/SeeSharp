// =====================================================================================================================
// == LICENSE:       Copyright (c) 2024 Kevin De Coninck
// ==
// ==                Permission is hereby granted, free of charge, to any person
// ==                obtaining a copy of this software and associated documentation
// ==                files (the "Software"), to deal in the Software without
// ==                restriction, including without limitation the rights to use,
// ==                copy, modify, merge, publish, distribute, sublicense, and/or sell
// ==                copies of the Software, and to permit persons to whom the
// ==                Software is furnished to do so, subject to the following
// ==                conditions:
// ==
// ==                The above copyright notice and this permission notice shall be
// ==                included in all copies or substantial portions of the Software.
// ==
// ==                THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// ==                EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// ==                OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// ==                NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// ==                HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// ==                WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// ==                FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// ==                OTHER DEALINGS IN THE SOFTWARE.
// =====================================================================================================================

// Package assert contains functions for performing assertions in Go's standard testing framework.
package assert

import (
	"fmt"
	"testing"
)

// NotNil compares got against nil.
// If they are equal, t is marked as failed, and it's execution is terminated.
func NotNil(tb testing.TB, got any, name string, msg ...any) {
	tb.Helper()

	if got == nil {
		failT(tb, got, "NOT <nil>", name, "%s = %v, want %s", msg...)
	}
}

// Equal compares got against want for equality.
// If they are not equal, tb is marked as failed, and it's execution is terminated.
func Equal[V comparable](tb testing.TB, got, want V, name string, msg ...any) {
	tb.Helper()

	if got != want {
		failT(tb, got, want, name, "%s = %v, want %v", msg...)
	}
}

// EqualS compares got against want for equality.
// If they are not equal, tb is marked as failed, and it's execution is terminated.
func EqualS[S ~[]E, E comparable](tb testing.TB, got, want S, name string, msg ...any) {
	tb.Helper()

	if len(got) != len(want) {
		failT(tb, len(got), len(want), name, "%s - Unequal slice length = %v, want %v", msg...)
	}

	if tb.Failed() {
		return
	}

	for idx, el := range got {
		if tb.Failed() {
			break
		}

		if el != want[idx] {
			failT(tb, el, want[idx], name, fmt.Sprintf("%%s - Idx #%d = %%v, want %%v", idx), msg...)
		}
	}
}

// Marks t as failed and terminates its execution.
func failT[V any](tb testing.TB, got, want V, name, msgTemplate string, msg ...any) {
	tb.Helper()

	if name != "" {
		tb.Fatalf(msgTemplate, name, got, want)
	} else {
		tb.Fatalf(msg[0].(string), msg[1:]...)
	}
}
