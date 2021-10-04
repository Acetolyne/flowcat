package main

import (
	"flag"
	"fmt"
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
		fmt.Println("    The project top level directory, where flowcat should start recursing from. (default '.')")
		fmt.Println("-l    If line numbers should be shown with todo items in output.")
		fmt.Println("-m string")
		fmt.Println("    The string to match to do items on. (default '//@todo')")
		fmt.Println("-o string")
		fmt.Println("    Optional output file to dump results to, note output will still be shown on terminal.")
	}
}
