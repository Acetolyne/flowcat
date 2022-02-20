package main_test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

// TestInitsettings tests that we can properly use the settings from the .flowcat settings file
func TestCanChooseFile(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-f", "tests/assets/test.go")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	//fmt.Println(string(stdout))
	want := "tests/assets/test.go //@todo Regular comment /* Some multiline  @todo comment  comment2 */ //@todo comment 3"
	cur := string(stdout)
	ret := strings.Replace(cur, "\n", "", -1)
	ret = strings.Replace(ret, "\t", " ", -1)
	if want != ret {
		t.Errorf("got %s, want %s", ret, want)
	}
}

func TestCanUseLineOption(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-f", "tests/assets/test.go", "-l")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	//fmt.Println(string(stdout))
	want := "tests/assets/test.go 4)//@todo Regular comment 6)/* Some multiline  @todo comment  comment2 */ 9)//@todo comment 3"
	cur := string(stdout)
	ret := strings.Replace(cur, "\n", "", -1)
	ret = strings.Replace(ret, "\t", " ", -1)
	if want != ret {
		t.Errorf("got %s, want %s", ret, want)
	}

}

// func TestCanUseMatching(t *testing.T) {

// }

// func TestCanParseMultipleFiles(t *testing.T) {

// }

//tests that we can send the output to a file and that file is not included in the scan
// func TestCanOutputToFile(t *testing.T) {

// }
