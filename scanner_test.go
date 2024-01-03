package main

import (
	"os"
	"testing"
)

// (text []byte, path string, showlines bool, outputFile string)
var scantests = []struct {
	text       []byte
	path       string
	showlines  bool
	outputFile string
	match      string
	exp        string
}{
	//Golang tests
	{[]byte("//@todo some comment"), "test.go", true, "output.txt", "@todo", "test.go\n 1)//@todo some comment\n"},                                                       //can get single line comment
	{[]byte("some words //@todo some comment"), "test.go", true, "output.txt", "@todo", "test.go\n 1)//@todo some comment\n"},                                            //can get comment at end of line
	{[]byte("//@todo some comment"), "test.go", true, "output.txt", "todo", ""},                                                                                          //ignores comments when a different match is used
	{[]byte("fmt.Println(\"//@todo some comment\")"), "test.go", true, "output.txt", "@todo", ""},                                                                        //ignores comments inside a print statement or string
	{[]byte("/* Some multiline \n@todo comment \ncomment2 */"), "test.go", true, "output.txt", "@todo", "test.go\n 1)/* Some multiline \n@todo comment \ncomment2 */\n"}, //can get multi-line comments
	//Python tests
	{[]byte("#@todo some comment"), "test.py", true, "output.txt", "@todo", "test.py\n 1)#@todo some comment\n"},
	{[]byte("some words first #@todo some comment"), "test1.py", true, "output.txt", "@todo", "test1.py\n 1)#@todo some comment\n"},
	{[]byte("#@todo some comment"), "test2.py", true, "output.txt", "todo", ""},
	{[]byte("print(\"#@todo some comment\")"), "test3.py", true, "output.txt", "@todo", ""},
	{[]byte("#some multiline \n#@todo comment\n#with 3 lines"), "test4.py", true, "output.txt", "@todo", "test4.py\n 1)#some multiline \n#@todo comment\n#with 3 lines\n"},
	//HTML tests
	{[]byte("<!-- @todo some comment -->"), "test.html", true, "output.txt", "@todo", "test.html\n 1)<!-- @todo some comment -->\n"},                                                       //can get single line comment
	{[]byte("<!-- stuff @todo some comment -->"), "test2.html", true, "output.txt", "@todo", "test2.html\n 1)<!-- stuff @todo some comment -->\n"},                                         //can get comment at end of line
	{[]byte("<!-- @todo some comment -->"), "test3.html", true, "output.txt", "note", ""},                                                                                                  //ignores comments when a different match is used
	{[]byte("fmt.Println(\"<!-- @todo some comment\")"), "test4.html", true, "output.txt", "@todo", ""},                                                                                    //ignores comments inside a print statement or string
	{[]byte("<!--  Some multiline \n@todo comment \ncomment2 \n-->"), "test5.html", true, "output.txt", "@todo", "test5.html\n 1)<!--  Some multiline \n@todo comment \ncomment2 \n-->\n"}, //can get multi-line comments
}

func TestScan(t *testing.T) {
	var Cfg Config
	Cfg.Linenum = "false"
	for _, e := range scantests {
		//fmt.Println("TESTING", e)
		lexer = newLexer(e.match)                              //sets the matching string
		err := Scan(e.text, e.path, e.showlines, e.outputFile) //e.text, e.path, e.showlines, e.outputFile
		if err != nil {
			t.Errorf("Scan failed %s", err.Error())
		}
		b, err := os.ReadFile(e.outputFile) // just pass the output file name
		if err != nil {
			t.Errorf("Could not read output file %s", e.outputFile)
		}
		if e.exp != string(b) {
			t.Errorf("Got: %s Expected: %s", string(b), e.exp)
		}
		err = os.Truncate(e.outputFile, 0)
		if err != nil {
			t.Errorf("%s", err.Error())

		}
	}
	err := os.Remove("output.txt")
	if err != nil {
		t.Errorf("Could not remove old test file output.txt %s", err.Error())
	}
}
