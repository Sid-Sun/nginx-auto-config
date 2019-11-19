package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func getInput(EmptyAllowed bool, SingleWorded bool) string {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	if SingleWorded {
		ans := strings.Fields(input.Text())
		if len(ans) == 0 {
			if EmptyAllowed {
				return ""
			}
			fmt.Println("This cannot be empty, please enter some text")
			return getInput(false, true)
		}
		if len(ans) > 1 {
			fmt.Println("More than one words entered, only first word will be used as name")
		}
		return ans[0]
	}
	if !EmptyAllowed && input.Text() == "" {
		fmt.Println("This cannot be empty, please enter some text")
		return getInput(false, false)
	}
	return input.Text()
}

func getConsent() bool {
	consent := string([]rune(getInput(false, true))[:1])
	if strings.ToLower(consent) == "y" {
		return true
	}
	return false
}
func getInt() int {
	input, err := strconv.Atoi(getInput(false, true))
	if err != nil {
		color.Red("Please enter a number")
		return getInt()
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
