package lexer

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
	"path/filepath"
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
	Ext  []string //list of all the file extentions associated with the TYPES for lexer tokens
	Type int      //list of the token types used for this typically one regex for single line comments and one for multiline
}

var tokens = []string{
	"IGNORE", "SL-COMMENT-COMMON-A", "ML-COMMENT-COMMON-A", "SL-SHELL-STYLE", "SL-HTML-STYLE", "ML-HTML-STYLE", "SL-LUA-STYLE", "ML-LUA-STYLE", "ML-RUBY-STYLE",
}
var tokmap map[string]int
var lexer *lexmachine.Lexer

var Extensions = []CommentValues{
	{
		Ext: []string{"", "go", "py", "js", "rs", "html", "gohtml", "php", "c", "cpp", "h", "class", "jar", "java", "jsp", "php"},
		// startSingle: "//",
		Type: 1,
	},
	{
		Ext: []string{"", "go", "py", "js", "rs", "html", "gohtml", "php", "c", "cpp", "h", "class", "jar", "java", "jsp", "php"},
		// startMulti:  "/*",
		// endMulti:    "*/",
		Type: 2,
	},
	{
		Ext: []string{"sh", "php", "rb"},
		// startSingle: "#",
		Type: 3,
	},
	{
		Ext: []string{"html", "gohtml", "md"},
		// Singleline HTML STYLE:  "<!--" COMMENT "-->",
		Type: 4,
	},
	{
		Ext: []string{"html", "gohtml", "md"},
		// startMulti:  "<!--",
		// endMulti:    "-->",
		Type: 5,
	},
	{
		Ext: []string{"lua"},
		// startSingle: "--",
		Type: 6,
	},
	{
		Ext: []string{"lua"},
		// startMulti:  "--[[",
		// endMulti:    "--]]",
		Type: 7,
	},
	{
		Ext: []string{"rb"},
		// startMulti:  "=begin",
		// endMulti:    "=end",
		Type: 8,
	},
	{
		Ext: []string{".py"},
		// startSingle: "#",
		Type: 0,
	},
	{
		Ext: []string{".tmpl"},
		// startMulti: "{{/*",
		// endMulti:   "*/}}",
		Type: 0,
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
	lexer.Add(lexReg([]byte(`[\"]//[ ]*`), match, []byte(`[^\n]*[\"][^\n]*`)), getToken(tokmap["IGNORE"]))                                                          //IGNORE THE MATCH WHEN IT IS BETWEEN DOUBLE QUOTES
	lexer.Add(lexReg([]byte(`[\']//[ ]*`), match, []byte(`[^\n]*[\'][^\n]*`)), getToken(tokmap["IGNORE"]))                                                          //IGNORE THE MATCH WHEN IT IS BETWEEN SINGLE QUOTES
	lexer.Add(lexReg([]byte(`//[ ]*`), match, []byte(`[^\n]*`)), getToken(tokmap["SL-COMMENT-COMMON-A"]))                                                           //SL-COMMENT-COMMON-A
	lexer.Add(lexReg([]byte(`\/\*([^\*]|\*[^\/])*`), match, []byte(`([^\*]|\*[^\/])*\*\/`)), getToken(tokmap["ML-COMMENT-COMMON-A"]))                               //ML-COMMENT-COMMON-A
	lexer.Add(lexReg([]byte(`#[ ]*`), match, []byte(`[^\n]*`)), getToken(tokmap["SL-SHELL-STYLE"]))                                                                 //SL-SHELL-STYLE
	lexer.Add(lexReg([]byte(`<!--[ ]*`), match, []byte(`[^\n]*-->`)), getToken(tokmap["SL-HTML-STYLE"]))                                                            //SL-HTML-STYLE
	lexer.Add(lexReg([]byte(`<!--([^-]|-[^-]|--[^>])*`), match, []byte(`([^-]|-[^-]|--[^>])*-->`)), getToken(tokmap["ML-HTML-STYLE"]))                              //ML-HTML-STYLE
	lexer.Add(lexReg([]byte(`\-\-[ ]*`), match, []byte(`[^\n]*`)), getToken(tokmap["SL-LUA-STYLE"]))                                                                //SL-LUA-STYLE
	lexer.Add(lexReg([]byte(`--\[\[([^-]|-[^-]|--[^\]|\-\-\][^\]])*`), match, []byte(`([^-]|-[^-]|--[^\]]|\-\-\][^\]])*--\]\]`)), getToken(tokmap["ML-LUA-STYLE"])) //ML-LUA-STYLE
	lexer.Add(lexReg([]byte(`\=begin([^=]|=[^e]|=e[^n]|=en[^d])*`), match, []byte(`([^=]|=[^e]|=e[^n]|=en[^d])*`)), getToken(tokmap["ML-RUBY-STYLE"]))
	//Gets all the token types and their cooresponding ids
	bs, _ := json.Marshal(tokmap)
	fmt.Println(string(bs))
	//{"AT":0,"BACKSLASH":5,"BACKTICK":7,"BUS":11,"CARROT":6,"CHIP":13,"COMMA":8,"COMMENT":19,"COMPUTE":12,"DASH":3,"IGNORE":14,"LABEL":15,"LPAREN":9,"NAME":18,"NUMBER":17,"PLUS":1,"RPAREN":10,"SET":16,"SLASH":4,"SPACE":20,"STAR":2}

	err := lexer.CompileDFA()
	if err != nil {
		panic(err)
	}
	return lexer
}

func scan(text []byte, path string, showlines bool) error {
	printfile := true
	_, curfile := filepath.Split(path)
	ext := strings.Split(curfile, ".")
	//var CommentValue *CommentValues
	scanner, err := lexer.Scanner(text)
	if err != nil {
		return err
	}
	for tk, err, eof := scanner.Next(); !eof; tk, err, eof = scanner.Next() {
		if ui, is := err.(*machines.UnconsumedInput); ui != nil && is {
			scanner.TC = ui.FailTC
			//log.Printf("skipping %v", ui)
		} else if err != nil {
			return err
		} else {
			curtok := tk.(*lexmachine.Token)
			if ext[1] == "" {
				// log.Println("Logging to custom file")
				fmt.Println("No extension for file not parsing")
			} else {
				//fmt.Println("EXT:", ext[1])
				for _, CommentValue := range Extensions {
					for _, curext := range CommentValue.Ext {
						if curext == ext[1] {
							//fmt.Println(CommentValue.Type, curtok.Type)
							if CommentValue.Type == curtok.Type {
								if printfile {
									fmt.Println(path)
									printfile = false
								}
								if showlines {
									fmt.Println(" ", curtok.StartLine, ")", curtok.Value)
								} else {
									fmt.Println(" ", curtok.Value)
								}
							}
						}
					}
				}
			}
			// if curtok.Type == 1 {
			// 	fmt.Println(curtok.Value)
			// 	AllTokens = append(AllTokens, curtok)
			// }
		}
	}
	return nil
}

func GetComments(text []byte, match string, path string, showlines bool) {
	fmt.Println("matching on", match)
	lexer = newLexer(match)
	err := scan(text, path, showlines)
	if err != nil {
		fmt.Println("Error scanning text", err)
	}

	// scanner, err := lexer.Scanner(text)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for tok, err, eos := scanner.Next(); !eos; tok, err, eos = scanner.Next() {
	// 	if err != nil {
	// 		// handle the error and exit the loop. For example:
	// 		fmt.Println(err)
	// 	}
	// 	// do some processing on tok or store it somewhere. eg.
	// 	curtok := tok.(*lexmachine.Token)
	// 	fmt.Println(string(curtok.Lexeme))
	// }
	// for i := 0; i < 1000; i++ {
	// 	err = scan(text)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
}
