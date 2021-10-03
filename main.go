package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Println("FlowCat")

	folderFlag := flag.String("-f", ".", "The project top level directory, where flowcat should start recursing from.")
	outputFlag := flag.String("-o", "", "Optional output file to dump results to, note output will still be shown on terminal.")
	matchFlag := flag.String("-m", "//@todo", "The string to match to do items on, defaults to //@todo")

	fmt.Println("Folder: " + *folderFlag)
	fmt.Println("Output: " + *outputFlag)
	fmt.Println("Match: " + *matchFlag)

}
