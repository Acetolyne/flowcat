package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestInitsettings tests that we can properly use the settings from the .flowcat settings file
func TestCanChooseFile(t *testing.T) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	curcmd := exec.Command("rm", "-f", dirname+"/.flowcat")
	_, err = curcmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	curcmd = exec.Command("go", "run", "main.go", "init")
	_, err = curcmd.Output()
	if err != nil {
		fmt.Println(err)
	}
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

func TestCanUseMatching(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-f", "tests/assets/multitest/test.php", "-m", "Comment")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	//fmt.Println(string(stdout))
	want := "tests/assets/multitest/test.php // Comment no todo"
	cur := string(stdout)
	ret := strings.Replace(cur, "\n", "", -1)
	ret = strings.Replace(ret, "\t", " ", -1)
	if want != ret {
		t.Errorf("got %s, want %s", ret, want)
	}

}

func TestCanParseMultipleFiles(t *testing.T) {

	cmd := exec.Command("go", "run", "main.go", "-f", "tests/assets/multitest/")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	//fmt.Println(string(stdout))
	want := "tests/assets/multitest/test.go //@todo Regular comment /* Some multiline  @todo comment  comment2 */ //@todo comment 3tests/assets/multitest/test.php # @todo comment 1 //@todo comment 2 /*@todo multiline comment value */ //@todo another comment # @todo comment 5 //@todo comment 4"
	cur := string(stdout)
	ret := strings.Replace(cur, "\n", "", -1)
	ret = strings.Replace(ret, "\t", " ", -1)
	if want != ret {
		t.Errorf("got %s, want %s", ret, want)
	}

}

//tests that we can send the output to a file and that file is not included in the scan
func TestCanOutputToFile(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-f", "tests/assets/multitest/", "-o", "tests/assets/multitest/output.txt")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	//fmt.Println(string(stdout))
	want := "tests/assets/multitest/test.go //@todo Regular comment /* Some multiline  @todo comment  comment2 */ //@todo comment 3tests/assets/multitest/test.php # @todo comment 1 //@todo comment 2 /*@todo multiline comment value */ //@todo another comment # @todo comment 5 //@todo comment 4"
	cur := string(stdout)
	ret := strings.Replace(cur, "\n", "", -1)
	ret = strings.Replace(ret, "\t", " ", -1)
	if want != ret {
		t.Errorf("got %s, want %s", ret, want)
	}
	_ = exec.Command("rm", "-f", "tests/assets/multitest/output.txt")

}

func TestCanUseCustomMatchInSettingsFile(t *testing.T) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	curcmd := exec.Command("cp", "tests/assets/.flowcat1", dirname+"/.flowcat")
	_, err = curcmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	cmd := exec.Command("go", "run", "main.go", "-f", "tests/assets/multitest/")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	//fmt.Println(string(stdout))
	want := "tests/assets/multitest/test.php // Comment no todo"
	cur := string(stdout)
	ret := strings.Replace(cur, "\n", "", -1)
	ret = strings.Replace(ret, "\t", " ", -1)
	if want != ret {
		t.Errorf("got %s, want %s", ret, want)
	}

	dirname, err = os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	curcmd = exec.Command("rm", "-f", dirname+"/.flowcat")
	_, err = curcmd.Output()
	if err != nil {
		fmt.Println(err)
	}

}
