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

// Package main implements SeeSharp, a CLI application for inspecting .NET / C# projects.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kdeconinck/xunit"
)

// The main entry point for this function.
func main() {
	// Prints the ASCII header.
	fmt.Println(" _____           _____ _")
	fmt.Println("/  ___|         /  ___| |")
	fmt.Println("\\ `--.  ___  ___\\ `--.| |__   __ _ _ __ _ __")
	fmt.Println(" `--. \\/ _ \\/ _ \\`--. \\ '_ \\ / _` | '__| '_ \\")
	fmt.Println("/\\__/ /  __/  __/\\__/ / | | | (_| | |  | |_) |")
	fmt.Println("\\____/ \\___|\\___\\____/|_| |_|\\__,_|_|  | .__/")
	fmt.Println("                                       | |")
	fmt.Println("  Version: 1.0.0                       |_|")
	fmt.Println("  Author: Kevin De Coninck")
	fmt.Println("")

	logFile := "data.xml"

	rdr, err := os.Open(logFile)

	// TODO: Handle the error gracefully.
	if err != nil {
		fmt.Println(err)
	}

	// Close the XML file as the file has been read.
	defer rdr.Close()

	// TODO: Handle the error gracefully.
	if err != nil {
		fmt.Println(err)
	}

	resultSet, err := xunit.Load(rdr)

	if err != nil {
		fmt.Println("Error encountered ...")
	}

	fmt.Printf("Input source:         %s\r\n", logFile)
	fmt.Printf("Amount of assemblies: %v\r\n", len(resultSet.Assemblies))

	if resultSet.Computer != "" {
		fmt.Printf("Computer:             %s\r\n", resultSet.Computer)
	}

	if resultSet.User != "" {
		fmt.Printf("User:                 %s\r\n", resultSet.User)
	}

	if resultSet.StartTimeRTF != "" {
		fmt.Printf("Start time:           %s\r\n", resultSet.StartTimeRTF)
	}

	if resultSet.EndTimeRTF != "" {
		fmt.Printf("End time:             %s\r\n", resultSet.EndTimeRTF)
	} else if resultSet.Timestamp != "" {
		fmt.Printf("End time:             %s\r\n", resultSet.Timestamp)
	}

	// Loop over the assemblies and print the results accordingly.
	for _, assembly := range resultSet.Assemblies {
		fmt.Println("")
		fmt.Printf("  Assembly:         %s\r\n", assembly.Name)

		if assembly.FailedCount != 0 {
			fmt.Printf("  Status:           \033[1;31m⛌ Failed (%v of %v failed).\033[0m\r\n", assembly.FailedCount, assembly.TotalCount)
		} else {
			fmt.Printf("  Status:           \033[1;32m✓ Passed (%v of %v passed).\033[0m \r\n", assembly.PassedCount, assembly.TotalCount)
		}

		fmt.Printf("  Date / time:      %s %s\r\n", assembly.RunDate, assembly.RunTime)
		fmt.Printf("  Total time:       %v seconds.\r\n", assembly.Time)

		// Print information about the assembly.
		fmt.Println("")
		fmt.Printf("    # tests:        %v\r\n", assembly.TotalCount)
		fmt.Printf("    # Passed tests: %v\r\n", assembly.PassedCount)
		fmt.Printf("    # Failed tests: %v\r\n", assembly.FailedCount)
		fmt.Printf("    # Errors:       %v\r\n", assembly.ErrorCount)
		fmt.Println("")

		// Loop the groups of the assembly.
		for _, key := range assembly.TestGroups {
			printGroup(key, 0)
		}
	}

	fmt.Println("")
}

// Prints the group with the given indentLevel to stdOut.
func printGroup(group *xunit.TestGroup, indentLevel int) {
	if group.Name != "" {
		fmt.Printf("%sGroup: %s\r\n", strings.Repeat(" ", indentLevel+1), group.Name)
	}

	for _, test := range group.Tests {
		if test.Time <= 0.05 {
			fmt.Printf("%s 🚀 %s %s (%v seconds)\r\n", strings.Repeat(" ", indentLevel), "", test.Name, test.Time)
		} else if test.Time <= 0.1 {
			fmt.Printf("%s 🕐 %s %s (%v seconds)\r\n", strings.Repeat(" ", indentLevel), "", test.Name, test.Time)
		} else {
			fmt.Printf("%s 🐌 %s %s (%v seconds)\r\n", strings.Repeat(" ", indentLevel), "", test.Name, test.Time)
		}
	}

	if len(group.Tests) > 0 {
		fmt.Println("")
	}

	for _, g := range group.Groups {
		printGroup(g, indentLevel+1)
	}
}
