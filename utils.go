package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getInput(EmptyAllowed bool, SingleWorded bool, RepeatMessage string) string {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	if SingleWorded {
		ans := strings.Fields(input.Text())
		if len(ans) == 0 {
			if EmptyAllowed {
				return ""
			}
			fmt.Println("This cannot be empty")
			_, _ = cyan.Print(RepeatMessage)
			return getInput(false, true, RepeatMessage)
		}
		if len(ans) > 1 {
			fmt.Println("More than one words entered, only first word will be used as name")
		}
		return ans[0]
	}
	if !EmptyAllowed && input.Text() == "" {
		fmt.Println("This cannot be empty")
		_, _ = cyan.Print(RepeatMessage)
		return getInput(false, false, RepeatMessage)
	}
	return input.Text()
}

func getConsent(Default bool) bool {
	consent := getInput(true, true, "")
	if consent == "" {
		return Default
	}
	return strings.ToLower(string([]rune(consent)[0])) == "y"
}

func getInt(EmptyAllowed bool, RepeatMessage string) int {
	input, err := strconv.Atoi(getInput(EmptyAllowed, true, RepeatMessage))
	if err != nil {
		color.Red("Please enter a number")
		return getInt(EmptyAllowed, RepeatMessage)
	}
	return input
}

func inRange(value int, span []int) bool {
	for _, i := range span {
		if value == i {
			return true
		}
	}
	return false
}

func testWritePermissions() {
	newFile, err := os.Create("nginxAutoConfig.test.txt")
	if err != nil {
		if os.IsPermission(err) {
			log.Println("Error: Write permission denied, please go to a workable dir, exiting.")
			os.Exit(1)
		}
		println(err.Error())
		os.Exit(1)
	} else {
		_ = newFile.Close()
		err := os.Remove("nginxAutoConfig.test.txt")
		if err != nil {
			println(err.Error())
		}
	}
}

func writeContentToFile(fileName string, fileContents string) {
	testWritePermissions()
	err := ioutil.WriteFile(fileName, []byte(fileContents), 0644)
	if err != nil {
		fmt.Println("Something went wrong, please send the log below to Sid Sun.")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Config written to %s, move it to the appropriate config folder and reload the nginx webserver, Enjoy!\n", fileName)
}

func printCautionSSL() {
	_, _ = yellow.Println("Caution: SSL config is commented out by default, please generate the key and point to it correctly as necessary.")
}

func dirExists(path string) bool {
	// Stat the file
	info, err := os.Stat(path)

	// If there is error check if the file is non existent
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	// Else check if the path is actually a directory and return accordingly
	return info.IsDir()
}

func verifyDirInput() string {
	_, _ = cyan.Print("Root path: ")
	dirName := getInput(false, false, "Root path: ")

	if !dirExists(dirName) {
		_, _ = red.Printf("Directory '%v' is non existent, please try again.\n", dirName)
		_, _ = cyan.Print("Root path: ")
		return verifyDirInput()
	}

	// Return a absloute path
	absPath, _ := filepath.Abs(dirName)
	return absPath
}
