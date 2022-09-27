package main

import (
	"fmt"
	"mgrep/workers"
	"os"
	"sync"
	"time"
)

func main() {
	times := time.Now()
	var numberOfMatches int

	//gather command-line arguments except name of the programm
	args := os.Args[1:]

	//break the program if number of arguments is different from what is needed

	if len(args) != 2 {
		fmt.Println("missing arguments - have", len(args), "- wants-", args)
		return
	}

	//readCh is the Channel used to communicate & synchronize between the directory scanner (ListFiles()) and the file reader function (LinesSplitter())
	readCh := make(chan string)

	//Lines is the channel used to push and pull the info of the Lines (with filename and line number) to be compared by the MatchFinder function
	Lines := make(chan workers.Line)

	var filesToRead, lines, errors int // Miscelanous Counters. filesToRead tracks the number of files to Read (not counting files that had Error messages when opening), lines counts the total number of lines to be compared with the matching string, errors track the number of Errors found when opening files.

	var wg, ng, nng sync.WaitGroup // Waitgroups for the the three main axis of GoRoutines (ListFiles(wg) in a Directory, SplitLines(ng) from a File into a array of type Line, MatchFinder(nng))

	wg.Add(1)
	go func() {
		defer wg.Done()
		workers.ListFiles(args[0], readCh, &wg, &filesToRead, &errors)

	}() //Launches initial goroutine that will scan the directory and subdirectories for files and send them into a channel to be read by the LineSplitter function

	ng.Add(1)
	go func() {
		for elem := range readCh {
			ng.Add(1)
			go func(filePath string) {
				workers.LineSplitter(filePath, Lines, &lines)
				defer ng.Done()
			}(elem)
		}

		defer ng.Done()
	}() //Pulls filenames and paths from the readCh directory through the range syntax, launching a new goroutine for each file.

	var goLines int

	nng.Add(1)
	go func() {
		for elem := range Lines {
			goLines++
			nng.Add(1)
			go func(s workers.Line) {
				workers.MatchFinder(s, args[1], &numberOfMatches)
				defer nng.Done()
			}(elem)
		}
		defer nng.Done()
	}() // Pulls variables of type Line (which contains the int line number(lineno), string text to be compared with the match (line), and the string name of the file (filename)) from the Lines Channel and launches a new goroutine with the function MatchFinder in package workers.go to examine for matches of the matching string within such Line.

	wg.Wait()
	close(readCh) // waits on ListFiles to be executed and closes the readCh when done.
	ng.Wait()
	close(Lines) // waits on LineSplitter to be executed and closes the Line when done.
	nng.Wait()

	since := time.Since(times) // time since program started - tracking for perfomance purposes.

	fmt.Println("time elapsed", since, "no. of Matches found", numberOfMatches, "files", filesToRead, "Lines", goLines, "errors", errors)
}
