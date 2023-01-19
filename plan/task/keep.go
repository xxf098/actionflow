package task

import (
	"context"
	"fmt"
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

// keep files: /abc/def/*.txt
// keep folder: /abc/def/hij/
// FIXME: conflict
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
			// keep path Separator
			if os.IsPathSeparator(path[len(path)-1]) {
				p = fmt.Sprintf("%s%s", p, path[len(path)-1:])
			}
			fullPaths = append(fullPaths, p)
		}
	}

	// get root dirs
	dirs := []string{}
	for _, v := range fullPaths {
		dir := filepath.Dir(v)
		// set /abc/def/ parent to /abc
		if len(v) > 1 && os.IsPathSeparator(v[len(v)-1]) {
			dir = filepath.Dir(dir)
		}
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

	// skip
	if len(fullPaths) < 1 || len(dirs) < 1 {
		value := compiler.NewValue()
		output := value.FillPath(cue.ParsePath("output"), ".")
		return &output, nil
	}

	// sort dir by length
	sort.Slice(dirs, func(i, j int) bool {
		return len(dirs[i]) < len(dirs[j])
	})

	// glob remove files
	for _, dir := range dirs {
		patterns := []string{}
		for _, path := range fullPaths {
			if strings.HasPrefix(path, dir) {
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
	dir := ":"
	filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if strings.HasPrefix(path, dir) {
			return nil
		}
		if m, err := filepath.Match(rootDir, path); err != nil || m {
			return nil
		}

		for _, p := range patterns {
			m, err := filepath.Match(p, path)
			if err != nil || m {
				return nil
			}
			// /abc/def == /abc/def/
			if os.IsPathSeparator(p[len(p)-1]) {
				if strings.HasPrefix(path, p) || (d.IsDir() && fmt.Sprintf("%s%s", path, p[len(p)-1:]) == p) {
					return nil
				}
			}
		}
		if d.IsDir() {
			dir = fmt.Sprintf("%s%c", path, filepath.Separator)
		}
		paths = append(paths, path)
		return nil
	})
	return paths, nil
}
