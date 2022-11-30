package actionflow

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue"
)

func TestWriteFile(t *testing.T) {
	doTest("./testcues/writefile.cue", "hello")
}

func TestWriteFile1(t *testing.T) {
	flowTest("./testcues/writefile1.cue", "hello")
}

func TestWriteFile2(t *testing.T) {
	flowTest("./testcues/writefile2.cue", "hello")
}

func TestExec(t *testing.T) {
	doTest("./testcues/exec.cue", "touch")
}

func TestRun(t *testing.T) {
	doTest("./testcues/exec.cue", "touchRun")
}

func TestRmFile(t *testing.T) {
	_, err := doTest("./testcues/rmfile.cue", "test.rmFile.verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRmFileWildcard(t *testing.T) {
	err := flowTest("./testcues/rmfile.cue", "test.rmWildcard.verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRmDir(t *testing.T) {
	err := flowTest("./testcues/rmfile.cue", "test.rmDir.verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRmMulti(t *testing.T) {
	err := flowTest("./testcues/rmfile.cue", "test.rmMulti.rm")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMkDir1(t *testing.T) {
	err := flowTest("./testcues/mkdir1.cue", "readChecker")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMkDirParent(t *testing.T) {
	err := flowTest("./testcues/mkdir_parents.cue", "readChecker")
	if err != nil {
		t.Fatal(err)
	}
}
func TestMkdir(t *testing.T) {
	output, err := doTest("./testcues/mkdir.cue", "hello")
	if err != nil {
		t.Fatal(err)
	}
	v := output.LookupPath(cue.ParsePath("output"))
	s, _ := v.String()
	fmt.Println(s)
}

func TestGitPull(t *testing.T) {
	doTest("./testcues/git.cue", "hello")
}

func TestGitCheckout(t *testing.T) {
	doTest("./testcues/git.cue", "lite")
}

func TestWorkDir(t *testing.T) {
	err := flowTest("./testcues/exec/workdir.cue", "verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestEnv(t *testing.T) {
	err := flowTest("./testcues/exec/env.cue", "verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestArgs(t *testing.T) {
	err := flowTest("./testcues/exec/args.cue", "verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestOutput(t *testing.T) {
	err := flowTest("./testcues/exec/output.cue", "save")
	if err != nil {
		t.Fatal(err)
	}
}

func TestThen(t *testing.T) {
	err := flowTest("./testcues/then.cue", "mkdir")
	if err != nil {
		t.Fatal(err)
	}
}

func TestThen1(t *testing.T) {
	err := flowTest("./testcues/then.cue", "gitrm")
	if err != nil {
		t.Fatal(err)
	}
}

func TestThen2(t *testing.T) {
	err := flowTest("./testcues/then.cue", "clone")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAll(t *testing.T) {
	err := flowTest("./testcues/all.cue", "writeAll")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStep(t *testing.T) {
	err := flowTest("./testcues/step.cue", "setup")
	if err != nil {
		t.Fatal(err)
	}
}
