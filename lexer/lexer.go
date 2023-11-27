package lexer

//@todo create map of file extensions or file headers and comment types
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

var tokens = []string{
	"IGNORE", "COMMENT",
}
var tokmap map[string]int
var lexer *lexmachine.Lexer

func init() {
	tokmap = make(map[string]int)
	for id, name := range tokens {
		tokmap[name] = id
	}
}

func newLexer() *lexmachine.Lexer {
	getToken := func(tokenType int) lexmachine.Action {
		return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return s.Token(tokenType, string(m.Bytes), m), nil
		}
	}
	var lexer = lexmachine.NewLexer()
	//lexer.Add([]byte(`#[^\n]*`), getToken(tokmap["COMMENT"]))
	lexer.Add([]byte(`[\"]//[ ]*@todo[^\n]*[\"][^\n]*`), getToken(tokmap["IGNORE"]))
	lexer.Add([]byte(`//[ ]*@todo[^\n]*`), getToken(tokmap["COMMENT"]))
	//lexer.Add([]byte(`/\*([^*]|\r|\n|(\*+([^*/]|\r|\n)))*\*+/`), getToken(tokmap["COMMENT"])) //Multi line comments
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

func scan(text []byte) error {
	var AllTokens []Token
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
			if curtok.Type == 1 {
				fmt.Println(curtok.Value)
				AllTokens = append(AllTokens, curtok)
			}
		}
	}
	return nil
}

func GetComments(text []byte) {

	lexer = newLexer()
	err := scan(text)
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
