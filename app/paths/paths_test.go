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

// QA: Verify the public API of the `paths` package.
package paths_test

import (
	"testing"

	"github.com/kdeconinck/assert"
	"github.com/kdeconinck/paths"
)

// UT: Get the full name (filename and extension) of a string representing a file path.
func TestName(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for tcName, tc := range map[string]struct {
		fullPathInput string
		want          string
	}{
		"When using a Linux file path \"/home/user/development/file.dll\"": {
			fullPathInput: "/home/user/development/file.dll",
			want:          "file.dll",
		},
		"When using a Windows file path \"/home/user/development/file.dll\"": {
			fullPathInput: "C:\\user\\development\\file.dll",
			want:          "file.dll",
		},
	} {
		t.Run(tcName, func(t *testing.T) {
			tc, tcName := tc, tcName // Rebind the `tc` & `tcName` variables. Required to support parallel exceution.
			t.Parallel()             // Enable parallel execution.

			// ACT.
			got := paths.Name(tc.fullPathInput)

			// ASSERT.
			assert.Equal(t, got, tc.want, "", "\n\n"+
				"UT Name:  %s\n"+
				"Input:    %v\n"+
				"\033[32mExpected: %v\033[0m\n"+
				"\033[31mActual:   %v\033[0m\n\n", tcName, tc.fullPathInput, tc.want, got)
		})
	}
}
