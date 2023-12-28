package main_test

//This file holds test performed on binaries

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

////Helper functions

// //Testing starts
// Performs pre-test actions like creating the binaries and a tmp folder for test files
func TestPre(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

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
	//Create the binaries
	fmt.Println("Building binary for darwin-arm64")
	//env GOOS=darwin GOARCH=arm64 GO111MODULE=auto go build *.go -o bin/flowcat-darwin-arm64/flowcat
	cmd := exec.Command("/usr/bin/bash", "-c", "export PATH=$PATH:/usr/local/go/bin; env GOOS=darwin GOARCH=arm64 GO111MODULE=auto go build -o bin/flowcat-darwin-arm64/flowcat")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
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

	//Create tmp file folder
	err = os.MkdirAll("tests/files", 0775)
	if err != nil {
		t.Fatal("Could not perform pre-test actions, creating the tmp file folder failed", err.Error())
	}
}

//@todo init settings works
//@todo init settings fails if ran a second time and the file is already there
//@todo can get config settings from file
