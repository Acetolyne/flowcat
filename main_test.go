package main

//This file holds test performed on binaries

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// Create the tests
// basic tests
var tests = []struct {
	path       string
	outfile    string
	folderFlag string
	exp        string
}{
	{"test.go", "", "", "test.go false"},
	{"todo", "", "./", "todo true"},
}

var outputfiletests = []struct {
	text       []byte
	path       string
	showlines  bool
	outputFile string
	exp        string
}{
	{[]byte("//@todo some comment"), "tests/files/test.go", true, "output.txt", "tests/files/test.go\n 1)//@todo some comment\n"},
	{[]byte("//@todo some comment"), "tests/files/test.go", true, "./output.txt", "tests/files/test.go\n 1)//@todo some comment\n"},
}

// //Testing starts
// Performs pre-test actions like creating the binaries and a tmp folder for test files
func TestPre(t *testing.T) {
	//Create binaries folders
	err := os.MkdirAll("bin/flowcat-darwin-arm64", 0775)
	if err != nil {
		t.Fatal("[PRE-TEST][BIN FOLDER] creating the folder bin/flowcat-darwin-arm64 failed", err.Error())
	}
	err = os.MkdirAll("bin/flowcat-linux-386", 0775)
	if err != nil {
		t.Fatal("[PRE-TEST][BIN FOLDER] creating the folder bin/flowcat-linux-386 failed", err.Error())
	}
	err = os.MkdirAll("bin/flowcat-linux-amd64", 0775)
	if err != nil {
		t.Fatal("[PRE-TEST][BIN FOLDER] creating the folder bin/flowcat-linux-amd64 failed", err.Error())
	}

	//Create tmp file folder
	err = os.MkdirAll("tests/files", 0775)
	if err != nil {
		t.Fatal("Could not perform pre-test actions, creating the tmp file folder failed", err.Error())
	}
	// TmpDir := t.TempDir()

	// //Create a tmp file
	// TmpFile := t.TempFile("tests", "files", "test.go")
}

// func TestCfg(t *testing.T) {
// 	var c Config
// 	c.Linenum := false
// 	c.Match := ""
// }

func TestCheckExclude(t *testing.T) {
	//Run the tests
	Cfg.IgnoredItems = append(Cfg.IgnoredItems, "todo")
	for _, e := range tests {
		file, exc := CheckExclude(e.path, e.outfile, e.folderFlag)
		res := file + " " + strconv.FormatBool(exc)
		if res != e.exp {
			t.Errorf("Got: %s Expected: %s", res, e.exp)
		}
	}
}

func TestOutputFile(t *testing.T) {
	for _, e := range outputfiletests {
		lexer = newLexer("@todo") //sets the matching string
		err := Scan(e.text, e.path, e.showlines, e.outputFile)
		if err != nil {
			t.Errorf("Scan failed %s", err.Error())
		}
		dir, _ := filepath.Split(e.outputFile)
		if dir == "" {
			folder, _ := filepath.Split(e.path)
			b, err := os.ReadFile(folder + e.outputFile) // pass path at -f plus output filename
			if err != nil {
				t.Errorf("Could not read output file %s %s", e.outputFile, err.Error())
			}
			if e.exp != string(b) {
				t.Errorf("Got: %s Expected: %s", string(b), e.exp)
			}
			err = os.Remove(folder + e.outputFile)
			if err != nil {
				t.Errorf("Could not remove old file %s %s", folder+e.outputFile, err.Error())
			}
		} else {
			b, err := os.ReadFile(e.outputFile) // just pass the output file name
			if err != nil {
				t.Errorf("Could not read output file %s", e.outputFile)
			}
			if e.exp != string(b) {
				t.Errorf("Got: %s Expected: %s", string(b), e.exp)
			}
			err = os.Remove(e.outputFile)
			if err != nil {
				t.Errorf("Could not remove old file %s %s", e.outputFile, err.Error())
			}
		}

	}
}

//@todo init settings works
//@todo init settings fails if ran a second time and the file is already there
//@todo can get config settings from file
