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

// Package gosentence contains functions for converting a slice of strings into a human readble sentence.
package gosentence

import (
	"slices"
	"strings"
)

// Transform combines each element of v inside a human readble sentence.
// Each element in v after the first one is compared against the noTransform set.
// If the element is found, it's used as is, if not, it's converted to lowercase.
func Transform(v []string, noTransform ...string) string {
	retVal := make([]string, 0, len(v))

	for idx, word := range v {
		if idx == 0 {
			retVal = append(retVal, word)
		} else if slices.Contains(noTransform, word) {
			retVal = append(retVal, word)
		} else {
			retVal = append(retVal, strings.ToLower(word))
		}
	}

	return strings.Join(retVal, " ")
}
