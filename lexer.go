package main

//@todo use the map like if map[FILE TYPE] == tok.Type then add token to a slice of token structs
//@todo return a slice of structs {Line, Token}
//@todo add current todo regex into middle of comment regex when adding to lexer
//@todo pass in the users regex to call to GetComments from main file
//@todo create more types of comments
//@todo make a go test file
//@todo make the test.txt file be one file to test all types of comments
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type Token struct {
	Type        int
	Value       interface{}
	Lexeme      []byte
	TC          int
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

// add function and the type will point to the comment with the appropriate regex
type CommentValues struct {
	Ext  []string //list of all the file extensions associated with the TYPES for lexer tokens
	Type int      //list of the token types used for this typically one regex for single line comments and one for multiline
}

var tokens = []string{
	"IGNORE", "ML-COMMENT-COMMON-B", "ML-SINGLE-STYLE", "SL-COMMENT-COMMON-A", "ML-COMMENT-COMMON-A", "SL-SHELL-STYLE", "SL-HTML-STYLE", "ML-HTML-STYLE", "SL-LUA-STYLE", "ML-LUA-STYLE", "ML-RUBY-STYLE", "TEMPLATE-STYLE",
}
var tokmap map[string]int
var lexer *lexmachine.Lexer

var Extensions = []CommentValues{
	{
		Ext: []string{"py"},
		// Multiline each starting with #
		Type: 1,
	},
	{
		Ext: []string{"rs"},
		// Multiline each starting with //
		Type: 2,
	},
	{
		Ext: []string{"", "go", "py", "js", "rs", "html", "gohtml", "php", "c", "cpp", "h", "class", "jar", "java", "jsp", "php"},
		// startSingle: "//",
		Type: 3,
	},
	{
		Ext: []string{"", "go", "py", "js", "rs", "html", "gohtml", "php", "c", "cpp", "h", "class", "jar", "java", "jsp", "php"},
		// startMulti:  "/*",
		// endMulti:    "*/",
		Type: 4,
	},
	{
		Ext: []string{"sh", "php", "rb", "py"},
		// startSingle: "#",
		Type: 5,
	},
	{
		Ext: []string{"html", "gohtml", "md"},
		// Singleline HTML STYLE:  "<!--" COMMENT "-->",
		Type: 6,
	},
	{
		Ext: []string{"html", "gohtml", "md"},
		// startMulti:  "<!--",
		// endMulti:    "-->",
		Type: 7,
	},
	{
		Ext: []string{"lua"},
		// startSingle: "--",
		Type: 8,
	},
	{
		Ext: []string{"lua"},
		// startMulti:  "--[[",
		// endMulti:    "--]]",
		Type: 9,
	},
	{
		Ext: []string{"rb"},
		// startMulti:  "=begin",
		// endMulti:    "=end",
		Type: 10,
	},
	{
		Ext: []string{"tmpl"},
		// startMulti: "{{/*",
		// endMulti:   "*/}}",
		Type: 11,
	},
	{
		Ext: []string{"rs"},
		// Multiline each starting with //
		Type: 12,
	},
}

func init() {
	tokmap = make(map[string]int)
	for id, name := range tokens {
		tokmap[name] = id
	}
	// logFile, err := os.OpenFile("LOGPATH", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	//     log.Panic(err)
	// }
	// defer logFile.Close()
	// log.SetOutput(logfile)
	// log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func lexReg(a []byte, b string, c []byte) []byte {
	reg := append(a, b...)
	reg = append(reg, c...)
	return reg
}

func newLexer(match string) *lexmachine.Lexer {
	logger := GetLoggerType()
	getToken := func(tokenType int) lexmachine.Action {
		return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return s.Token(tokenType, string(m.Bytes), m), nil
		}
	}
	var lexer = lexmachine.NewLexer()
	//NEGATIVE LOOKAHEADS NOT SUPPORTED use START([^E]|E[^N]|EN[^D])*MATCH([^E]|E[^N]|EN[^D])*END to generate a long POSIX compatible regex
	//@todo create tests for:
	//partial ending token in middle of comment does not cut comment off
	//comment in double quotes
	//comment in single quotes
	//lexer.Add([]byte(`#[^\n]*`), getToken(tokmap["COMMENT"]))
	lexer.Add(lexReg([]byte(`[\"]//[ ]*`), match, []byte(`[^\n]*[\"][^\n]*`)), getToken(tokmap["IGNORE"])) //IGNORE THE MATCH WHEN IT IS BETWEEN DOUBLE QUOTES
	lexer.Add(lexReg([]byte(`[\']//[ ]*`), match, []byte(`[^\n]*[\'][^\n]*`)), getToken(tokmap["IGNORE"])) //IGNORE THE MATCH WHEN IT IS BETWEEN SINGLE QUOTES
	lexer.Add(lexReg([]byte(`(#(@todo)*.*#`), match, []byte(`([^\n]|\n[#])*)`)), getToken(tokmap["ML-COMMENT-COMMON-B"]))
	//lexer.Add(lexReg([]byte(`//.*(\n[\/][\/].*)*`), match, []byte(`.*(\n[\/][\/].*)*`)), getToken(tokmap["ML-SINGLE-STYLE"]))
	lexer.Add(lexReg([]byte(`//[ ]*`), match, []byte(`[^\n]*`)), getToken(tokmap["SL-COMMENT-COMMON-A"]))                                                           //SL-COMMENT-COMMON-A
	lexer.Add(lexReg([]byte(`\/\*([^\*]|\*[^\/])*`), match, []byte(`([^\*]|\*[^\/])*\*\/`)), getToken(tokmap["ML-COMMENT-COMMON-A"]))                               //ML-COMMENT-COMMON-A
	lexer.Add(lexReg([]byte(`#[ ]*`), match, []byte(`[^\n]*`)), getToken(tokmap["SL-SHELL-STYLE"]))                                                                 //SL-SHELL-STYLE
	lexer.Add(lexReg([]byte(`<!--[ ]*`), match, []byte(`[^\n]*-->`)), getToken(tokmap["SL-HTML-STYLE"]))                                                            //SL-HTML-STYLE
	lexer.Add(lexReg([]byte(`<!--([^-]|-[^-]|--[^>])*`), match, []byte(`([^-]|-[^-]|--[^>])*-->`)), getToken(tokmap["ML-HTML-STYLE"]))                              //ML-HTML-STYLE
	lexer.Add(lexReg([]byte(`\-\-[ ]*`), match, []byte(`[^\n]*`)), getToken(tokmap["SL-LUA-STYLE"]))                                                                //SL-LUA-STYLE
	lexer.Add(lexReg([]byte(`--\[\[([^-]|-[^-]|--[^\]|\-\-\][^\]])*`), match, []byte(`([^-]|-[^-]|--[^\]]|\-\-\][^\]])*--\]\]`)), getToken(tokmap["ML-LUA-STYLE"])) //ML-LUA-STYLE
	lexer.Add(lexReg([]byte(`\=begin([^=]|=[^e]|=e[^n]|=en[^d])*`), match, []byte(`([^=]|=[^e]|=e[^n]|=en[^d])*`)), getToken(tokmap["ML-RUBY-STYLE"]))
	//@todo regex below has a bug so does not return the last token when doing the match on all chars except */}} when we add the last } it makes the asterisk not match in ending of token.
	lexer.Add(lexReg([]byte(`\{\{\/\*([^\*]|\*[^\/]|\*\/[^\}]|\*\/\}[^\}])*`), match, []byte(`([^\*]|\*[^\/]|\*\/[^\}]|\*\/}[^}])*`)), getToken(tokmap["TEMPLATE-STYLE"]))
	// lexer.Add(lexReg([]byte(`<!--[ ]*`), match, []byte(`[^\n]*-->`)), getToken(tokmap["WHAT-STYLE"]))
	// lexer.Add(lexReg([]byte(`<!--[ ]*`), match, []byte(`[^\n]*-->`)), getToken(tokmap["WHAT-STYLE"]))
	// lexer.Add(lexReg([]byte(`<!--[ ]*`), match, []byte(`[^\n]*-->`)), getToken(tokmap["WHAT-STYLE"]))
	// lexer.Add(lexReg([]byte(`<!--[ ]*`), match, []byte(`[^\n]*-->`)), getToken(tokmap["WHAT-STYLE"]))
	// lexer.Add(lexReg([]byte(`<!--[ ]*`), match, []byte(`[^\n]*-->`)), getToken(tokmap["WHAT-STYLE"]))

	//Gets all the token types and their corresponding ids
	bs, _ := json.Marshal(tokmap)
	logger.Info.Println("LEXER REGEX MAP\n", string(bs))
	//@todo create conf variable to assign NFA or DFA and compile accordingly below
	// DFA but DFA uses less memory
	// real    1m13.999s
	// user    1m21.119s
	// sys     0m3.862s
	// NFA better times but more memory intensive
	// real    0m6.187s
	// user    0m5.289s
	// sys     0m1.491s
	err := lexer.CompileNFA()
	if err != nil {
		logger.Fatal.Println("Could not compile the lexer", err.Error())
		panic(err)
	}
	return lexer
}

