package actionflow

import (
	"fmt"
	"path/filepath"
	"testing"

	"cuelang.org/go/cue"
)

func TestWriteFile(t *testing.T) {
	err := doFlowTest("./testcues/writefile.cue", "hello")
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteFile1(t *testing.T) {
	doFlowTest("./testcues/writefile1.cue", "hello")
}

func TestWriteFile2(t *testing.T) {
	doFlowTest("./testcues/writefile2.cue", "hello")
}

func TestExec(t *testing.T) {
	err := doFlowTest("./testcues/exec.cue", "touch")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRun(t *testing.T) {
	err := doFlowTest("./testcues/exec.cue", "touchRun")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRmFile(t *testing.T) {
	err := doFlowTest("./testcues/rmfile.cue", "test.rmFile.verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRmFileWildcard(t *testing.T) {
	err := doFlowTest("./testcues/rmfile.cue", "test.rmWildcard.verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRmDir(t *testing.T) {
	err := doFlowTest("./testcues/rmfile.cue", "test.rmDir.verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRmMulti(t *testing.T) {
	err := doFlowTest("./testcues/rmfile.cue", "test.rmMulti.rm")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMkDir1(t *testing.T) {
	err := doFlowTest("./testcues/mkdir1.cue", "readChecker")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMkDirs(t *testing.T) {
	err := doFlowTest("./testcues/mkdirs.cue", "zero")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMkDirParent(t *testing.T) {
	err := doFlowTest("./testcues/mkdir_parents.cue", "readChecker")
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
	err := doFlowTest("./testcues/exec/workdir.cue", "verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestEnv(t *testing.T) {
	err := doFlowTest("./testcues/exec/env.cue", "verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestArgs(t *testing.T) {
	err := doFlowTest("./testcues/exec/args.cue", "verify")
	if err != nil {
		t.Fatal(err)
	}
}

func TestOutput(t *testing.T) {
	err := doFlowTest("./testcues/exec/output.cue", "save")
	if err != nil {
		t.Fatal(err)
	}
}

func TestThen(t *testing.T) {
	err := doFlowTest("./testcues/then.cue", "mkdir")
	if err != nil {
		t.Fatal(err)
	}
}

func TestThen1(t *testing.T) {
	err := doFlowTest("./testcues/then.cue", "gitrm")
	if err != nil {
		t.Fatal(err)
	}
}

func TestThen2(t *testing.T) {
	err := doFlowTest("./testcues/then.cue", "clone")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAll(t *testing.T) {
	err := doFlowTest("./testcues/all.cue", "writeAll")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStep(t *testing.T) {
	err := doFlowTest("./testcues/step.cue", "setup")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStdout(t *testing.T) {
	err := doFlowTest("./testcues/stdout.cue", "print")
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadfile(t *testing.T) {
	err := doFlowTest("./testcues/attr.cue", "hellofile")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAPICall(t *testing.T) {
	err := doFlowTest("./testcues/api.cue", "print")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeps(t *testing.T) {
	err := doFlowTest("./testcues/deps.cue", "read")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeps1(t *testing.T) {
	err := doFlowTest("./testcues/deps1.cue", "all")
	if err != nil {
		t.Fatal(err)
	}
}

func TestKeep(t *testing.T) {
	err := doFlowTest("./testcues/keep.cue", "keepFile")
	if err != nil {
		t.Fatal(err)
	}
}

func TestKeep1(t *testing.T) {
	err := doFlowTest("./testcues/keep.cue", "keepSub")
	if err != nil {
		t.Fatal(err)
	}
}

func TestKeep2(t *testing.T) {
	p := "/home/abc/github/def/sub/trials/*"
	name := "/home/abc/github/def/sub/trials/abc.txt"
	m, err := filepath.Match(p, name)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m)
}
