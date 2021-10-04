package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	fmt.Println("FlowCat")

	folderFlag := flag.String("f", ".", "The project top level directory, where flowcat should start recursing from.")
	outputFlag := flag.String("o", "", "Optional output file to dump results to, note output will still be shown on terminal.")
	matchFlag := flag.String("m", "//@todo", "The string to match to do items on.")
	lineFlag := flag.Bool("l", false, "If line numbers should be shown with todo items in output.")
	helpFlag := flag.Bool("h", false, "Shows the help menu.")
	flag.Parse()

	//@todo remove below after testing
	fmt.Println("Folder: ", *folderFlag)
	fmt.Println("Output: ", *outputFlag)
	fmt.Println("Match: ", *matchFlag)
	fmt.Println("Lines: ", *lineFlag)
	fmt.Println("Help: ", *helpFlag)

	//Helpflag implemented because the default help flag from the flag package returns status code 2
	if *helpFlag {
		fmt.Println("Usage of Flowcat:")
		fmt.Println("-f string")
		fmt.Println("    The project top level directory, where flowcat should start recursing from. (default '.' Current Directory)")
		fmt.Println("-l    If line numbers should be shown with todo items in output.")
		fmt.Println("-m string")
		fmt.Println("    The string to match to do items on. (default '//@todo')")
		fmt.Println("-o string")
		fmt.Println("    Optional output file to dump results to, note output will still be shown on terminal.")
	}

	getParsableFiles := func(path string, info os.FileInfo, _ error) (err error) {
		var f []string
		f = append(f, "^\\..*")
		f = append(f, "todo")
		if info.Mode().IsRegular() {
			//@todo later filter differently below is temporary
			for _, i := range f {
				//fmt.Println("From slice", i)
				v, _ := regexp.Compile(i)
				//fmt.Println("After changed", v)
				regCheck := v.MatchString(info.Name())
				fmt.Println(regCheck)

			}
			fmt.Println(info.Name())
		}
		return nil
	}

	//Start crawling the base directory
	//@todo change below to use filepath.WalkDir instead
	err := filepath.Walk(*folderFlag, getParsableFiles)
	if err != nil {
		fmt.Println(err)
	}
}