func Scan(text []byte, path string, showlines bool, outputFile string) error {
	logger.Info.Println("Scanning file:", path)
	var f *os.File
	var err error
	//open output file in preparation
	if outputFile != "" {
		dir, _ := filepath.Split(outputFile)
		if dir == "" {
			folder, _ := filepath.Split(path)
			outputFile = folder + outputFile
		}
		logger.Info.Println("Output will be sent to", outputFile)
		f, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
		if err != nil {
			logger.Err.Println("Could not open output file", err.Error())
		}
		//defer f.Close()
	}
	//ensures we print the filename only once per file
	printfile := true
	_, curfile := filepath.Split(path)
	ext := strings.Split(curfile, ".")
	if len(ext) < 2 {
		ext = append(ext, "")
	}
	scanner, err := lexer.Scanner(text)
	if err != nil {
		logger.Err.Println("Error while scanning text", string(text), err.Error())
		return err
	}
	for tk, err, eof := scanner.Next(); !eof; tk, err, eof = scanner.Next() {
		if ui, is := err.(*machines.UnconsumedInput); ui != nil && is {
			scanner.TC = ui.FailTC
			//fmt.Println("skipping", ui)
		} else if err != nil {
			logger.Err.Println("Error while lexing for comments", err.Error())
			return err
		} else {
			curtok := tk.(*lexmachine.Token)
			for _, CommentValue := range Extensions {
				for _, curext := range CommentValue.Ext {
					if curext == ext[1] {
						//fmt.Println(CommentValue.Type, curtok.Type)
						if CommentValue.Type == curtok.Type {
							if printfile {
								fmt.Println(path)
								if outputFile != "" {
									_, err = f.WriteString(path + "\n")
									if err != nil {
										logger.Err.Println("Could not write output file", err.Error())
									}
								}
								printfile = false
							}
							if showlines {
								fmt.Println(" ", curtok.StartLine, ")", curtok.Value)
								if outputFile != "" {
									str := fmt.Sprintf(" %s)%s\n", strconv.Itoa(curtok.StartLine), curtok.Value.(string))
									_, err = f.WriteString(str)
									if err != nil {
										logger.Err.Println("Could not write to output file", err.Error())
									}
								}
							} else {
								fmt.Println(" ", curtok.Value)
								if outputFile != "" {
									str := fmt.Sprintf("  %s\n", curtok.Value.(string))
									_, err = f.WriteString(str)
									if err != nil {
										logger.Err.Println("Could not write output file", err.Error())
									}
								}
							}
						}
					}
					// }
				}
			}
		}
	}
	//f.Close()
	return nil
}

func GetComments(text []byte, match string, path string, showlines bool, outputFile string) {
	logger.Info.Println("matching on", match)
	lexer = newLexer(match)
	err := Scan(text, path, showlines, outputFile)
	if err != nil {
		logger.Err.Println("Error scanning text", err.Error())
	}
}
