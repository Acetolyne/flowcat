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
	"text/scanner"
	"unicode/utf8"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Linenums     string              `yaml:"linenum"`
	Match        string              `yaml:"match"`
	IgnoredItems map[string][]string `yaml:"ignore"`
}

type Comments struct {
	FileExt        []string
	SingleLine     string
	MultiLineStart string
	MultiLineEnd   string
}

//comment types from https://geekflare.com/how-to-add-comments/
//@todo we only need custom file mappings that are not covered in the lexer
var CommentTypes = map[int]Comments{
	0: {
		FileExt:        []string{".go", ".py", ".kt", ".kts", ".ktm"},
		SingleLine:     "//",
		MultiLineStart: "/*",
		MultiLineEnd:   "*/",
	},
	1: {
		FileExt:        []string{".js", ".ts", ".tsx", ".jsx", ".html", ".css", ".scss", ".sass", ".less", ".styl", ".stylus", ".json", ".yml", ".yaml", ".xml", ".toml", ".md", ".sh", ".bat", ".ini", ".conf", ".c", ".cpp", ".h", ".hpp", ".hxx", ".h++", ".cs", ".java", ".csproj", ".csproj.user", ".csproj.references", ".csproj.debug", ".csproj.release", ".csproj.unspecified", ".csproj.nuget", ".csproj.nuget.unspecified", ".csproj.nuget.debug", ".csproj.nuget.release", ".csproj.nuget.references", ".csproj.nuget.user", ".csproj.nuget.unspecified", ".csproj.nuget.debug", ".csproj.nuget.release", ".csproj.nuget.references", ".csproj.nuget.user", ".csproj.nuget.unspecified", ".csproj.nuget.debug", ".csproj.nuget.release", ".csproj.nuget.references", ".csproj.nuget.user", ".csproj.nuget.unspecified", ".csproj.nuget.debug", ".csproj.nuget.release", ".csproj.nuget.references", ".csproj.nuget.user", ".csproj.nuget.unspecified", ".csproj.nuget.debug", ".csproj.nuget.release", ".csproj.nuget.references", ".csproj.nuget.user", ".csproj.nuget.unspecified", ".csproj.nuget.debug", ".csproj.nuget.release", ".csproj.nuget.references", ".csproj.nuget"},
		SingleLine:     "//",
		MultiLineStart: "/*",
		MultiLineEnd:   "*/",
	},
	2: {
		FileExt:        []string{""},
		SingleLine:     "//",
		MultiLineStart: "/*",
		MultiLineEnd:   "*/",
	},
}

//@todo create Comments struct to match comments on different file types
//@todo convert from regex to lexers to parse comments
//@todo convert Config struct to get additional user defined comment matching and extend the comments struct with it if not nil
//@todo when we have multi line comments include each line until we hit the comment end

//@todo update master branch build badges
//@todo add unit testing
//@todo add github workflows for testing binaries on different OS's
var ListedFiles []string
var Cfg Config

func checkExclude(path string, outfile string) (string, bool) {
	//@todo add builds to autorun in VSCode
	//@todo make matching workflows for each build on Github to show status for each arch
	//@todo implement http://tools.ietf.org/html/draft-ietf-websec-mime-sniff to check for files we cant parse like binary data
	m := Cfg.IgnoredItems["ignore"]
	//If we are outputting to a file ignore the output file by default
	if outfile != "" {
		m = append(m, outfile)
	}
	for _, i := range m {
		v, _ := regexp.Compile(i)
		f, err := os.Open(path)
		contentType, err := GetFileContentType(f)
		if err != nil {
			panic(err)
		}
		fmt.Println(contentType)
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

func checkLine(line string, flag string) bool {
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

func checklines(s *scanner.Scanner, ext string, Showlines bool) string {
	//@todo if ext in Comments struct then get comment characters
	//@todo if ext not in Comments struct then get default comment characters
	if Showlines {
		l := "\t" + strconv.Itoa(s.Position.Line) + ")" + s.TokenText() + "\n"
		return l
	} else {
		l := "\t" + s.TokenText() + "\n"
		return l
	}
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

	//If the -f argument contains a specific file then return the directory its in
	dir, _ := path.Split(*folderFlag)
	//Get settings if there is a settings file in the current directory
	settings, err := ioutil.ReadFile(dir + "/.flowcat")
	//If there is a settings file then get the values
	if err == nil {
		err = yaml.Unmarshal(settings, &Cfg)
		//Ignore errors
		err = yaml.Unmarshal(settings, &Cfg.IgnoredItems)
		//Ignore errors
		Showlines, err = strconv.ParseBool(Cfg.Linenums)
		//@todo should this block? else what should the default be
		if err != nil {
			fmt.Println("linenum should be true or false", err)
		}
	}

	if *lineFlag != false {
		Showlines = *lineFlag
	}
	matchexp = Cfg.Match
	if *matchFlag != "" {
		matchexp = *matchFlag //@todo create matchflag from Comments struct
	}
	//Fallback incase settings is missing the match and one is not specified per an argument
	//@todo remove the fallback as we will get everything from the Comments Struct
	if matchexp == "" {
		matchexp = "//@todo"
	}
	//reg, _ := regexp.Compile(matchexp)

	parseFiles := func(path string, info os.FileInfo, _ error) (err error) {

		if *outputFlag != "" {
			F, err = os.OpenFile(*outputFlag, os.O_WRONLY|io.SeekStart|os.O_CREATE, 0755)
			defer F.Close()
			if err != nil && F != nil {
				fmt.Println("ERROR: could not write output to", *outputFlag)
			}
		}
		if info.Mode().IsRegular() {

			file, exc := checkExclude(path, *outputFlag)

			//If the file does not match our exclusion regex then use it.
			if !exc {
				curfile, err := os.Open(file)
				if err != nil {
					fmt.Println("ERROR: could not open file", file, err)
				}
				contents, err := os.ReadFile(path)
				if err != nil {
					fmt.Println("ERROR: could not read file", file, err)
				}
				contentbytes := []byte(contents)
				if utf8.Valid(contentbytes) {
					var s scanner.Scanner
					s.Error = func(*scanner.Scanner, string) {} // ignore errors
					s.Init(curfile)
					s.Mode = scanner.ScanComments

					//@todo make Showlines Exportable so we dont need to pass it thru in below function as it does not change after initially set
					checklines := func(s scanner.Scanner, path string, Showlines bool) string {
						//ext := filepath.Ext(path)
						tok := s.Scan()
						var line string
						// if line {
						// 	fmt.Println(path)
						// 	fmt.Println(line)
						// }
						for tok != scanner.EOF {
							if tok == scanner.Comment {
								if Showlines {
									line += "\t" + strconv.Itoa(s.Position.Line) + ")" + s.TokenText() + "\n"
								} else {
									line += "\t" + s.TokenText() + "\n"
									//@todo only return strings that start with the -m argument

								}
							}
							tok = s.Scan()
						}

						return line
					}
					filelines := checklines(s, path, Showlines)
					if filelines != "" {
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
	//@todo fix tests for new functionality
	err = filepath.Walk(*folderFlag, parseFiles)
	if err != nil {
		fmt.Println(err)
	}
}
