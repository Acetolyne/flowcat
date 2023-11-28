package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/Acetolyne/flowcat/lexer"
	"gopkg.in/yaml.v2"
)

// Config structure of configuration from a yaml settings file.
type Config struct {
	Match        string              `yaml:"match"`
	IgnoredItems map[string][]string `yaml:"ignore"`
}

// ListedFiles returns a string of all files in a directory.
var ListedFiles []string

// Cfg returns the user configurations from a file.
var Cfg Config

func checkExclude(path string, outfile string, folderFlag string) (string, bool) {
	regpath := strings.TrimPrefix(path, folderFlag)
	m := Cfg.IgnoredItems["ignore"]
	reg := []bool{}
	//If we are outputting to a file ignore the output file by default if it is in the project path
	if outfile != "" {
		if strings.Contains(outfile, folderFlag) {
			m = append(m, outfile)
		}
	}
	for _, i := range m {
		v, _ := regexp.Compile(i)
		regCheck := v.MatchString(strings.TrimSpace(regpath))
		if regCheck {
			reg = append(reg, true)
		} else {
			reg = append(reg, false)
		}
	}
	for _, r := range reg {
		if r {
			return path, true
		}
	}
	return path, false
}

func initSettings() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Open(dirname + "/.flowcat")
	f.Close()
	if err != nil {
		var SetFile *os.File
		SetFile, err = os.OpenFile(dirname+"/.flowcat", os.O_WRONLY|io.SeekStart|os.O_CREATE, 0755)
		if err != nil {
			return errors.New("ERROR: could not create settings file")
		}
		_, err = SetFile.WriteString("# Settings\n")
		if err != nil {
			return errors.New("ERROR: could not create settings file")
		}
		_, err = SetFile.WriteString("match: \"@todo\"\n\n")
		if err != nil {
			return errors.New("ERROR: could not create settings file")
		}
		_, err = SetFile.WriteString("# File patterns to ignore\n")
		if err != nil {
			return errors.New("ERROR: could not create settings file")
		}
		_, err = SetFile.WriteString("ignore:\n")
		if err != nil {
			return errors.New("ERROR: could not create settings file")
		}
		_, err = SetFile.WriteString("  - \"^\\\\..*\"\n")
		if err != nil {
			return errors.New("ERROR: could not create settings file")
		}
		SetFile.Close()
		fmt.Println("Settings file created at ~/.flowcat")
		return nil
	}
	return errors.New("setting file already exists consider editing the .flowcat file or delete it before running init")
}

func main() {
	var F *os.File
	var Showlines bool = false
	var matchexp string

	folderFlag := flag.String("f", "./", "The project top level directory, where flowcat should start recursing from.")
	outputFlag := flag.String("o", "", "Optional output file to dump results to, note output will still be shown on terminal.")
	matchFlag := flag.String("m", "", "The string to match to do items on.")
	lineFlag := flag.Bool("l", false, "If line numbers should be shown with todo items in output.")
	helpFlag := flag.Bool("h", false, "Shows the help menu.")
	flag.Parse()

	//Helpflag implemented because the default help flag from the flag package returns status code 2
	if *helpFlag {
		fmt.Println("Flowcat version 3.1.1")
		fmt.Println("")
		fmt.Println("Options for Flowcat:")
		fmt.Println("init")
		fmt.Println("using flowcat init creates a settings file for the current user, settings can be changed later in the ~/.flowcat file")
		fmt.Println("-f string")
		fmt.Println("   The project top level directory, where flowcat should start recursing from. (default Current Directory)")
		fmt.Println("-l")
		fmt.Println("	Display line numbers in the output.")
		fmt.Println("-m string")
		fmt.Println("   The string to match to do items on. (default '@todo')")
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

	//Get settings from .flowcat file in users home directory
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	settings, err := os.OpenFile(dirname+"/.flowcat", os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println("Error opening configuration file")
	}
	defer settings.Close()
	configuration := yaml.NewDecoder(settings)
	err = configuration.Decode(&Cfg)
	if err != nil {
		fmt.Println("Unable to get settings from configuration file.")
	}
	// if err == nil {
	// 	_ = yaml.Unmarshal(settings, &Cfg)
	// 	// 	//Ignore errors
	// 	_ = yaml.Unmarshal(settings, &Cfg.IgnoredItems
	// 	// 	//Ignore errors
	// }

	if *lineFlag {
		Showlines = *lineFlag
	}

	matchexp = Cfg.Match
	if *matchFlag != "" {
		matchexp = *matchFlag
	} else {
		matchexp = "@todo"
	}

	parseFiles := func(path string, info os.FileInfo, _ error) (err error) {

		if *outputFlag != "" {
			F, err = os.OpenFile(*outputFlag, os.O_WRONLY|io.SeekStart|os.O_CREATE, 0755)
			if err != nil {
				panic(err)
			}
			defer F.Close()
			if err != nil && F != nil {
				fmt.Println("ERROR: could not write output to", *outputFlag)
			}
		}
		if info.Mode().IsRegular() {
			file, exc := checkExclude(path, *outputFlag, *folderFlag)

			//If the file does not match our exclusion regex then use it.
			if !exc {
				contents, err := os.ReadFile(path)
				if err != nil {
					fmt.Println("ERROR: could not read file", file, err)
				}
				contentbytes := []byte(contents)
				if utf8.Valid(contentbytes) {
					lexer.GetComments(contentbytes, matchexp)
					// var s lexer.Scanner
					// s.Match = matchexp
					// s.Error = func(*lexer.Scanner, string) {} // ignore errors
					// s.Init(file)
					// s.Mode = lexer.ScanComments

					// checklines := func(s lexer.Scanner, path string, Showlines bool) string {
					// 	tok := s.Scan()
					// 	var line string
					// 	for tok != lexer.EOF {
					// 		if tok == lexer.Comment {
					// 			//remove newlines
					// 			linetext := strings.Replace(s.TokenText(), "\n", " ", -1)
					// 			linetext = strings.Replace(linetext, "\t", " ", -1)
					// 			if Showlines {
					// 				line += "\t" + strconv.Itoa(s.Position.Line) + ")" + linetext + "\n"
					// 			} else {
					// 				line += "\t" + linetext + "\n"
					// 			}
					// 		}
					// 		tok = s.Scan()
					// 	}

					// 	return line
					// }
					// filelines := checklines(s, path, Showlines)
					// if filelines != "" {
					// 	if *outputFlag != "" {
					// 		F.WriteString(path)
					// 		F.WriteString("\n")
					// 		F.WriteString(filelines)
					// 	}
					// 	fmt.Println(path)
					// 	fmt.Println(filelines)
					// }
					//TEMP
					if Showlines {
						fmt.Println("showlines")
					}
				}
			}
			return nil
		}
		return nil
	}
	//Start crawling the base directory
	//@todo change below to use filepath.WalkDir instead
	err = filepath.Walk(*folderFlag, parseFiles)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(0)
}
