package task

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog/log"
	"github.com/xxf098/actionflow/compiler"
)

func init() {
	Register("Keep", func() Task { return &keepTask{} })
}

type keepTask struct {
}

func (t *keepTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	paths := []string{}
	pValue := v.LookupPath(cue.ParsePath("path"))
	lg := log.Ctx(ctx)
	start := time.Now()
	if pValue.IncompleteKind().IsAnyOf(cue.ListKind) {
		iter, err := pValue.List()
		if err != nil {
			return nil, err
		}
		for iter.Next() {
			path, err := iter.Value().String()
			if err != nil {
				return nil, err
			}

			paths = append(paths, path)
		}
	} else {
		path, err := pValue.String()
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	fullPaths := []string{}
	for _, path := range paths {
		if p, err := filepath.Abs(path); err == nil {
			fullPaths = append(fullPaths, p)
		}
	}

	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg("get root dirs")
	// get root dirs
	dirs := []string{}
	for _, v := range fullPaths {
		dir := filepath.Dir(v)
		found := false
		for _, dir1 := range dirs {
			if dir1 == dir {
				found = true
				break
			}
		}
		if !found {
			dirs = append(dirs, dir)
		}
	}

	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg("skip")
	// skip
	if len(fullPaths) < 1 || len(dirs) < 1 {
		value := compiler.NewValue()
		output := value.FillPath(cue.ParsePath("output"), ".")
		return &output, nil
	}

	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg("sort dir by length")
	// sort dir by length
	sort.Slice(dirs, func(i, j int) bool {
		return len(dirs[i]) < len(dirs[j])
	})

	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg("glob remove files")
	// glob remove files
	for _, dir := range dirs {
		lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(dir)
		patterns := []string{}
		for _, path := range fullPaths {
			if strings.HasPrefix(path, dir) {
				lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(path)
				patterns = append(patterns, path)
			}
		}
		if len(patterns) < 1 {
			continue
		}
		removeFiles, err := reverseGlob(dir, patterns)
		if err != nil {
			continue
		}
		for _, f := range removeFiles {
			lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(f)
			os.RemoveAll(f)
		}
	}

	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(t.Name())
	Then(ctx, v)
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), fullPaths[0])
	return &output, nil
}

func (t *keepTask) Name() string {
	return "Keep"
}

func reverseGlob(rootDir string, patterns []string) ([]string, error) {
	paths := []string{}
	filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if m, err := filepath.Match(rootDir, path); err != nil || m {
			return nil
		}

		for _, p := range patterns {
			m, err := filepath.Match(p, path)
			if err != nil || m {
				return nil
			}
		}
		paths = append(paths, path)
		return nil
	})
	return paths, nil
}
