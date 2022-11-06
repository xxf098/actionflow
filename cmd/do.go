package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/containerd/console"
	"github.com/rs/zerolog"
	"github.com/xxf098/dagflow"
	"github.com/xxf098/dagflow/cmd/logger"
	"github.com/xxf098/dagflow/plan"
	"golang.org/x/term"
)

// https://cuelang.org/docs/concepts/packages/#import-path
func Do(dir string, action string) {
	lg, ctx := setupLog()
	targetPath := getTargetPath([]string{action})
	daggerPlan, err := loadPlan(ctx, dir)
	if err != nil {
		lg.Fatal().Err(err).Msg("failed to load plan")
	}
	lg.Info().Msg("load plan")
	err = daggerPlan.Do(ctx, targetPath)
	if err != nil {
		lg.Fatal().Err(err).Msg("failed to exec plan")
	}
	lg.Info().Msg("finish plan")
}

func Flow(dir string, action string) {
	mainCue := path.Join(dir, "main.cue")
	fmt.Println(mainCue)
	v := loadFile(mainCue)
	iter, err := v.Fields()
	for iter.Next() {
		fmt.Println(iter.Label())
	}
	// setup log
	lg, ctx := setupLog()
	target := cue.ParsePath(fmt.Sprintf(`actions.%s`, action))
	runner := dagflow.NewRunner(target)
	err = runner.Run(ctx, v)
	if err != nil {
		lg.Fatal().Err(err).Msg("failed to execute plan")
	}
}

func setupLog() (zerolog.Logger, context.Context) {
	cfg := logger.LogConfig{
		Level:  "info", // panic fatal error warn info debug trace
		Format: "plain",
	}
	lg := logger.New(cfg)
	ctx := lg.WithContext(context.Background())
	var tty *logger.TTYOutput
	var tty2 *logger.TTYOutputV2
	var err error

	f := cfg.Format
	switch {
	case f == "tty" || f == "auto" && term.IsTerminal(int(os.Stdout.Fd())):
		tty, err = logger.NewTTYOutput(os.Stderr)
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to initialize TTY logger")
		}
		tty.Start()
		defer tty.Stop()

		lg = lg.Output(logger.TeeCloud(tty))
		ctx = lg.WithContext(ctx)

	case f == "tty2":
		// FIXME: dolanor: remove once it's more stable/debuggable
		f, err := ioutil.TempFile("/tmp", "dagger-console-*.log")
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to create TTY file logger")
		}
		defer func() {
			err := f.Close()
			if err != nil {
				lg.Fatal().Err(err).Msg("failed to close TTY file logger")
			}
		}()

		cons, err := console.ConsoleFromFile(os.Stderr)
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to create TTY console")
		}

		c := logger.ConsoleAdapter{Cons: cons, F: f}
		tty2, err = logger.NewTTYOutputConsole(&c)
		if err != nil {
			lg.Fatal().Err(err).Msg("failed to initialize TTYv2 logger")
		}
		tty2.Start()
		defer tty2.Stop()

		lg = lg.Output(logger.TeeCloud(tty2))
		ctx = lg.WithContext(ctx)

	}
	return lg, ctx
}

func loadFile(filePath string) cue.Value {
	ctx := cuecontext.New()
	entrypoints := []string{filePath}

	bis := load.Instances(entrypoints, nil)
	return ctx.BuildInstance(bis[0])
}

func loadPlan(ctx context.Context, planPath string) (*plan.Plan, error) {
	absPlanPath, err := filepath.Abs(planPath)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(absPlanPath)
	if err != nil {
		return nil, err
	}
	return plan.Load(ctx, plan.Config{
		Args: []string{planPath},
	})
}

func getTargetPath(args []string) cue.Path {
	selectors := []cue.Selector{plan.ActionSelector}
	for _, arg := range args {
		selectors = append(selectors, cue.Str(arg))
	}
	return cue.MakePath(selectors...)
}
