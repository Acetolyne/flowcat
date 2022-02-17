package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	lexer "github.com/Acetolyne/commentlex"

	"gopkg.in/yaml.v2"
)

type Config struct {
	IgnoredItems map[string][]string `yaml:"ignore"`
}

//@todo update master branch build badges
//@todo add unit testing
var ListedFiles []string
var Cfg Config

func checkExclude(path string, outfile string) (string, bool) {
	//@todo add builds to autorun in VSCode
	//@todo make matching workflows for each build on Github to show status for each arch
	m := Cfg.IgnoredItems["ignore"]
	//If we are outputting to a file ignore the output file by default
	if outfile != "" {
		m = append(m, outfile)
	}
	for _, i := range m {
		v, _ := regexp.Compile(i)
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		contentType, err := GetFileContentType(f)
		if err != nil {
			panic(err)
		}
		regCheck := v.MatchString(strings.TrimSpace(path))
		if regCheck {
			if strings.Contains(contentType, "utf-8") {
				return path, true
			} else {
				return path, false
			}
		} else {
			return path, false
		}
	}

	return path, false
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
		fmt.Println("Regex to match on by default for this project(ex: @todo):")
		fmt.Scanln(&user_match)
		//write the settings to the file
		SetFile, err = os.OpenFile(".flowcat", os.O_WRONLY|io.SeekStart|os.O_CREATE, 0755)
		if err != nil {
			return errors.New("ERROR: could not create settings file")
		}
		//@todo where should this file be stored having it per directory is a bit messy for projects as we may run it specifying different directories
		SetFile.WriteString("# Settings\n")
		// SetFile.WriteString("linenum: \"" + user_lines + "\"\n")
		// SetFile.WriteString("match: \"" + user_match + "\"\n\n")
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

	//If the -f argument contains a specific file then return the directory its in
	dir, _ := path.Split(*folderFlag)
	//Get settings if there is a settings file in the current directory
	settings, err := ioutil.ReadFile(dir + "/.flowcat")
	//If there is a settings file then get the values
	if err == nil {
		_ = yaml.Unmarshal(settings, &Cfg)
		//Ignore errors
		_ = yaml.Unmarshal(settings, &Cfg.IgnoredItems)
		//Ignore errors
		// Showlines, err = strconv.ParseBool(Cfg.Linenums)
		// //@todo should this block? else what should the default be
		// if err != nil {
		// 	fmt.Println("linenum should be true or false", err)
		// }
	}

	if *lineFlag {
		Showlines = *lineFlag
	}
	// matchexp = Cfg.Match
	if *matchFlag != "" {
		matchexp = *matchFlag
	}
	// //Fallback incase settings is missing the match and one is not specified per an argument
	if matchexp == "" {
		matchexp = "@todo"
	}
	//reg, _ := regexp.Compile(matchexp)

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

			file, exc := checkExclude(path, *outputFlag)

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
					//fmt.Println(s.srcType)
					s.Mode = lexer.ScanComments

					checklines := func(s lexer.Scanner, path string, Showlines bool) string {
						tok := s.Scan()
						var line string
						//fmt.Println(path)
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
