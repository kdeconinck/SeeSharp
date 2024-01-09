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
	"io"
	"strings"

	"github.com/kdeconinck/camelcase"
	"github.com/kdeconinck/gosentence"
	"github.com/kdeconinck/maps"
	"github.com/kdeconinck/paths"
)

// TestRun contains the relevant information stored in xUnit's v2+ XML format.
type TestRun struct {
	Computer     string
	User         string
	StartTimeRTF string
	EndTimeRTF   string
	Timestamp    string
	Assemblies   []Assembly
}

// Assembly contains information about the run of a single test assembly including environmental information.
type Assembly struct {
	Name        string
	ErrorCount  int
	PassedCount int
	FailedCount int
	NotRunCount int
	TotalCount  int
	RunDate     string
	RunTime     string
	Time        float32
	TimeRTF     string
	TestGroups  []*TestGroup
}

// TestGroup is a group of tests.
type TestGroup struct {
	Name   string
	Tests  []TestCase
	Groups []*TestGroup
}

// TestCase contains information about a single test.
type TestCase struct {
	Name   string
	Result string
	Time   float32

	// Internal fields.
	groups []string
}

// Load returns a TestRun constructed from the data in rdr.
// It reads and unmarshals the data in rdr and then converts it into structs that are optimized for further processing.
// If an error occurs during the process, an empty TestRun and the corresponding error are returned.
func Load(rdr io.Reader) (TestRun, error) {
	result, err := unmarshal(rdr)

	if err != nil {
		return TestRun{}, err
	}

	return readResult(result), nil
}

// Returns a TestRun that's constructed from r which represents the root of an .NET test results in xUnit's v2+ XML
// format.
func readResult(r result) TestRun {
	testRun := TestRun{
		Computer:     r.Computer,
		User:         r.User,
		StartTimeRTF: r.StartRTF,
		EndTimeRTF:   r.FinishRTF,
		Timestamp:    r.Timestamp,
		Assemblies:   make([]Assembly, 0, len(r.Assemblies)),
	}

	for _, assembly := range r.Assemblies {
		testRun.Assemblies = append(testRun.Assemblies, Assembly{
			Name:        paths.Name(assembly.FullName),
			ErrorCount:  assembly.ErrorCount,
			PassedCount: assembly.PassedCount,
			FailedCount: assembly.FailedCount,
			NotRunCount: assembly.NotRunCount,
			TotalCount:  assembly.Total,
			RunDate:     assembly.RunDate,
			RunTime:     assembly.RunTime,
			TimeRTF:     assembly.TimeRTF,
			Time:        assembly.Time,
			TestGroups:  assembly.groupTests(),
		})
	}

	return testRun
}

// Returns an hierarchical representation of all the tests in the assembly.
// If the assembly has no tests, an empty slice is returned.
func (assembly *assembly) groupTests() []*TestGroup {
	if !assembly.hasTests() {
		return make([]*TestGroup, 0)
	}

	uniqueTraits := assembly.uniqueTraits()
	resultSet := make([]*TestGroup, 0, len(uniqueTraits))

	for idx, trait := range uniqueTraits {
		cGroup := &TestGroup{Name: trait, Tests: make([]TestCase, 0, len(assembly.testMap[trait]))}
		resultSet = append(resultSet, cGroup)

		for _, tc := range assembly.testMap[trait] {
			if len(tc.groups) == 0 {
				cGroup.Tests = append(cGroup.Tests, TestCase{Name: tc.Name, Result: tc.Result, Time: tc.Time})
			} else {
				for idx, nn := range tc.groups {
					var sGroup *TestGroup

					for _, group := range cGroup.Groups {
						if group.Name == tc.groups[idx] {
							sGroup = group

							break
						}
					}

					if sGroup == nil {
						sGroup = &TestGroup{Name: nn}
						cGroup.Groups = append(cGroup.Groups, sGroup)
					}

					if idx == len(tc.groups)-1 {
						sGroup.Tests = append(sGroup.Tests, TestCase{Name: tc.Name, Result: tc.Result, Time: tc.Time})
					}

					cGroup = sGroup
				}

				cGroup = resultSet[idx]
			}
		}
	}

	return resultSet
}

