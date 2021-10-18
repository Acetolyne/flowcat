package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Linenums     string              `yaml:"linenum"`
	Match        string              `yaml:"match"`
	IgnoredItems map[string][]string `yaml:"ignore"`
}

type Ignored struct {
	Name       string
	Properties []string
}

//@todo push binaries from github workflow
//@todo update master branch build badges
//@todo add unit testing
//@todo add github workflows for testing binaries on different OS's
var ListedFiles []string
var cfg Config

func testExclude(path string, outfile string, cfg Config) (string, bool) {
	//@todo add builds to autorun in VSCode
	//@todo make matching workflows for each build on Github to show status for each arch
	m := cfg.IgnoredItems["ignore"]
	//If we are outputting to a file ignore the output file by default
	if outfile != "" {
		m = append(m, outfile)
	}
	for _, i := range m {
		v, _ := regexp.Compile(i)
		regCheck := v.MatchString(strings.TrimSpace(path))
		if regCheck {
			return path, true
		}

	}
	return path, false
}

func testLine(line string, flag string) bool {
	//flag = "^" + flag
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

//@todo finish init create file with standard settings if not found after input from user
func initSettings() error {
	var user_lines string
	var user_match string
	f, err := os.Open(".flowcat")
	f.Close()
	if err != nil {
		var SetFile *os.File
		fmt.Println("Display line numbers by default(true/false):")
		fmt.Scanln(&user_lines)
		fmt.Println("Regex to match on by default for this project(ex: //@todo):")
		fmt.Scanln(&user_match)
		//write the settings to the file
		SetFile, err = os.OpenFile(".flowcat", os.O_WRONLY|io.SeekStart|os.O_CREATE, 0755)
		if err != nil {
			return errors.New("ERROR: could not create settings file")
		}
		SetFile.WriteString("# Settings\n")
		SetFile.WriteString("linenum: \"" + user_lines + "\"\n")
		SetFile.WriteString("match: \"" + user_match + "\"\n\n")
		SetFile.WriteString("# File patterns to ignore\n")
		SetFile.WriteString("ignore:\n")
		SetFile.WriteString("  - \"\\\\.flowcat\"\n")
		SetFile.WriteString("  - \"^\\\\.\"\n")
		SetFile.Close()
		return nil
	} else {
		return errors.New("setting file already exists consider editing the .flowcat file or delete it before running init")
	}
}

func main() {
	var F *os.File
	var showlines bool = false
	var matchexp string

	folderFlag := flag.String("f", ".", "The project top level directory, where flowcat should start recursing from.")
	outputFlag := flag.String("o", "", "Optional output file to dump results to, note output will still be shown on terminal.")
	matchFlag := flag.String("m", "", "The string to match to do items on.")
	lineFlag := flag.Bool("l", false, "If line numbers should be shown with todo items in output.")
	helpFlag := flag.Bool("h", false, "Shows the help menu.")
	flag.Parse()

	//Helpflag implemented because the default help flag from the flag package returns status code 2
	if *helpFlag {
		fmt.Println("Flowcat version 2.0.0")
		fmt.Println("")
		fmt.Println("Options for Flowcat:")
		fmt.Println("init")
		fmt.Println("using flowcat init allows you to create a settings file used when flowcat is run in the current directory, settings can be changed later in the .flowcat file")
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
	//if we are using the init argument then run init function
	if len(os.Args) > 1 && os.Args[1] == "init" {
		err := initSettings()
		if err != nil {
			fmt.Println(err)
		}
		//always exit without running if we were using init argument
		os.Exit(0)
	}

	//Get settings if there is a settings file in the current directory
	settings, err := ioutil.ReadFile(".flowcat")
	if err == nil {
		showlines, err = strconv.ParseBool(cfg.Linenums)
		//@todo should this block? else what should the default be
		if err != nil {
			fmt.Println("linenum should be true or false", err)
		}
	}
	err = yaml.Unmarshal(settings, &cfg)
	err = yaml.Unmarshal(settings, &cfg.IgnoredItems)

	if *lineFlag != false {
		showlines = *lineFlag
	}
	matchexp = cfg.Match
	fmt.Println(matchexp)
	if *matchFlag != "" {
		matchexp = *matchFlag
	}
	//Fallback incase settings is missing the match and one is not specified per an argument
	if matchexp == "" {
		matchexp = "//@todo"
	}

	reg, _ := regexp.Compile(matchexp)

	parseFiles := func(path string, info os.FileInfo, _ error) (err error) {
		if *outputFlag != "" {
			F, err = os.OpenFile(*outputFlag, os.O_WRONLY|io.SeekStart|os.O_CREATE, 0755)
			defer F.Close()
			if err != nil && F != nil {
				fmt.Println("ERROR: could not write output to", *outputFlag)
			}
		}
		if info.Mode().IsRegular() {
			file, exc := testExclude(path, *outputFlag, cfg)

			//If the file does not match our exclusion regex then use it.
			if !exc {
				curfile, err := os.Open(file)
				if err == nil {
					fscanner := bufio.NewScanner(curfile)
					var linenum = 0
					var ln string
					for fscanner.Scan() {
						if showlines {
							linenum++
							ln = fmt.Sprint(linenum)
							ln = ln + ")"
						}
						incline := testLine(fscanner.Text(), matchexp)
						if incline {
							listFile(path, F)
							l := "\t" + ln + reg.Split(strings.TrimSpace(fscanner.Text()), 2)[1]
							fmt.Println("\t", ln, reg.Split(strings.TrimSpace(fscanner.Text()), 2)[1])
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
	err = filepath.Walk(*folderFlag, parseFiles)
	if err != nil {
		fmt.Println(err)
	}
}
