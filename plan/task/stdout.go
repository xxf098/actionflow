package task

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog/log"
	"github.com/xxf098/actionflow/compiler"
)

// @stdout
func init() {
	Register("Stdout", func() Task { return &Stdout{} })
}

type Stdout struct {
}

func (t *Stdout) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	lg := log.Ctx(ctx)
	start := time.Now()
	text, err := v.LookupPath(cue.ParsePath("text")).String()
	if err != nil {
		return nil, err
	}

	bufStdout := bufio.NewWriter(os.Stdout)
	defer bufStdout.Flush()
	fmt.Fprintln(bufStdout, text)

	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(t.Name())
	Then(ctx, v)
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), "")
	return &output, nil
}

func (t *Stdout) Name() string {
	return "Stdout"
}

func attrStdout(v *cue.Value, text string) {
	// read attr
	attrs := v.Attributes(cue.ValueAttr)
	for _, attr := range attrs {
		if attr.Name() == "stdout" {
			bufStdout := bufio.NewWriter(os.Stdout)
			defer bufStdout.Flush()
			fmt.Fprintln(bufStdout, text)
			break
		}
	}
}
