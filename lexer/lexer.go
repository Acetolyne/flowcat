package main

import (
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
	"AT", "PLUS", "STAR", "DASH", "SLASH", "BACKSLASH", "CARROT", "BACKTICK", "COMMA", "LPAREN", "RPAREN",
	"BUS", "COMPUTE", "CHIP", "IGNORE", "LABEL", "SET", "NUMBER", "NAME",
	"COMMENT", "SPACE",
}
var tokmap map[string]int
var lexer *lexmachine.Lexer

func init() {
	tokmap = make(map[string]int)
	for id, name := range tokens {
		tokmap[name] = id
	}
}

func newLexer(dfa bool) *lexmachine.Lexer {
	getToken := func(tokenType int) lexmachine.Action {
		return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return s.Token(tokenType, string(m.Bytes), m), nil
		}
	}
	var lexer = lexmachine.NewLexer()
	// lexer.Add([]byte("@"), getToken(tokmap["AT"]))
	// lexer.Add([]byte(`\+`), getToken(tokmap["PLUS"]))
	// lexer.Add([]byte(`\*`), getToken(tokmap["STAR"]))
	// lexer.Add([]byte("-"), getToken(tokmap["DASH"]))
	// lexer.Add([]byte("/"), getToken(tokmap["SLASH"]))
	// lexer.Add([]byte("\\"), getToken(tokmap["BACKSLASH"]))
	// lexer.Add([]byte(`\^`), getToken(tokmap["CARROT"]))
	// lexer.Add([]byte("`"), getToken(tokmap["BACKTICK"]))
	// lexer.Add([]byte(","), getToken(tokmap["COMMA"]))
	// lexer.Add([]byte(`\(`), getToken(tokmap["LPAREN"]))
	// lexer.Add([]byte(`\)`), getToken(tokmap["RPAREN"]))
	// lexer.Add([]byte("bus"), getToken(tokmap["BUS"]))
	// lexer.Add([]byte("chip"), getToken(tokmap["CHIP"]))
	// lexer.Add([]byte("label"), getToken(tokmap["LABEL"]))
	// lexer.Add([]byte("compute"), getToken(tokmap["COMPUTE"]))
	// lexer.Add([]byte("ignore"), getToken(tokmap["IGNORE"]))
	// lexer.Add([]byte("set"), getToken(tokmap["SET"]))
	// lexer.Add([]byte(`[0-9]*\.?[0-9]+`), getToken(tokmap["NUMBER"]))
	// lexer.Add([]byte(`[a-zA-Z_][a-zA-Z0-9_]*`), getToken(tokmap["NAME"]))
	// lexer.Add([]byte(`"[^"]*"`), getToken(tokmap["NAME"]))
	lexer.Add([]byte(`#[^\n]*`), getToken(tokmap["COMMENT"]))
	// lexer.Add([]byte(`\s+`), getToken(tokmap["SPACE"]))
	lexer.Add([]byte(`//[^\n]*\n?`), getToken(tokmap["COMMENT"]))
	lexer.Add([]byte(`/\*([^*]|\r|\n|(\*+([^*/]|\r|\n)))*\*+/`), getToken(tokmap["COMMENT"]))
	// bs, _ := json.Marshal(tokmap)
	// fmt.Println(string(bs))
	//{"AT":0,"BACKSLASH":5,"BACKTICK":7,"BUS":11,"CARROT":6,"CHIP":13,"COMMA":8,"COMMENT":19,"COMPUTE":12,"DASH":3,"IGNORE":14,"LABEL":15,"LPAREN":9,"NAME":18,"NUMBER":17,"PLUS":1,"RPAREN":10,"SET":16,"SLASH":4,"SPACE":20,"STAR":2}
	var err error
	if dfa {
		err = lexer.CompileDFA()
	} else {
		err = lexer.CompileNFA()
	}
	if err != nil {
		panic(err)
	}
	return lexer
}

func scan(text []byte) error {
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
			if curtok.Type == 19 {
				fmt.Println(curtok.Value)
			}
		}
	}
	return nil
}

func GetComments(dfa bool, text []byte) {

	lexer = newLexer(dfa)
	scan(text)

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
