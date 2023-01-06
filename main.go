package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
	"flag"
)

func main() {
	disableColors := false
	showFailuresOnly := false
	flag.BoolVar(&disableColors, "no-color", false, "disable color output")
	flag.BoolVar(&showFailuresOnly, "failures-only", false, "show only failures")
	flag.Parse()
	passed := 0
	failedTests := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	currentTest := ""
	currentTestOutput := []string{}
	currentTestStart := time.Now()
	entering := regexp.MustCompile(`..?/..?/.... .?.:..:.. .M Entering[^:]*: (?P<Name>.*):`)
	failure := regexp.MustCompile(`..?/..?/.... .?.:..:.. .M Fail[^:]*: (?P<Name>.*)`)
	pass := regexp.MustCompile(`..?/..?/.... .?.:..:.. .M Pass[^:]*: (?P<Name>.*)`)
	latchPrefix := regexp.MustCompile(`^..?/..?/.... .?.:..:.. .M [^:]*: `)
	help := regexp.MustCompile(`To re-run:`)
	capturingHelp := false
	helpText := ""

	failMessage := "\033[31mFail\033[0m"
	if disableColors {
		failMessage = "Fail"
	}
	successMessage := "\033[32mOK\033[0m"
	if disableColors {
		successMessage = "OK"
	}

	for scanner.Scan() {
		txt := scanner.Text()
		enteringMatches := entering.FindStringSubmatch(txt)
		passingMatches := pass.FindStringSubmatch(txt)
		failureMatches := failure.FindStringSubmatch(txt)
		rerunMatches := help.FindStringSubmatch(txt)
		if enteringMatches != nil {
			index := entering.SubexpIndex("Name")
			currentTest = enteringMatches[index]
			currentTestOutput = []string{}
			currentTestStart = time.Now()
			if !showFailuresOnly {
				fmt.Print(currentTest)
			}
		} else if failureMatches != nil {
			if showFailuresOnly {
				fmt.Println()
				fmt.Println(failMessage, currentTest, time.Since(currentTestStart))
			}
			failedTests = append(failedTests, currentTest)
			// fmt.Print(currentTestOutput)
			if len(currentTestOutput) > 0 {
				fmt.Println()
			}
			for _, line := range(currentTestOutput) {
				fmt.Println("    ", line)
			}
			if !showFailuresOnly {
				fmt.Println()
				fmt.Println(currentTest, failMessage, time.Since(currentTestStart))
			}
			currentTestOutput = []string{}
			currentTest = ""
		} else if passingMatches != nil {
			if !showFailuresOnly {
				fmt.Println("", successMessage, time.Since(currentTestStart))
			}
			currentTest = ""
			currentTestOutput = []string{}
			passed += 1
		} else if rerunMatches != nil {
			capturingHelp = true
		} else {
			if capturingHelp {
				helpText = helpText + "\n" + txt
			} else {
				if showFailuresOnly {
					// currentTestOutput = currentTestOutput + "\n" + latchPrefix.ReplaceAllString(txt, "")
					currentTestOutput = append(currentTestOutput, latchPrefix.ReplaceAllString(txt, ""))
				} else {
					currentTestOutput = append(currentTestOutput, txt)
				}
			}
		}

	}
	exitCode := 0
	if len(failedTests) > 0 {
		if !showFailuresOnly {
			fmt.Println("\n\nFailed tests:")
			for i := range failedTests {
				fmt.Println(failedTests[i])
			}
		}
		exitCode = 1
	}
	if currentTest != "" {
		fmt.Println("\n\nUnfinished tests:")
		fmt.Println(currentTest)
		fmt.Println(currentTestOutput)
		exitCode = 1
	}
	if !showFailuresOnly {
		fmt.Printf("\nTOTAL: %d passed / %d failed\n", passed, len(failedTests))

		if helpText != "" {
			fmt.Println("\nTo re-run:")
			fmt.Println(helpText)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	os.Exit(exitCode)
}
