package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const (
	tab        = "\t"
	trait      = "│\t"
	lastFile   = "└───"
	middleFile = "├───"
)

func getFilteredFiles(files []os.FileInfo, isPrintFile bool) []os.FileInfo {
	filteredFiles := make([]os.FileInfo, 0)
	for _, value := range files {
		if value.IsDir() || isPrintFile {
			filteredFiles = append(filteredFiles, value)
		}
	}
	return filteredFiles
}

func printDirectories(output io.Writer, path string, isPrintFiles bool, line string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	files = getFilteredFiles(files, isPrintFiles)

	sliceSize := len(files)
	for index, value := range files {
			currentLine := line
			if index == sliceSize - 1 {
				currentLine += lastFile
			} else {
				currentLine += middleFile
			}

			currentLine += value.Name()
			if !value.IsDir() {
				fileSize := value.Size()
				if fileSize != 0 {
					currentLine += " (" + strconv.FormatInt(value.Size(), 10) + "b)"
				} else {
					currentLine += " (" + "empty" + ")"
				}
			}

			_, err := fmt.Fprint(output, currentLine, "\n")
			if err != nil {
				return err
			}

			if value.IsDir() {
			nextLine := line
			if index == sliceSize - 1 {
				nextLine += tab
			} else {
				nextLine += trait
			}
			dirPath := fmt.Sprintf("%s/%s", path, value.Name())
			err = printDirectories(output, dirPath, isPrintFiles, nextLine)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func dirTree(output io.Writer, path string, isPrintFiles bool) error {
	return printDirectories(output, path, isPrintFiles, "")
}

func main() {
	out := new(bytes.Buffer)
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
