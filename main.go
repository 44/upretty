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
	currentTestOutput := ""
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
			currentTestOutput = ""
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
			fmt.Print(currentTestOutput)
			if !showFailuresOnly {
				fmt.Println("\n", currentTest, failMessage, time.Since(currentTestStart))
			}
			currentTestOutput = ""
			currentTest = ""
		} else if passingMatches != nil {
			if !showFailuresOnly {
				fmt.Println(successMessage, time.Since(currentTestStart))
			}
			currentTest = ""
			currentTestOutput = ""
			passed += 1
		} else if rerunMatches != nil {
			capturingHelp = true
		} else {
			if capturingHelp {
				helpText = helpText + "\n" + txt
			} else {
				if showFailuresOnly {
					currentTestOutput = currentTestOutput + "\n" + latchPrefix.ReplaceAllString(txt, "")
				} else {
					currentTestOutput = currentTestOutput + "\n" + txt
				}
			}
		}

	}
	if len(failedTests) > 0 {
		if !showFailuresOnly {
			fmt.Println("\n\nFailed tests:")
			for i := range failedTests {
				fmt.Println(failedTests[i])
			}
		}
	}
	if currentTest != "" {
		fmt.Println("\n\nUnfinished tests:")
		fmt.Println(currentTest)
		fmt.Println(currentTestOutput)
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
}
