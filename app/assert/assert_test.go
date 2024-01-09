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

// QA: Verify the public API of the `assert` package.
package assert_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kdeconinck/assert"
)

// Wraps the testing.TB struct and add a field for storing the failure message.
type testableT struct {
	testing.TB
	isFailed   bool
	failureMsg string
}

// Fatal flags t as failed and formats args using fmt.Sprintf and stores the result in t.
func (t *testableT) Fatalf(format string, args ...any) {
	t.isFailed = true
	t.failureMsg = fmt.Sprintf(format, args...)
}

// Failed returns true if t is marked as failed, false otherwise.
func (t *testableT) Failed() bool {
	return t.isFailed
}

// UT: Compare 2 values for equality.
func TestEqual(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for tcName, tc := range map[string]struct {
		gotInput, wantInput any
		nameInput           string
		want                string
	}{
		"When `got` and `want` are NOT equal.": {
			gotInput: false, wantInput: true,
			nameInput: "IsDigit(\"0\")",
			want:      "IsDigit(\"0\") = false, want true",
		},
		"When `got` and `want` are equal.": {
			gotInput: true, wantInput: true,
			nameInput: "IsDigit(\"0\")",
		},
		"When comparing `got` against `nil`.": {
			gotInput: errors.New("IO Error"), wantInput: nil,
			nameInput: "GetErr(nil)",
			want:      "GetErr(nil) = IO Error, want <nil>",
		},
	} {
		t.Run(tcName, func(t *testing.T) {
			tc := tc     // Rebind the `tc` variable. Required to support parallel exceution.
			t.Parallel() // Enable parallel execution.

			// ARRANGE.
			testingT := &testableT{TB: t}

			// ACT.
			assert.Equal(testingT, tc.gotInput, tc.wantInput, tc.nameInput)

			// ASSERT.
			if testingT.failureMsg != tc.want {
				t.Fatalf("Failure message = \"%s\", want \"%s\"", testingT.failureMsg, tc.want)
			}
		})
	}
}

// UT: Compare 2 values for equality (with a custom message).
func TestEqualWithCustomMessage(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	// ARRANGE.
	testingT := &testableT{TB: t}

	// ACT.
	assert.Equal(testingT, false, true, "", "UT Failed: `IsDigit(\"0\")` - got %t, want %t.", false, true)

	// ASSERT.
	if testingT.failureMsg != "UT Failed: `IsDigit(\"0\")` - got false, want true." {
		t.Fatalf("Failure message = \"%s\", want \"%s\"", testingT.failureMsg, "UT Failed: `IsDigit(\"0\")` - got false, want true.")
	}
}

// UT: Compare 2 slices for equality.
func TestEqualS(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for tcName, tc := range map[string]struct {
		gotInput, wantInput []int
		nameInput           string
		want                string
	}{
		"When `got` and `want` have a different amount of elements.": {
			gotInput: []int{1, 2, 3}, wantInput: []int{},
			nameInput: "Equal([]int{1, 2, 3}, []int{})",
			want:      "Equal([]int{1, 2, 3}, []int{}) - Unequal slice length = 3, want 0",
		},
		"When `got` and `want` are NOT equal.": {
			gotInput: []int{1, 2, 3}, wantInput: []int{3, 2, 1},
			nameInput: "Equal([]int{1, 2, 3}, []int{3, 2, 1})",
			want:      "Equal([]int{1, 2, 3}, []int{3, 2, 1}) - Idx #0 = 1, want 3",
		},
		"When `got` and `want` are equal.": {
			gotInput: []int{1, 2, 3}, wantInput: []int{1, 2, 3},
			nameInput: "Equal([]int{1, 2, 3}, []int{1, 2, 3})",
		},
	} {
		t.Run(tcName, func(t *testing.T) {
			tc := tc     // Rebind the `tc` variable. Required to support parallel exceution.
			t.Parallel() // Enable parallel execution.

			// ARRANGE.
			testingT := &testableT{TB: t}

			// ACT.
			assert.EqualS(testingT, tc.gotInput, tc.wantInput, tc.nameInput)

			// ASSERT.
			if testingT.failureMsg != tc.want {
				t.Fatalf("Failure message = \"%s\", want \"%s\"", testingT.failureMsg, tc.want)
			}
		})
	}
}

// UT: Compare 2 slices for equality (with a custom message).
func TestEqualSWithCustomMessage(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for tcName, tc := range map[string]struct {
		gotInput, wantInput []int
		msgInput            []any
		want                string
	}{
		"When `got` and `want` have a different amount of elements.": {
			gotInput: []int{1, 2, 3}, wantInput: []int{},
			msgInput: []any{"UT Failed: `Equal([]int{1, 2, 3}, []int{})` - got (length) %d, want (length) %d.", 3, 0},
			want:     "UT Failed: Equal([]int{1, 2, 3}, []int{}) - got (length) 3, want (length) 0.",
		},
		"When `got` and `want` are NOT equal.": {
			gotInput: []int{1, 2, 3}, wantInput: []int{3, 2, 1},
			msgInput: []any{"UT Failed: `Equal([]int{1, 2, 3}, []int{3, 2, 1})` - Element @ idx#%d = %d, want %d.", 0, 1, 3},
			want:     "Equal([]int{1, 2, 3}, []int{}) - Element @ idx #0 = 1, want 3.",
		},
	} {
		t.Run(tcName, func(t *testing.T) {
			tc := tc     // Rebind the `tc` variable. Required to support parallel exceution.
			t.Parallel() // Enable parallel execution.

			// ARRANGE.
			testingT := &testableT{TB: t}

			// ACT.
			assert.EqualS(testingT, tc.gotInput, tc.wantInput, "", tc.msgInput...)

			// ASSERT.
			if testingT.failureMsg != fmt.Sprintf(tc.msgInput[0].(string), tc.msgInput[1:]...) {
				t.Fatalf("Failure message = \"%s\", want \"%s\"", testingT.failureMsg, fmt.Sprintf(tc.msgInput[0].(string), tc.msgInput[1:]...))
			}
		})
	}
}
