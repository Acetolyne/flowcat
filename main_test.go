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
	fmt.Printf("Building binary for darwin-arm64")
	//env GOOS=darwin GOARCH=arm64 GO111MODULE=auto go build *.go -o bin/flowcat-darwin-arm64/flowcat
	cmd := exec.Command("go build *.go -o bin/flowcat-darwin-arm64/flowcat")
	//cmd.Dir = os.Getenv("TMP_PATH") + u
	env := []string{"GOOS=darwin", "GOARCH=arm64", "GO111MODULE=auto", "PATH=$PATH:/usr/local/go/bin"}
	cmd.Env = env
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		t.Fatal("Error comiling binary darwin-arm64", err.Error())
	}
	if stdout.String() != "" {
		fmt.Println(stdout.String())
	}
	fmt.Println("Building binary for linux-386")
	// env GOOS=linux GOARCH=386 GO111MODULE=auto go build -o bin/flowcat-linux-386/flowcat

	fmt.Println("Building binary for linux-amd64")
	// env GOOS=linux GOARCH=amd64 GO111MODULE=auto go build -o bin/flowcat-linux-amd64/flowcat
	//Create tmp file folder
	err = os.MkdirAll("tests/files", 0775)
	if err != nil {
		t.Fatal("Could not perform pre-test actions, creating the tmp file folder failed", err.Error())
	}
}

//@todo init settings works
//@todo init settings fails if ran a second time and the file is already there
//@todo can get config settings from file
