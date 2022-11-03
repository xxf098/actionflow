package dagflow

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue"
)

func TestWriteFile(t *testing.T) {
	Do("./testcues/writefile.cue", "hello")
}

func TestWriteFile1(t *testing.T) {
	Do("./testcues/writefile1.cue", "hello")
}

func TestExec(t *testing.T) {
	Do("./testcues/exec.cue", "touch")
}

func TestRmFile(t *testing.T) {
	_, err := Do("./testcues/rmfile.cue", "test.rmFile.verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRmFileWildcard(t *testing.T) {
	err := Flow("./testcues/rmfile.cue", "test.rmWildcard.verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRmDir(t *testing.T) {
	err := Flow("./testcues/rmfile.cue", "test.rmDir.verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMkDir1(t *testing.T) {
	err := Flow("./testcues/mkdir1.cue", "readChecker")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMkDirParent(t *testing.T) {
	err := Flow("./testcues/mkdir_parents.cue", "readChecker")
	if err != nil {
		t.Fatal(err)
	}
}
func TestMkdir(t *testing.T) {
	output, err := Do("./testcues/mkdir.cue", "hello")
	if err != nil {
		t.Fatal(err)
	}
	v := output.LookupPath(cue.ParsePath("output"))
	s, _ := v.String()
	fmt.Println(s)
}

func TestGitPull(t *testing.T) {
	Do("./testcues/gitpull.cue", "hello")
}

func TestWorkDir(t *testing.T) {
	err := Flow("./testcues/exec_workdir.cue", "verify")
	if err != nil {
		t.Fatal(err)
	}
}
