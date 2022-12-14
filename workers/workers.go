package workers

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

// Used by functions LineSplitter and MatchFinder. It contains the lines of the files opened by LineSplitter to be send to the MatchFinder function for match finding.
type Line struct {
	lineNo   int    // Number of the Line in the file
	line     string // The string contained in the line.
	filePath string // Path of the file where the line string is located
}

// Function ListFiles takes a directory string provided by arg[0], a channel to send the results, a WaitGroup to synchronize with other goRoutines, a pointer to a filesToRead int that tracks the number of files founded in the directory it recursed, and the pointer to errors int that tracks the number of errors found. This functions runs recursively through subdirectories and launches new goroutines for each recursion.
func ListFiles(directory string, readCh chan string, wg *sync.WaitGroup, filesToRead *int, errors *int) {
	files, err := os.ReadDir(directory)

	if err != nil {
		*errors++
		//fmt.Println("ERROR", *errors, directory)
	} else {
		for _, file := range files {
			filePath := filepath.Join(directory, file.Name())
			if !file.IsDir() {
				*filesToRead++
				readCh <- filePath
			} else if file.IsDir() { // if File is a Directory, format the
				wg.Add(1)
				go func() {
					defer wg.Done()
					ListFiles(filePath, readCh, wg, filesToRead, errors)
				}()
			}
		}
	}
}

// LineSplitter takes a filePath string to read a file, receives a chan of type Lines to send the lines in the format of the Line struct to be compared on the MatchFinder function, and counts the total number of lines to be compared during an invocation on the lines pointer to an int variable.
func LineSplitter(filePath string, Lines chan Line, lines *int) {
	readFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var k int
	for fileScanner.Scan() {
		k++
		lineToCompare := Line{
			k,
			fileScanner.Text(),
			filePath,
		}
		Lines <- lineToCompare
	}

	defer readFile.Close()
}

// MatchFinder receives a variable of type Line to be compared, a value match of type string to be compared with, and a pointer to the matches of type to count the number of matches founded during an invocation of the program.
func MatchFinder(elem Line, match string, matches *int) {
	matched, _ := regexp.MatchString(match, elem.line)
	if matched {
		*matches = *matches + 1
		fmt.Println(*matches, "MATCHED @ Line", elem.lineNo, elem.filePath, elem.line)
	}
}
