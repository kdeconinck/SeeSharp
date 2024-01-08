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

// QA: Verify the public API of the `camelcase` package.
package camelcase_test

import (
	"strings"
	"testing"

	"github.com/kdeconinck/assert"
	"github.com/kdeconinck/camelcase"
)

// UT: Split a "CamelCase" string.
func TestEqual(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for tcName, tc := range map[string]struct {
		vInput       string
		noSplitInput []string
		want         []string
	}{
		"When reading \"\"": {
			vInput: "",
			want:   []string{""},
		},
		"When reading a string that's NOT valid UTF-8.": {
			vInput: "BadUTF8\xe2\xe2\xa1",
			want:   []string{"BadUTF8\xe2\xe2\xa1"},
		},
		"When reading \"lowercase\".": {
			vInput: "lowercase",
			want:   []string{"lowercase"},
		},
		"When reading \"Uppercase\"": {
			vInput: "Uppercase",
			want:   []string{"Uppercase"},
		},
		"When reading \"MultipleWords\"": {
			vInput: "MultipleWords",
			want:   []string{"Multiple", "Words"},
		},
		"When reading \"MyID\"": {
			vInput: "MyID",
			want:   []string{"My", "ID"},
		},
		"When reading \"HTML\"": {
			vInput: "HTML",
			want:   []string{"HTML"},
		},
		"When reading \"PDFLoader\"": {
			vInput: "PDFLoader",
			want:   []string{"PDF", "Loader"},
		},
		"When reading \"ASample\"": {
			vInput: "ASample",
			want:   []string{"A", "Sample"},
		},
		"When reading \"EasyXMLParser\"": {
			vInput: "EasyXMLParser",
			want:   []string{"Easy", "XML", "Parser"},
		},
		"When reading \"vimRPCPlugin\"": {
			vInput: "vimRPCPlugin",
			want:   []string{"vim", "RPC", "Plugin"},
		},
		"When reading \"GL11Version\".": {
			vInput: "GL11Version",
			want:   []string{"GL", "11", "Version"},
		},
		"When reading \"10\".": {
			vInput: "10",
			want:   []string{"10"},
		},
		"When reading \"10Validators\"": {
			vInput: "10Validators",
			want:   []string{"10", "Validators"},
		},
		"When reading \"May5\".": {
			vInput: "May5",
			want:   []string{"May", "5"},
		},
		"When reading \"BFG9000\".": {
			vInput: "BFG9000",
			want:   []string{"BFG", "9000"},
		},
		"When reading \"Html2Version\"": {
			vInput: "Html2Version",
			want:   []string{"Html", "2", "Version"},
		},
		"When reading \"5May2000\"": {
			vInput: "5May2000",
			want:   []string{"5", "May", "2000"},
		},
		"When reading \"Two  spaces\"": {
			vInput: "Two  spaces",
			want:   []string{"Two", "  ", "spaces"},
		},
	} {
		t.Run(tcName, func(t *testing.T) {
			tc, tcName := tc, tcName // Rebind the `tc` & `tcName` variables. Required to support parallel exceution.
			t.Parallel()             // Enable parallel execution.

			// ACT.
			got := camelcase.Split(tc.vInput)

			// ASSERT.
			assert.EqualS(t, got, tc.want, "", "\n\n"+
				"UT Name:  %s\n"+
				"Input:    %v\n"+
				"\033[32mExpected: %v\033[0m\n"+
				"\033[31mActual:   %v\033[0m\n\n", tcName, tc.vInput, tc.want, got)
		})
	}
}

// Benchmark: Split a "CamelCase" string.
func BenchmarkSplit(b *testing.B) {
	for bName, bench := range map[string]struct {
		vInput       string
		noSplitInput []string
	}{
		"When reading \"\"": {
			vInput: "",
		},
		"When reading a string that's NOT valid UTF-8.": {
			vInput: "BadUTF8\xe2\xe2\xa1",
		},
		"When reading \"lowercase\".": {
			vInput: "lowercase",
		},
		"When reading \"Uppercase\"": {
			vInput: "Uppercase",
		},
		"When reading \"MultipleWords\"": {
			vInput: "MultipleWords",
		},
		"When reading \"MyID\"": {
			vInput: "MyID",
		},
		"When reading \"HTML\"": {
			vInput: "HTML",
		},
		"When reading \"PDFLoader\"": {
			vInput: "PDFLoader",
		},
		"When reading \"ASample\"": {
			vInput: "ASample",
		},
		"When reading \"EasyXMLParser\"": {
			vInput: "EasyXMLParser",
		},
		"When reading \"vimRPCPlugin\"": {
			vInput: "vimRPCPlugin",
		},
		"When reading \"GL11Version\".": {
			vInput: "GL11Version",
		},
		"When reading \"10\".": {
			vInput: "10",
		},
		"When reading \"10Validators\"": {
			vInput: "10Validators",
		},
		"When reading \"May5\".": {
			vInput: "May5",
		},
		"When reading \"BFG9000\".": {
			vInput: "BFG9000",
		},
		"When reading \"Html2Version\"": {
			vInput: "Html2Version",
		},
		"When reading \"5May2000\"": {
			vInput: "5May2000",
		},
		"When reading \"Two  spaces\"": {
			vInput: "Two  spaces",
		},
	} {
		b.Run(bName, func(b *testing.B) {
			// ARRANGE.
			var s strings.Builder

			for i := 0; i < 1_000_000; i++ {
				s.WriteString(bench.vInput)
			}

			input := s.String()

			// RESET.
			b.ResetTimer()

			// EXECUTION.
			for i := 0; i < b.N; i++ {
				_ = camelcase.Split(input)
			}
		})
	}
}
