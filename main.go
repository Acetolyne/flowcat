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

	"gopkg.in/yaml.v2"
)

// @todo add logger
// @todo add debug option in config file
// Config structure of configuration from a yaml settings file.
type Config struct {
	Linenum      string   `yaml:"linenum"`
	Match        string   `yaml:"match"`
	IgnoredItems []string `yaml:"ignore"`
}

// ListedFiles returns a string of all files in a directory.
var ListedFiles []string

// Cfg returns the user configurations from a file.
var Cfg Config
var Debug int

// @todo excluded regex not working when -f is specified because -f value is part of the path
func checkExclude(path string, outfile string, folderFlag string) (string, bool) {
	regpath := strings.TrimPrefix(path, folderFlag)
	m := Cfg.IgnoredItems
	reg := []bool{}
	//If we are outputting to a file ignore the output file by default if it is in the project path
	//@todo fix this logic see if path ends with the outfile
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
	logger := GetLoggerType()
	dirname, err := os.UserHomeDir()
	if err != nil {
		logger.Err.Println("Could not get the users home directory", err.Error())
	}
	f, err := os.Open(dirname + "/.flowcat/config")
	f.Close()
	if err != nil {
		var SetFile *os.File
		SetFile, err = os.OpenFile(dirname+"/.flowcatconfig", os.O_WRONLY|io.SeekStart|os.O_CREATE, 0755)
		if err != nil {
			logger.Err.Println("Could not create settings file when running init", err.Error())
			return errors.New("ERROR: could not create settings file")
		}
		defer SetFile.Close()
		_, err = SetFile.WriteString("# Settings\n")
		if err != nil {
			logger.Err.Println("Could not write to settings file during init", err.Error())
			return errors.New("ERROR: could not create settings file")
		}
		_, err = SetFile.WriteString("match: \"@todo\"\n\n")
		if err != nil {
			logger.Err.Println("Could not write to settings file during init", err.Error())
			return errors.New("ERROR: could not create settings file")
		}
		_, err = SetFile.WriteString("# File patterns to ignore\n")
		if err != nil {
			logger.Err.Println("Could not write to settings file during init", err.Error())
			return errors.New("ERROR: could not create settings file")
		}
		_, err = SetFile.WriteString("ignore:\n")
		if err != nil {
			logger.Err.Println("Could not write to settings file during init", err.Error())
			return errors.New("ERROR: could not create settings file")
		}
		_, err = SetFile.WriteString("  - \"^\\\\..*\"\n")
		if err != nil {
			logger.Err.Println("Could not write to settings file during init", err.Error())
			return errors.New("ERROR: could not create settings file")
		}
		SetFile.Close()
		logger.Info.Println("Settings file created at ~/.flowcat/config")
		return nil
	}
	logger.Err.Println("User ran init but the config file already exists")
	fmt.Println("setting file already exists consider editing the file ~/.flowcat/config or delete it before running init if you want to refresh it")
	return errors.New("setting file already exists")
}

func init() {
	logger := GetLoggerType()
	homedir, err := os.UserHomeDir()
	if err != nil {
		logger.Err.Println("Could not get the users home directory", err.Error())
	}
	//Make sure the user directory has a folder called .flowcat and there is a logs folder in it
	err = os.MkdirAll(homedir+"/.flowcat/logs", 0775)
	if err != nil {
		logger.Err.Println("Could not create flowcat directories in user folder", homedir+"/.flowcat/logs", err.Error())
	}
}

func main() {
	logger := GetLoggerType()
	var F *os.File
	var Showlines bool = false
	var matchexp string

	folderFlag := flag.String("f", "./", "The project top level directory, where flowcat should start recursing from.")
	outputFlag := flag.String("o", "", "Optional output file to dump results to, note output will still be shown on terminal.")
	matchFlag := flag.String("m", "", "The string to match to do items on.")
	lineFlag := flag.Bool("l", false, "If line numbers should be shown with todo items in output.") //@todo change this to string so we can override the default in the configuration file if needed
	helpFlag := flag.Bool("h", false, "Shows the help menu.")
	flag.Parse()

	//Helpflag implemented because the default help flag from the flag package returns status code 2
	if *helpFlag {
		fmt.Println("Flowcat version 3.1.1")
		fmt.Println("")
		fmt.Println("Options for Flowcat:")
		fmt.Println("init")
		fmt.Println("using flowcat init creates a settings file for the current user, settings can be changed later in the ~/.flowcat/config file")
		fmt.Println("-f string")
		fmt.Println("   The project top level directory, where flowcat should start recursing from or a specific file (default Current Directory)")
		fmt.Println("-l")
		fmt.Println("	Display line numbers in the output.")
		fmt.Println("-m string")
		fmt.Println("   The string to match to do items on. (default 'TODO')")
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
			logger.Warn.Println("Could not init settings when running flowcat with init argument", err.Error())
		}
		//always exit without running if we were using init argument
		os.Exit(0)
	}

	//Get settings from .flowcat file in users home directory
	dirname, err := os.UserHomeDir()
	if err != nil {
		logger.Warn.Println("Could not find user home directory", dirname, err.Error())
	}
	//@todo check if the folder .flowcat exists in the user dir
	//@todo if it does not exist create it
	//@todo if there is no .flowcat/settings file then create it with defaults
	//@todo logs should go in the .flowcat folder as debug.log and error.log each should be overwritten every time.
	settings, err := os.OpenFile(dirname+"/.flowcat/config", os.O_RDONLY, 0600)
	if err != nil {
		logger.Warn.Println("Could not open user configuration file", dirname+"/.flowcat/config", err.Error())
	}
	defer settings.Close()
	configuration := yaml.NewDecoder(settings)
	err = configuration.Decode(&Cfg)
	if err != nil {
		logger.Warn.Println("Unable to get settings from configuration file.", err.Error())
	}
	// else {
	// 	_ = yaml.Unmarshal(settings, &Cfg)
	// 	// 	//Ignore errors
	// 	_ = yaml.Unmarshal(settings, &Cfg.IgnoredItems)
	// 	// 	//Ignore errors
	// }

	if Cfg.Linenum == "true" {
		Showlines = true
	}
	if *lineFlag {
		Showlines = *lineFlag
	}

	if *matchFlag != "" {
		matchexp = *matchFlag
	} else if Cfg.Match != "" {
		matchexp = Cfg.Match
	} else {
		matchexp = "TODO"
	}
	parseFiles := func(path string, info os.FileInfo, _ error) (err error) {

		if *outputFlag != "" {
			F, err = os.OpenFile(*outputFlag, os.O_WRONLY|io.SeekStart|os.O_CREATE, 0755)
			if err != nil {
				logger.Err.Println("could not open specified output file", *outputFlag, err.Error())
				fmt.Println("WARNING could not create output file", err.Error())
			}
			defer F.Close()
			if err != nil && F != nil {
				logger.Err.Println("could not write to the specified output file", *outputFlag, err.Error())
			}
		}
		if info.Mode().IsRegular() {
			file, exc := checkExclude(path, *outputFlag, *folderFlag)

			//If the file does not match our exclusion regex then use it.
			if !exc {
				logger.Info.Println("Checking file", path)
				contents, err := os.ReadFile(path)
				if err != nil {
					logger.Err.Println("could not read file", file, err)
				}
				contentbytes := []byte(contents)
				if utf8.Valid(contentbytes) {
					GetComments(contentbytes, matchexp, path, Showlines)
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
		logger.Err.Println("An error occured while walking directory", err.Error())
	}
	os.Exit(0)
}
