package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

func main() {
	fmt.Println("\n \n \n")

	//gather command-line arguments except name of the programm
	args := os.Args[1:]

	//break the program if number of arguments is different from what is needed

	if len(args) != 2 {
		fmt.Println("missing arguments", len(args), args)
		return
	}

	/* 	var wg sync.WaitGroup
	 */
	/* 	wg.Add(1)
	go func() {

		filesSplitter(file.Name(), args[1], &wg)
	}() */
	readCh := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		arrayOfFiles(args[0], readCh)

		defer close(readCh)
		defer fmt.Println("channel closed")
		defer wg.Done()
	}()
	var i int
	for value := range readCh {
		fmt.Println("hjere  ____", value, i)
		i++
	}
	wg.Wait()

	fmt.Println("\n \n \n")

}

func filesSplitter(file string, match string, wg *sync.WaitGroup) {
	content, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}
	defer wg.Done()

	arrayOfLines := strings.Split(string(content), "\n")

	for k, v := range arrayOfLines {
		matchFinder(k, v, match, file)
	}

}

func matchFinder(lineno int, line string, match string, filename string) {
	matched, _ := regexp.MatchString(match, line)
	if matched {
		fmt.Println(matched, filename, lineno, line)
	}
}

func arrayOfFiles(directory string, readCh chan string) {
	files, err := ioutil.ReadDir(directory)
	fmt.Println(directory, "to be serarched")
	if err != nil {
		log.Fatal(err)
	}

	for k, file := range files {
		fmt.Println("-------", file.Name(), file.IsDir())
		if !file.IsDir() {
			directory = strings.Replace(directory, "//", "/", 1) + "/" + file.Name()
			readCh <- directory
			fmt.Println(k, directory)
			break
		} else {
			a := fmt.Sprintf("%s/%s", directory, file.Name())
			arrayOfFiles(a, readCh)
		}
	}

}
