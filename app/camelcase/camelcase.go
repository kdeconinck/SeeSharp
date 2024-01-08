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

// Package camelcase contains functions for working with "CamelCase" strings.
package camelcase

import (
	"unicode"
	"unicode/utf8"
)

// Split reads v treating it as "CamelCase" and returns the different words.
// If v isn't a valid UTF-8 string, when v is an empty string or when v is a string starting with a lowercase character,
// a slice with one element (v) is returned.
func Split(v string) []string {
	if !isValid(v) {
		return []string{v}
	}

	wordIndexes := make([]int, 0, len(v)/2)

	// Read the entire string, rune by rune.
	for pos := 0; pos < len(v); pos += 1 {
		if unicode.IsNumber(rune(v[pos])) {
			endPos := getNumberEndPos(v, pos)
			wordIndexes = append(wordIndexes, endPos+1)
			pos = endPos
		} else if unicode.IsUpper(rune(v[pos])) || unicode.IsLower(rune(v[pos])) {
			endPos := getWordEndPos(v, pos)
			wordIndexes = append(wordIndexes, endPos+1)
			pos = endPos
		} else if unicode.IsSpace(rune(v[pos])) {
			endPos := getSpaceEndPos(v, pos)
			wordIndexes = append(wordIndexes, endPos+1)
			pos = endPos
		}
	}

	return extractWordsByIndexes(v, wordIndexes)
}

// Returns a slice of words extracted from v based on the provided indexes.
// The provided indexes represents the end positions of each word in the string.
// The function iterates through the end indexes, calculates the corresponding start index for each word,
// and extracts the word from the original string. The extracted words are then returned in a slice.
func extractWordsByIndexes(v string, endIndexes []int) []string {
	words := make([]string, 0, len(endIndexes))

	for i := 0; i < len(endIndexes); i++ {
		endIndex := endIndexes[i]

		if endIndex >= 0 && endIndex <= len(v) {
			var startIndex int

			if i == 0 {
				startIndex = 0
			} else {
				startIndex = endIndexes[i-1]
			}

			word := v[startIndex:endIndex]
			words = append(words, word)
		}
	}

	return words
}

// Returns true if v is non-empy and valid UTF-8, false otherwise.
func isValid(v string) bool {
	return utf8.ValidString(v) && len(v) > 0
}

// Returns the end position of the numeric sequence in the string v, starting from the given index idx.
// It iterates through the string, rune by rune, to find the end position of the numeric sequence.
// The returned value points to the last numeric character in the sequence.
func getNumberEndPos(v string, idx int) int {
	for idx = idx + 1; idx < len(v) && unicode.IsNumber(rune(v[idx])); idx++ {
		// NOTE: Intentionally left blank.
	}

	return idx - 1
}

// Returns the end position of the current word in the string v, starting from the given index idx.
// The word is defined as the longest contiguous sequence of uppercase character starting from idx if the character at
// idx + 1 is uppercase. If the character at idx + 1 is lowercase, the function returns the longest contiguous sequence
// of lowercase characters starting from idx.
func getWordEndPos(v string, idx int) int {
	if idx+1 < len(v) && unicode.IsUpper(rune(v[idx+1])) {
		for idx = idx + 1; idx < len(v) && unicode.IsUpper(rune(v[idx])); idx++ {
			// NOTE: Intentionally left blank.
		}

		if idx < len(v) {
			if unicode.IsDigit(rune(v[idx])) {
				return idx - 1
			}

			return idx - 2
		} else {
			return idx - 1
		}
	}

	for idx = idx + 1; idx < len(v) && unicode.IsLower(rune(v[idx])); idx++ {
		// NOTE: Intentionally left blank.
	}

	return idx - 1
}

// Returns the end position of the space sequence in the string v, starting from the given index idx.
// It iterates through the string, rune by rune, to find the end position of the space sequence.
// The returned value points to the last space in the sequence.
func getSpaceEndPos(v string, idx int) int {
	for idx = idx + 1; idx < len(v) && unicode.IsSpace(rune(v[idx])); idx++ {
		// NOTE: Intentionally left blank.
	}

	return idx - 1
}
