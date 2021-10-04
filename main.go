package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//@todo add more architectures under bin folder
//@todo update master branch build badge
var ListedFiles []string

func testExclude(path string, outfile string) (string, bool) {
	var f []string
	//@todo make below come from a settings file passed in my argument -e
	//@todo make global and per folder settings file .flowcat
	//@todo add regex info to readme file
	//@todo add builds to autorun in VSCode
	//@todo make matching workflows for each build on Github to show status for each arch?
	//ignore paths starting with period by default
	f = append(f, "^\\.")
	//ignore paths starting with underscore by default
	f = append(f, "^_")
	//If we are outputting to a file ignore the output file by default
	if outfile != "" {
		f = append(f, outfile)
	}
	for _, i := range f {
		v, _ := regexp.Compile(i)
		regCheck := v.MatchString(strings.TrimSpace(path))
		if regCheck {
			return path, true
		}

	}
	return path, false
}

func testLine(line string, flag string) bool {
	flag = "^" + flag
	v, _ := regexp.Compile(flag)
	regCheck := v.MatchString(strings.TrimSpace(line))
	if regCheck {
		return true
	}
	return false
}

func listFile(file string, f *os.File) bool {
	for _, v := range ListedFiles {
		if v == file {
			return false
		}
	}
	ListedFiles = append(ListedFiles, file)
	fmt.Println(file)
	f.WriteString(file + "\n")
	return true
}

func main() {
	var F *os.File

	folderFlag := flag.String("f", ".", "The project top level directory, where flowcat should start recursing from.")
	outputFlag := flag.String("o", "", "Optional output file to dump results to, note output will still be shown on terminal.")
	matchFlag := flag.String("m", "//@todo", "The string to match to do items on.")
	lineFlag := flag.Bool("l", false, "If line numbers should be shown with todo items in output.")
	helpFlag := flag.Bool("h", false, "Shows the help menu.")
	flag.Parse()

	//Helpflag implemented because the default help flag from the flag package returns status code 2
	if *helpFlag {
		fmt.Println("Flowcat version 2.0.0")
		fmt.Println("")
		fmt.Println("Options for Flowcat:")
		fmt.Println("-f string")
		fmt.Println("   The project top level directory, where flowcat should start recursing from. (default '.' Current Directory)")
		fmt.Println("-l")
		fmt.Println("	If line numbers should be shown with todo items in output.")
		fmt.Println("-m string")
		fmt.Println("   The string to match to do items on. (default '//@todo')")
		fmt.Println("-o string")
		fmt.Println("   Optional output file to dump results to, note output will still be shown on terminal.")
		fmt.Println("-h")
		fmt.Println("   This help menu.")
		os.Exit(0)
	}

	parseFiles := func(path string, info os.FileInfo, _ error) (err error) {
		if *outputFlag != "" {
			F, err = os.OpenFile(*outputFlag, os.O_WRONLY|io.SeekStart|os.O_CREATE, 0755)
			defer F.Close()
			if err != nil && F != nil {
				fmt.Println("ERROR: could not write output to", *outputFlag)
			}
		}
		if info.Mode().IsRegular() {
			file, exc := testExclude(path, *outputFlag)
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
							listFile(path, F)
							l := "\t" + ln + strings.TrimSpace(fscanner.Text())
							fmt.Println("\t", ln, strings.TrimSpace(fscanner.Text()))
							if *outputFlag != "" {
								F.WriteString(l + "\n")
							}
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