// Returns true if the assembly has tests, false otherwise.
func (assembly *assembly) hasTests() bool {
	for _, collection := range assembly.Collections {
		if len(collection.Tests) > 0 {
			return true
		}
	}

	return false
}

// Returns all all the unique trait(s), sorted by their name.
// As a result of the assembly's testMap contains the test cases per trait that are found in this assembly.
func (assembly *assembly) uniqueTraits() []string {
	assembly.testMap = make(map[string][]TestCase)

	for _, collection := range assembly.Collections {
		for _, t := range collection.Tests {
			tCase := TestCase{Name: t.friendlyName(), groups: t.groups(), Result: t.Result, Time: t.Time}

			if len(t.TraitSet.Traits) == 0 {
				assembly.testMap[""] = append(assembly.testMap[""], tCase)
			}

			for _, tTrait := range t.TraitSet.Traits {
				traitName := tTrait.friendlyName()

				assembly.testMap[traitName] = append(assembly.testMap[traitName], tCase)
			}
		}
	}

	return maps.SortedKeys(assembly.testMap)
}

// Returns a slice of groups that t belongs to.
func (t *test) groups() []string {
	if !t.isNested() {
		return make([]string, 0)
	}

	groupName := strings.Split(t.Name, "+")

	groupNameParts := make([]string, 0, len(groupName[1:len(groupName)-1]))
	groupNameParts = append(groupNameParts, groupName[0][strings.LastIndex(groupName[0], ".")+1:])
	groupNameParts = append(groupNameParts, groupName[1:len(groupName)-1]...)
	groupNameParts = append(groupNameParts, strings.Split(groupName[len(groupName)-1], ".")[0])

	groupName = make([]string, 0, len(groupNameParts))

	for _, p := range groupNameParts {
		ccSplit := camelcase.Split(p)

		groupName = append(groupName, gosentence.Transform(ccSplit))
	}

	return groupName
}

// Returns true if t is nested, false otherwise.
// When t has any `+` character in its name and when it's NOT a display name, the test is considered nested.
func (t *test) isNested() bool {
	return !t.hasDisplayName() && strings.Contains(t.Name, "+")
}

// Returns the friendly name of the trait t.
func (t *trait) friendlyName() string {
	var b strings.Builder

	// Increase the size of the builder so that it has sufficient capacity to write the data to it.
	b.Grow(len(t.Name) + len(" -  ") + len(t.Value))

	b.WriteString(t.Name)
	b.WriteString(" - ")
	b.WriteString(t.Value)

	return b.String()
}

// Returns the friendly name of the test t.
// If t has a display name, the name is returned as is, if not, the name is split based on the `.` character.
// This gives us a slices of strings where each part contains a valid C# identifier. The last part would be the name of
// the function. We feed this name to the "CamelCase" package to turn it into a readable sentence.
func (t *test) friendlyName() string {
	if t.hasDisplayName() {
		return t.Name
	}

	fnName := t.Name[strings.LastIndex(t.Name, ".")+1:]

	if strings.ContainsRune(fnName, '(') {
		fnName = fnName[:strings.LastIndex(fnName, "(")]
	}

	fnNameWords := camelcase.Split(fnName)

	return gosentence.Transform(fnNameWords)
}

// Returns true if t has a display name, false otherwise.
// When t has any space in its name, it's considered to have display name.
// This ie because by design, C# doesn't allow to have spaces in any identifier and the default name of a test is the
// concatenation (with a `.`) of all identifiers (namespace, class, subclass(es) and methods).
func (t *test) hasDisplayName() bool {
	return !strings.Contains(t.Name, t.Type)
}
