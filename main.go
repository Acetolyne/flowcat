package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ListedFiles []string

func testExclude(path string) (string, bool) {
	var f []string
	f = append(f, "^\\.")
	f = append(f, "todo")
	for _, i := range f {
		//fmt.Println("From slice", i)
		v, _ := regexp.Compile(i)
		//fmt.Println("After changed", v)
		regCheck := v.MatchString(strings.TrimSpace(path))
		//fmt.Println(v, regCheck)
		if regCheck {
			return path, true
		}

	}
	return path, false
}

func testLine(line string, flag string) bool {
	//fmt.Println("From slice", i)
	flag = "^" + flag
	v, _ := regexp.Compile(flag)
	//fmt.Println("After changed", v)
	regCheck := v.MatchString(strings.TrimSpace(line))
	//fmt.Println(v, regCheck)
	if regCheck {
		return true
	}
	return false
}

func listFile(file string) bool {
	for _, v := range ListedFiles {
		if v == file {
			return false
		}
	}
	ListedFiles = append(ListedFiles, file)
	fmt.Println(file)
	return true
}

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

	parseFiles := func(path string, info os.FileInfo, _ error) (err error) {
		if info.Mode().IsRegular() {
			file, exc := testExclude(path)
			//If the file does not match our exclusion regex then use it.
			if !exc {
				curfile, err := os.Open(file)
				if err == nil {
					fscanner := bufio.NewScanner(curfile)
					var linenum = 0
					var ln string
					for fscanner.Scan() {
						if *lineFlag {
							linenum++
							ln = fmt.Sprint(linenum)
							ln = ln + ")"
						}
						incline := testLine(fscanner.Text(), *matchFlag)
						if incline {
							listFile(path)
							fmt.Println("\t", ln, strings.TrimSpace(fscanner.Text()))
						}
					}
				}
			}
		}
		return nil
	}

	//Start crawling the base directory
	//@todo change below to use filepath.WalkDir instead
	err := filepath.Walk(*folderFlag, parseFiles)
	if err != nil {
		fmt.Println(err)
	}
}
