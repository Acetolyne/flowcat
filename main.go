package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	//"path"

	//"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	lexer "github.com/Acetolyne/commentlex"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Match        string              `yaml:"match"`
	IgnoredItems map[string][]string `yaml:"ignore"`
}

//@todo update master branch build badges
//@todo add unit testing
var ListedFiles []string
var Cfg Config

func checkExclude(path string, outfile string, folderFlag string) (string, bool) {
	//@todo add builds to autorun in VSCode
	//@todo make matching workflows for each build on Github to show status for each arch
	regpath := strings.TrimPrefix(path, folderFlag)
	m := Cfg.IgnoredItems["ignore"]
	reg := []bool{}
	//If we are outputting to a file ignore the output file by default
	//@todo does the output file get ignored with the new path settings?
	if outfile != "" {
		m = append(m, outfile)
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

//@todo finish init create file with standard settings if not found after input from user
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
		SetFile.WriteString("# Settings\n")
		SetFile.WriteString("match: \"@todo\"\n\n")
		SetFile.WriteString("# File patterns to ignore\n")
		SetFile.WriteString("ignore:\n")
		SetFile.WriteString("  - \"^\\\\..*\"\n")
		SetFile.Close()
		fmt.Println("Settings file created at ~/.flowcat")
		return nil
	} else {
		return errors.New("setting file already exists consider editing the .flowcat file or delete it before running init")
	}
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
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
	//@todo update help menu and options
	if *helpFlag {
		fmt.Println("Flowcat version 3.0.0")
		fmt.Println("")
		fmt.Println("Options for Flowcat:")
		fmt.Println("init")
		fmt.Println("using flowcat init allows you to create a settings file used when flowcat is run in the current directory, settings can be changed later in the .flowcat file")
		fmt.Println("-f string")
		fmt.Println("   The project top level directory, where flowcat should start recursing from. (default '.' Current Directory)")
		fmt.Println("-l")
		fmt.Println("	If line numbers should be shown with todo items in output.")
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
	settings, err := ioutil.ReadFile(dirname + "/.flowcat")
	// //If there is a settings file then get the values
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		_ = yaml.Unmarshal(settings, &Cfg)
		// 	//Ignore errors
		_ = yaml.Unmarshal(settings, &Cfg.IgnoredItems)
		// 	//Ignore errors
	}

	if *lineFlag {
		Showlines = *lineFlag
	}
	//Cfg.IgnoredItems["ignore"]
	matchexp = Cfg.Match
	if *matchFlag != "" {
		matchexp = *matchFlag
	}
	// //Fallback incase settings is missing the match and one is not specified per an argument
	if matchexp == "" {
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
					var s lexer.Scanner
					s.Match = matchexp
					s.Error = func(*lexer.Scanner, string) {} // ignore errors
					s.Init(file)
					s.Mode = lexer.ScanComments

					checklines := func(s lexer.Scanner, path string, Showlines bool) string {
						tok := s.Scan()
						var line string
						for tok != lexer.EOF {
							// fmt.Println(":)", s.TokenText())
							// fmt.Println("tok:", tok)
							// fmt.Println("Comment?", lexer.Comment)
							if tok == lexer.Comment {
								//remove newlines
								linetext := strings.Replace(s.TokenText(), "\n", " ", -1)
								linetext = strings.Replace(linetext, "\t", " ", -1)
								if Showlines {
									line += "\t" + strconv.Itoa(s.Position.Line) + ")" + linetext + "\n"
								} else {
									line += "\t" + linetext + "\n"
								}
							}
							tok = s.Scan()
						}

						return line
					}
					filelines := checklines(s, path, Showlines)
					if filelines != "" {
						if *outputFlag != "" {
							F.WriteString(path)
							F.WriteString("\n")
							F.WriteString(filelines)
						}
						fmt.Println(path)
						fmt.Println(filelines)
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
}
