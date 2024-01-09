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

// Package xunit contains functions for parsing XML files containing .NET test results in xUnit's v2+ XML format.
// More information regarding this format can be found @ https://xunit.net/docs/format-xml-v2.
package xunit

import (
	"encoding/xml"
	"io"
)

// A result is the top-level element of the document. It's the result of a `dotnet test` operation in xUnit's v2+ XML
// format.
type result struct {
	XMLName       xml.Name   `xml:"assemblies"`
	Computer      string     `xml:"computer,attr"`
	FinishRTF     string     `xml:"finish-rtf,attr"`
	ID            string     `xml:"id,attr"`
	SchemaVersion string     `xml:"schema-version,attr"`
	StartRTF      string     `xml:"start-rtf,attr"`
	Timestamp     string     `xml:"timestamp,attr"`
	User          string     `xml:"user,attr"`
	Assemblies    []assembly `xml:"assembly"`
}

// An assembly contains information about the run of a single test assembly.
// This includes environmental information.
type assembly struct {
	ConfigFile      string       `xml:"config-file,attr"`
	Environment     string       `xml:"environment,attr"`
	ErrorCount      int          `xml:"errors,attr"`
	FailedCount     int          `xml:"failed,attr"`
	FinishRTF       string       `xml:"finish-rtf,attr"`
	ID              string       `xml:"id,attr"`
	FullName        string       `xml:"name,attr"`
	NotRunCount     int          `xml:"not-run,attr"`
	PassedCount     int          `xml:"passed,attr"`
	RunDate         string       `xml:"run-date,attr"`
	RunTime         string       `xml:"run-time,attr"`
	SkippedCount    int          `xml:"skipped,attr"`
	StartRTF        string       `xml:"start-rtf,attr"`
	TargetFramework string       `xml:"target-framework,attr"`
	TestFramework   string       `xml:"test-framework,attr"`
	Time            float32      `xml:"time,attr"`
	TimeRTF         string       `xml:"time-rtf,attr"`
	Total           int          `xml:"total,attr"`
	Collections     []collection `xml:"collection"`
	ErrorSet        errorSet     `xml:"errors"`

	// Calculated fields.
	testMap map[string][]TestCase // A map that contains all the tests of the assembly, grouped by trait.
}

// A collection contains information about the run of a single test collection.
type collection struct {
	ID           string `xml:"id,attr"`
	Name         string `xml:"name,attr"`
	FailedCount  int    `xml:"failed,attr"`
	NotRunCount  int    `xml:"not-run,attr"`
	PassedCount  int    `xml:"passed,attr"`
	SkippedCount int    `xml:"skipped,attr"`
	Time         string `xml:"time,attr"`
	TimeRTF      string `xml:"time-rtf,attr"`
	TotalCount   int    `xml:"total,attr"`
	Tests        []test `xml:"test"`
}

// An errorSet contains a collection of error elements.
type errorSet struct {
	Errors []err `xml:"error"`
}

// An err contains information about an environment failure that happened outside the scope of running a single unit
// test (for example, an exception thrown while disposing of a fixture object).
type err struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

// A test contains information about the run of a single test.
type test struct {
	ID         string     `xml:"id,attr"`
	Method     string     `xml:"method,attr"`
	Name       string     `xml:"name,attr"`
	Result     string     `xml:"result,attr"`
	SourceFile string     `xml:"source-file,attr"`
	SourceLine string     `xml:"source-line,attr"`
	Time       float32    `xml:"time,attr"`
	TimeRTF    string     `xml:"time-rtf,attr"`
	Type       string     `xml:"type,attr"`
	Failure    failure    `xml:"failure"`
	Output     string     `xml:"output"`
	Reason     string     `xml:"reason"`
	TraitSet   traitSet   `xml:"traits"`
	WarningSet warningSet `xml:"warnings"`
}

// A failure contains information a test failure.
type failure struct {
	ExceptionType string `xml:"exception-type,attr"`
	Message       string `xml:"message"`
	StackTrace    string `xml:"stack-trace"`
}

// A traitSet contains a collection of trait elements.
type traitSet struct {
	Traits []trait `xml:"trait"`
}

// A trait contains a single trait name/value pair.
type trait struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// A warningSet contains a collection of warning elements.
type warningSet struct {
	Warnings []string `xml:"warning"`
}

// Returns the unmarshalled result read from the provided io.Reader.
// If an error occurs during the process, an empty result and the corresponding error are returned.
func unmarshal(rdr io.Reader) (result, error) {
	var res result

	if bytes, err := io.ReadAll(rdr); err == nil {
		if err := xml.Unmarshal(bytes, &res); err != nil {
			return result{}, err
		}
	}

	return res, nil
}
