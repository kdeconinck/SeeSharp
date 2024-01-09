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

// QA: Verify the public API of the `gosentence` package.
package gosentence_test

import (
	"testing"

	"github.com/kdeconinck/assert"
	"github.com/kdeconinck/gosentence"
)

// UT: Combine each element of a slice inside a human readable sentence.
func TestName(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for tcName, tc := range map[string]struct {
		vInput           []string
		noTransformInput []string
		want             string
	}{
		"When using an empty slice.": {
			vInput: []string{""},
			want:   "",
		},
		"When using [\"A\", \"collection\", \"of\", \"words\"].": {
			vInput: []string{"A", "collection", "of", "words"},
			want:   "A collection of words",
		},
		"When using [\"A\", \"Collection\", \"Of\", \"Words\"].": {
			vInput: []string{"A", "collection", "of", "words"},
			want:   "A collection of words",
		},
		"When using [\"An\", \"HTTP\", \"&\", \"SqlQuery\"] and .": {
			vInput:           []string{"An", "HTTP", "&", "SqlQuery"},
			noTransformInput: []string{"HTTP", "SqlQuery"},
			want:             "An HTTP & SqlQuery",
		},
	} {
		t.Run(tcName, func(t *testing.T) {
			tc, tcName := tc, tcName // Rebind the `tc` & `tcName` variables. Required to support parallel exceution.
			t.Parallel()             // Enable parallel execution.

			// ACT.
			got := gosentence.Transform(tc.vInput, tc.noTransformInput...)

			// ASSERT.
			assert.Equal(t, got, tc.want, "", "\n\n"+
				"UT Name:  %s\n"+
				"Input:    %v\n"+
				"\033[32mExpected: %v\033[0m\n"+
				"\033[31mActual:   %v\033[0m\n\n", tcName, tc.vInput, tc.want, got)
		})
	}
}
