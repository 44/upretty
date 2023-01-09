package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
	"flag"
	"strings"
	"github.com/charmbracelet/lipgloss"
	// sanitize "github.com/mrz1836/go-sanitize"
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
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff3311"))
	successStyle :=  lipgloss.NewStyle().Foreground(lipgloss.Color("#11aa33"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#777"))
	if disableColors {
		errorStyle = lipgloss.NewStyle()
		successStyle = lipgloss.NewStyle()
		dimStyle = lipgloss.NewStyle()
	}

	failMessage := errorStyle.Render("Fail") //"\033[31mFail\033[0m"
	successMessage := successStyle.Render("OK") //"\033[32mOK\033[0m"

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
				fmt.Print(currentTest, ':')
			}
		} else if failureMatches != nil {
			if showFailuresOnly {
				fmt.Println()
				fmt.Println(failMessage, currentTest, dimStyle.Render(time.Since(currentTestStart).String()))
			}
			failedTests = append(failedTests, currentTest)
			fmt.Print("\x1b[2K")
			if len(currentTestOutput) > 0 {
				fmt.Println()
			}
			for _, line := range(currentTestOutput) {
				fmt.Println("    ", line)
			}
			if !showFailuresOnly {
				fmt.Println()
				fmt.Println(currentTest, failMessage, dimStyle.Render(time.Since(currentTestStart).String()))
			}
			currentTestOutput = []string{}
			currentTest = ""
		} else if passingMatches != nil {
			if !showFailuresOnly {
				fmt.Print("\x1b[2K\x1b[0G")
				fmt.Println(currentTest, successMessage, dimStyle.Render(time.Since(currentTestStart).String()))
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
					// fmt.Print()
					if currentTest != "" {
						fmt.Print("\x1b[2K\x1b[0G")
						s := strings.ReplaceAll(txt, "\n", " ")
						s = strings.ReplaceAll(txt, "\r", " ")
						s = strings.ReplaceAll(txt, "\t", " ")
						s = strings.ReplaceAll(txt, "  ", " ")
						s = strings.Trim(s, " \n\r\t")
						fmt.Printf("%s: %80.80s", currentTest, s)
					}
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
		fmt.Println()
		passedStyle := successStyle
		failedStyle := errorStyle
		if passed == 0 {
			passedStyle = dimStyle
		}
		if len(failedTests) == 0 {
			failedStyle = dimStyle
		}
		fmt.Println("TOTAL: ",
			passedStyle.Render(fmt.Sprintf("%d passed", passed)),
			"/",
			failedStyle.Render(fmt.Sprintf("%d failed", len(failedTests))))

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
