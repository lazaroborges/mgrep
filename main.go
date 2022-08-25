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

	var readIt []string
	/* 	var wg sync.WaitGroup
	 */
	/* 	wg.Add(1)
	go func() {

		filesSplitter(file.Name(), args[1], &wg)
	}() */
	arrayOfFiles(args[0], &readIt)

	fmt.Println("\n \n \n")

	for k, v := range readIt {
		fmt.Println(k, v)
	}

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

func arrayOfFiles(directory string, returnable *[]string) {
	files, err := ioutil.ReadDir(directory)
	fmt.Println(directory, "to be serarched")
	if err != nil {
		log.Fatal(err)

	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		if !file.IsDir() {
			directory = strings.Replace(directory, "//", "/", 1)
			fmt.Println("____________________", directory)
			*returnable = append(*returnable, directory+file.Name())
		} else {
			a := fmt.Sprintf("%s/%s", directory, file.Name())
			arrayOfFiles(a, returnable)
			fmt.Println("breaking here 2", a)

		}
		fmt.Println(file.Name(), file.IsDir())
	}

}
