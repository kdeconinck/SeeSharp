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

// QA: Verify the public API of the `maps` package.
package maps_test

import (
	"testing"

	"github.com/kdeconinck/assert"
	"github.com/kdeconinck/maps"
)

// UT: Get the keys of map, sorted alphabetically.
func TestSortedKeys(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for tcName, tc := range map[string]struct {
		mInput map[int]bool
		want   []int
	}{
		"When using a map with keys \"[5, 2, 8, 1, 7, 6, 3, 4]\"": {
			mInput: map[int]bool{5: true, 2: true, 8: false, 1: false, 7: true, 6: false, 3: true, 4: false},
			want:   []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
	} {
		t.Run(tcName, func(t *testing.T) {
			tc, tcName := tc, tcName // Rebind the `tc` & `tcName` variables. Required to support parallel exceution.
			t.Parallel()             // Enable parallel execution.

			// ACT.
			got := maps.SortedKeys(tc.mInput)

			// ASSERT.
			assert.EqualS(t, got, tc.want, "", "\n\n"+
				"UT Name:  %s\n"+
				"Input:    %v\n"+
				"\033[32mExpected: %v\033[0m\n"+
				"\033[31mActual:   %v\033[0m\n\n", tcName, tc.mInput, tc.want, got)
		})
	}
}

// Benchmark: Get the keys of map, sorted alphabetically.
func BenchmarkSortedKeys(b *testing.B) {
	// ARRANGE.
	input := make(map[int]bool, 1_000_000)

	for idx := 0; idx < 1_000_000; idx++ {
		input[1_000_000-idx] = true
	}

	// RESET.
	b.ResetTimer()

	// EXECUTION.
	for i := 0; i < b.N; i++ {
		_ = maps.SortedKeys(input)
	}
}
