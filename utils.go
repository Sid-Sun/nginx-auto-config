package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type inputConfig struct {
	EmptyAllowed  bool
	SingleWorded  bool
	RepeatMessage string
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func readFromFile(filePath string) []byte {
	// Check if file exists and if not, print
	if fileExists(filePath) {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err.Error())
		}
		return data
	}
	fmt.Println("File:", filePath, "seems to be nonexistent")
	os.Exit(0)
	return nil
}

func newInputConfig(EmptyAllowed bool, SingleWorded bool, RepeatMessage string) inputConfig {
	return inputConfig{
		EmptyAllowed:  EmptyAllowed,
		SingleWorded:  SingleWorded,
		RepeatMessage: RepeatMessage,
	}
}

func getRootPath() string {
	inputConfig := newInputConfig(false, false, "Root path: ")
	path := getInput(inputConfig)
	if pathExists, pathInfo := pathExists(path); !pathExists {
		_, _ = yellow.Printf("%s does not exist on this machine, do you want to keep this?\n", path)
		_, _ = cyan.Print("Keep non-existent directory (Y[es]/n[o]): ")
		verifyRoot := !getConsent(true)
		if !verifyRoot {
			return path
		}
		path = verifyDirInput()
	} else if !pathInfo.IsDir() {
		_, _ = yellow.Printf("Path %s is not a directory, do you want to keep this?\n", path)
		_, _ = cyan.Print("Keep non-directory path as root (Y[es]/n[o]): ")
		verifyPath := !getConsent(true)
		if !verifyPath {
			return getAbsolutePath(path)
		}
		path = verifyDirInput()
	}
	return getAbsolutePath(path)
}

func getInput(config inputConfig) string {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	if config.SingleWorded {
		ans := strings.Fields(input.Text())
		if len(ans) == 0 {
			if config.EmptyAllowed {
				return ""
			}
			fmt.Println("This cannot be empty")
			_, _ = cyan.Print(config.RepeatMessage)
			return getInput(config)
		}
		if len(ans) > 1 {
			fmt.Println("More than one words entered, only first word will be used as name")
		}
		return ans[0]
	}
	if !config.EmptyAllowed && input.Text() == "" {
		fmt.Println("This cannot be empty")
		_, _ = cyan.Print(config.RepeatMessage)
		return getInput(config)
	}
	return input.Text()
}

func getConsent(Default bool) bool {
	inputConfig := newInputConfig(true, true, "")
	consent := getInput(inputConfig)
	if consent == "" {
		return Default
	}
	return strings.ToLower(string([]rune(consent)[0])) == "y"
}

func getInt(EmptyAllowed bool, RepeatMessage string) int {
	inputConfig := newInputConfig(EmptyAllowed, true, RepeatMessage)
	input, err := strconv.Atoi(getInput(inputConfig))
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
	}
	_ = newFile.Close()
	err = os.Remove("nginxAutoConfig.test.txt")
	if err != nil {
		println(err.Error())
		os.Exit(0)
	}
}

func writeContentToFile(fileName string, fileContents []byte) error {
	err := ioutil.WriteFile(fileName, fileContents, 0644)
	if err != nil {
		return err
	}
	return nil
}

func printCautionSSL() {
	_, _ = yellow.Println("Caution: SSL config and listen directives are commented out by default, please generate the key, point to it correctly and uncomment them.")
}

func pathExists(path string) (bool, os.FileInfo) {
	var info os.FileInfo
	var err error
	if info, err = os.Stat(path); os.IsNotExist(err) {
		return false, info
	}
	// Thence, the path does exist, return true and info
	return true, info
}

func verifyDirInput() string {
	_, _ = cyan.Print("Root path: ")
	inputConfig := newInputConfig(false, false, "Root path: ")
	dirName := getInput(inputConfig)
	if pathExists, fileInfo := pathExists(dirName); !pathExists {
		_, _ = red.Printf("Path '%v' is non-existent, please try again.\n", dirName)
		_, _ = cyan.Print("Root path: ")
		return verifyDirInput()
	} else if !fileInfo.IsDir() {
		_, _ = red.Printf("Path '%v' is not a directory, please try again.\n", dirName)
		_, _ = cyan.Print("Root path: ")
		return verifyDirInput()
	}
	return dirName
}

func getAbsolutePath(path string) string {
	absPath, _ := filepath.Abs(path)
	return absPath
}
