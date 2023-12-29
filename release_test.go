package main

//This file holds test performed on binaries and should only be ran after the main_test.go file

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)

// fmt.Println("Test Init when not yet inited")
// cmd := exec.Command("/usr/bin/bash", "-c", "export PATH=$PATH:/usr/local/go/bin; go run main.go logger.go lexer.go init")
// cmd.Stdout = &stdout
// cmd.Stderr = &stderr
// err = cmd.Run()
// if err != nil {
// 	t.Fatal("Error initing flowcat", err.Error())
// }
// if stdout.String() != "" {
// 	fmt.Println(stdout.String())
// }

func TestBuildBinaries(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// Create the binaries
	fmt.Println("Building binary for darwin-arm64")
	// env GOOS=darwin GOARCH=arm64 GO111MODULE=auto go build *.go -o bin/flowcat-darwin-arm64/flowcat
	cmd := exec.Command("/usr/bin/bash", "-c", "export PATH=$PATH:/usr/local/go/bin; env GOOS=darwin GOARCH=arm64 GO111MODULE=auto go build -o bin/flowcat-darwin-arm64/flowcat")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		t.Fatal("Error compiling binary darwin-arm64", err.Error())
	}
	if stdout.String() != "" {
		fmt.Println(stdout.String())
	}
	fmt.Println("Building binary for linux-386")
	cmd = exec.Command("/usr/bin/bash", "-c", "export PATH=$PATH:/usr/local/go/bin; env GOOS=linux GOARCH=386 GO111MODULE=auto go build -o bin/flowcat-linux-386/flowcat")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		t.Fatal("Error compiling binary linux-386", err.Error())
	}
	if stdout.String() != "" {
		fmt.Println(stdout.String())
	}

	fmt.Println("Building binary for linux-amd64")
	cmd = exec.Command("/usr/bin/bash", "-c", "export PATH=$PATH:/usr/local/go/bin; env GOOS=linux GOARCH=amd64 GO111MODULE=auto go build -o bin/flowcat-linux-amd64/flowcat")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		t.Fatal("Error compiling binary linux-amd64", err.Error())
	}
	if stdout.String() != "" {
		fmt.Println(stdout.String())
	}
}
