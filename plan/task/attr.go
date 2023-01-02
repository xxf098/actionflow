package task

import (
	"bufio"
	"fmt"
	"os"

	"cuelang.org/go/cue"
)

func debug(v *cue.Value, text string) {
	// read attr
	attrs := v.Attributes(cue.ValueAttr)
	for _, attr := range attrs {
		if attr.Name() == "debug" {
			bufStdout := bufio.NewWriter(os.Stdout)
			defer bufStdout.Flush()
			fmt.Fprintln(bufStdout, text)
			break
		}
	}
}
