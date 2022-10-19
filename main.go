package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

func main() {
	failedTests := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	currentTest := ""
	currentTestOutput := ""
	currentTestStart := time.Now()
	entering := regexp.MustCompile(`../../.... .:..:.. .M Entering[^:]*: (?P<Name>.*):`)
	failure := regexp.MustCompile(`../../.... .:..:.. .M Fail[^:]*: (?P<Name>.*)`)
	pass := regexp.MustCompile(`../../.... .:..:.. .M Pass[^:]*: (?P<Name>.*)`)

	for scanner.Scan() {
		txt := scanner.Text()
		enteringMatches := entering.FindStringSubmatch(txt)
		passingMatches := pass.FindStringSubmatch(txt)
		failureMatches := failure.FindStringSubmatch(txt)
		if enteringMatches != nil {
			index := entering.SubexpIndex("Name")
			currentTest = enteringMatches[index]
			currentTestOutput = ""
			currentTestStart = time.Now()
			fmt.Print(currentTest)
		} else if failureMatches != nil {
			failedTests = append(failedTests, currentTest)
			fmt.Print(currentTestOutput)
			fmt.Println("\n", currentTest, " \033[31mFail\033[0m", time.Since(currentTestStart))
			currentTestOutput = ""
			currentTest = ""
		} else if passingMatches != nil {
			fmt.Println(" \033[32mOK\033[0m", time.Since(currentTestStart))
			currentTest = ""
			currentTestOutput = ""
		} else {
			currentTestOutput = currentTestOutput + "\n" + txt
		}

	}
	if len(failedTests) > 0 {
		fmt.Println("\n\nFailed tests:")
		for i := range failedTests {
			fmt.Println(failedTests[i])
		}
	}
	if currentTest != "" {
		fmt.Println("\n\nUnfinished tests:")
		fmt.Println(currentTest)
		fmt.Println(currentTestOutput)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
