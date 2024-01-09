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

// QA: Verify the public API of the `xunit` package.
package xunit_test

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/kdeconinck/assert"
	"github.com/kdeconinck/xunit"
)

// UT: Load an XML file containing .NET test results in xUnit's v2+ XML format.
func TestLoad(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for tcName, tc := range map[string]struct {
		xmlData string
		want    xunit.TestRun
		wantErr bool
	}{
		"When using an empty string.": {
			xmlData: "",
			wantErr: true,
		},
		"When using an invalid XML document.": {
			xmlData: "{}",
			wantErr: true,
		},
		"When using an empty XML document.": {
			xmlData: "<assemblies />",
			want: xunit.TestRun{
				Computer:     "",
				User:         "",
				StartTimeRTF: "",
				EndTimeRTF:   "",
				Timestamp:    "",
				Assemblies:   make([]xunit.Assembly, 0),
			},
		},
		"When using a simple XML document": {
			xmlData: "<assemblies computer=\"WIN11\" user=\"Kevin\" timestamp=\"07/10/2023 20:53:19\" start-rtf=\"2000-12-01\" finish-rtf=\"2001-12-01\" timestamp=\"2001-12-02\">\n" +
				"  <assembly name=\"C:\\Parent\\Sub\\App.dll\" errors=\"1\" failed=\"2\" passed=\"3\" not-run=\"4\" total=\"5\" run-date=\"07/10/2023\" run-time=\"20:53:19\" time-rtf=\"2000-12-01\">\n" +
				"  </assembly>\n" +
				"</assemblies>",
			want: xunit.TestRun{
				Computer:     "WIN11",
				User:         "Kevin",
				StartTimeRTF: "2000-12-01",
				EndTimeRTF:   "2001-12-01",
				Timestamp:    "2001-12-02",
				Assemblies: []xunit.Assembly{
					{
						Name:        "App.dll",
						ErrorCount:  1,
						PassedCount: 3,
						FailedCount: 2,
						NotRunCount: 4,
						TotalCount:  5,
						RunDate:     "07/10/2023",
						RunTime:     "20:53:19",
						TimeRTF:     "2000-12-01",
						TestGroups:  make([]*xunit.TestGroup, 0),
					},
				},
			},
		},
		"When using a complex XML document": {
			xmlData: "<assemblies computer=\"WIN11\" user=\"Kevin\" timestamp=\"07/10/2023 20:53:19\" start-rtf=\"2000-12-01\" finish-rtf=\"2001-12-01\" timestamp=\"2001-12-02\">\n" +
				"  <assembly name=\"~/parent/sub/app.dll\" errors=\"1\" failed=\"2\" passed=\"3\" not-run=\"4\" total=\"5\" run-date=\"07/10/2023\" run-time=\"20:53:19\" time-rtf=\"2000-12-01\">\n" +
				"    <collection>\n" +

				// NOTE: A test which has a display name (the name of the test doesn't start with the test's type).
				"      <test name=\"A test with a display name.\" type=\"NS\" result=\"Pass\">\n" +
				"        <traits />\n" +
				"      </test>\n" +

				// NOTE: A NON nested test without a display name (the name of the test starts with the test's type).
				"      <test name=\"NS1.Class.SubClass.TestClass.TestMethod\" type=\"NS1.Class.SubClass.TestClass.TestMethod\" result=\"Fail\">\n" +
				"        <traits />\n" +
				"      </test>\n" +

				// NOTE: A NON nested parameterized test without a display name (the name of the test starts with the test's type).
				"      <test name=\"NS1.Class.SubClass.TestClass.ParameterizedTestMethod(arg: null)\" type=\"NS1.Class.SubClass.TestClass.ParameterizedTestMethod\" result=\"Fail\">\n" +
				"        <traits />\n" +
				"      </test>\n" +

				// NOTE: A nested test (the name of the test starts with the test's type).
				"      <test name=\"NS1.Class.SubClass.TestClass+Method+Scenario+SubScenario.Result\" type=\"NS1.Class.SubClass.TestClass+Method+Scenario+SubScenario\" result=\"Pass\">\n" +
				"        <traits />\n" +
				"      </test>\n" +

				// NOTE: A nested test (it contains the `+` character), in already existing group.
				"      <test name=\"NS1.Class.SubClass.TestClass+Method+Scenario2+SubScenario.Result\" type=\"NS1.Class.SubClass.TestClass+Method+Scenario2+SubScenario\" result=\"Pass\">\n" +
				"        <traits />\n" +
				"      </test>\n" +

				// NOTE: A test which has a display name (the name of the test doesn't start with the test's type), belonging to a single trait.
				"      <test name=\"A test with a display name (with a trait).\" type=\"NS\" result=\"Pass\">\n" +
				"        <traits>\n" +
				"          <trait name=\"Category\" value=\"Unit\" />\n" +
				"        </traits>\n" +
				"      </test>\n" +

				// NOTE: A test which has a display name (the name of the test doesn't start with the test's type), belonging to multiple traits.
				"      <test name=\"A test with a display name (with multiple traits).\" type=\"NS\" result=\"Pass\">\n" +
				"        <traits>\n" +
				"          <trait name=\"Category\" value=\"Unit\" />\n" +
				"          <trait name=\"Timing\" value=\"Slow\" />\n" +
				"        </traits>\n" +
				"      </test>\n" +
				"    </collection>\n" +
				"  </assembly>\n" +
				"</assemblies>",
			want: xunit.TestRun{
				Computer:     "WIN11",
				User:         "Kevin",
				StartTimeRTF: "2000-12-01",
				EndTimeRTF:   "2001-12-01",
				Timestamp:    "2001-12-02",
				Assemblies: []xunit.Assembly{
					{
						Name:        "app.dll",
						ErrorCount:  1,
						PassedCount: 3,
						FailedCount: 2,
						NotRunCount: 4,
						TotalCount:  5,
						RunDate:     "07/10/2023",
						RunTime:     "20:53:19",
						TimeRTF:     "2000-12-01",
						TestGroups: []*xunit.TestGroup{
							{
								Name: "",
								Tests: []xunit.TestCase{
									{
										Name:   "A test with a display name.",
										Result: "Pass",
									},
									{
										Name:   "Test method",
										Result: "Fail",
									},
									{
										Name:   "Parameterized test method",
										Result: "Fail",
									},
								},
								Groups: []*xunit.TestGroup{
									{
										Name:  "Test class",
										Tests: nil,
										Groups: []*xunit.TestGroup{
											{
												Name:  "Method",
												Tests: nil,
												Groups: []*xunit.TestGroup{
													{
														Name:  "Scenario",
														Tests: nil,
														Groups: []*xunit.TestGroup{
															{
																Name: "Sub scenario",
																Tests: []xunit.TestCase{
																	{
																		Name:   "Result",
																		Result: "Pass",
																	},
																},
															},
														},
													},
													{
														Name:  "Scenario 2",
														Tests: nil,
														Groups: []*xunit.TestGroup{
															{
																Name: "Sub scenario",
																Tests: []xunit.TestCase{
																	{
																		Name:   "Result",
																		Result: "Pass",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							{
								Name: "Category - Unit",
								Tests: []xunit.TestCase{
									{
										Name:   "A test with a display name (with a trait).",
										Result: "Pass",
									},
									{
										Name:   "A test with a display name (with multiple traits).",
										Result: "Pass",
									},
								},
							},
							{
								Name: "Timing - Slow",
								Tests: []xunit.TestCase{
									{
										Name:   "A test with a display name (with multiple traits).",
										Result: "Pass",
									},
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tcName, func(t *testing.T) {
			tc, tcName := tc, tcName // Rebind the `tc` & `tcName` variables. Required to support parallel exceution.
			t.Parallel()             // Enable parallel execution.

			// HELPER FUNCTIONS.
			fmtJson := func(v xunit.TestRun) string {
				b, _ := json.MarshalIndent(v, "", "  ")
				res := strings.Replace(string(b), "\n", "\n            ", -1)

				return res
			}

			fmtXml := func(v string) string {
				v = strings.Replace(v, "\n", "\n            ", -1)

				return v
			}

			// ARRANGE.
			rdr := strings.NewReader(tc.xmlData)

			// ACT.
			got, err := xunit.Load(rdr)

			// ASSERT.
			if tc.wantErr {
				assert.NotNil(t, err, "", "\n\n"+
					"UT Name:    %s\n"+
					"XML Input:  %s\n"+
					"\033[32mExpected:   Error, NOT <nil>\033[0m\n"+
					"\033[31mActual:     Error, %v\033[0m\n\n", tcName, fmtXml(tc.xmlData), err)
			}

			if !tc.wantErr {
				assert.Equal(t, err, nil, "", "\n\n"+
					"UT Name:    %s\n"+
					"XML Input:  %s\n"+
					"\033[32mExpected:   Error, <nil>\033[0m\n"+
					"\033[31mActual:     Error, %v\033[0m\n\n", tcName, fmtXml(tc.xmlData), err)
			}

			assert.EqualFn(t, got, tc.want, func(got xunit.TestRun, want xunit.TestRun) bool {
				return reflect.DeepEqual(got, want)
			}, "", "\n\n"+
				"UT Name:    %s\n"+
				"XML Input:  %s\n"+
				"\033[32mExpected:   %s\033[0m\n"+
				"\033[31mActual:     %s\033[0m\n\n",
				tcName, fmtXml(tc.xmlData), fmtJson(tc.want), fmtJson(got))
		})
	}
}

// Benchmark: Load an XML file containing .NET test results in xUnit's v2+ XML format.
func BenchmarkLoad_MultipleAssemblies(b *testing.B) {
	xmlData := "<assemblies>\n"

	for i := 0; i < 1_000; i++ {
		xmlData += "  <assembly />\n"
	}

	xmlData += "</assemblies>"

	benchmarkLoad(xmlData, b)
}

// Benchmark: Load an XML file containing .NET test results in xUnit's v2+ XML format.
func BenchmarkLoad_MultipleTraits(b *testing.B) {
	xmlData := "<assemblies>\n" +
		"  <assembly>\n" +
		"    <collection>\n" +
		"      <test>\n" +
		"        <traits>\n"

	for i := 0; i < 1_000; i++ {
		xmlData += "		   <trait name=\"Idx\" value=\"" + strconv.Itoa(i) + "\" />\n"
	}

	xmlData += "        </traits>\n" +
		"      </test>\n" +
		"    </collection>\n" +
		"  </assembly>\n" +
		"</assemblies>"

	benchmarkLoad(xmlData, b)
}

// Benchmark: Load an XML file containing .NET test results in xUnit's v2+ XML format.
func BenchmarkLoad_MultipleTests_MultipleTraits(b *testing.B) {
	xmlData := "<assemblies>\n" +
		"  <assembly>\n" +
		"    <collection>\n"

	for tcIdx := 0; tcIdx < 100; tcIdx++ {
		xmlData += "      <test>\n" +
			"        <traits>\n"

		for i := 0; i < 10; i++ {
			xmlData += "		   <trait name=\"Idx\" value=\"" + strconv.Itoa(i) + "\" />\n"
		}

		xmlData += "        </traits>\n" +
			"      </test>\n"
	}

	xmlData += "    </collection>\n" +
		"  </assembly>\n" +
		"</assemblies>"

	benchmarkLoad(xmlData, b)
}

// Benchmark: Load an XML file containing .NET test results in xUnit's v2+ XML format.
func benchmarkLoad(xmlData string, b *testing.B) {
	for i := 0; i < b.N; i++ {
		rdr := strings.NewReader(xmlData)
		_, _ = xunit.Load(rdr)
	}
}
