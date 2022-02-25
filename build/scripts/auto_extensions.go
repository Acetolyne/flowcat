package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	lexer "github.com/Acetolyne/commentlex"
)

func main() {
	var s lexer.Scanner
	var buffer string
	s.Mode = lexer.ScanComments
	allext := s.GetExtensions()

	file, err := os.Open("../../README.md")
	if err != nil {
		fmt.Println(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		buffer += scanner.Text() + "\n"
		if strings.Contains(buffer, "##### Supported Filetypes") {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	for l := range allext {
		curext := allext[l]
		buffer += curext + "\n"
	}
	file.Close()
	file, err = os.OpenFile("../../README.md", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	err = file.Truncate(0)
	if err != nil {
		fmt.Println(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
	}
	_, err = fmt.Fprintf(file, "%s", buffer)
	if err != nil {
		fmt.Println(err)
	}
}
