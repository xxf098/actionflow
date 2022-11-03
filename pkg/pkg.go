package pkg

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gofrs/flock"
	"github.com/rs/zerolog/log"
)

// ln -s ~/github/dagflow/pkg/github.com/xxf098/dagflow dagflow
var (
	// FS contains the filesystem of the stdlib.
	//go:embed github.com
	FS embed.FS
)

var (
	lockFilePath = "dagger.lock"
)

func Vendor(ctx context.Context, p string) error {
	if p == "" {
		p, _ = GetCueModParent()
	}
	cuePkgDir := path.Join(p, "cue.mod", "pkg")
	if err := os.MkdirAll(cuePkgDir, 0755); err != nil {
		return err
	}

	// Lock this function so no more than 1 process can run it at once.
	lockFile := path.Join(cuePkgDir, lockFilePath)
	l := flock.New(lockFile)
	if err := l.Lock(); err != nil {
		return err
	}
	defer func() {
		l.Unlock()
		os.Remove(lockFile)
	}()

	// ensure cue module is initialized
	if err := CueModInit(ctx, p, ""); err != nil {
		return err
	}

	if err := extractModules(cuePkgDir); err != nil {
		return err
	}

	return nil
}

// extract files
func extractModules(dest string) error {
	return fs.WalkDir(FS, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.Type().IsRegular() {
			return nil
		}

		// Do not vendor the package's `cue.mod/pkg`
		if strings.Contains(p, "cue.mod/pkg") {
			return nil
		}

		contents, err := fs.ReadFile(FS, p)
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}

		overlayPath := path.Join(dest, p)
		// fmt.Println(overlayPath)

		if err := os.MkdirAll(filepath.Dir(overlayPath), 0755); err != nil {
			return err
		}

		// Give exec permission on embedded file to freely use shell script
		// Exclude permission linter
		//nolint
		return os.WriteFile(overlayPath, contents, 0700)
	})
}

// GetCueModParent traverses the directory tree up through ancestors looking for a cue.mod folder
func GetCueModParent(args ...string) (string, bool) {
	cwd, _ := os.Getwd()
	parentDir := cwd

	if len(args) == 1 {
		parentDir = args[0]
	}

	found := false

	for {
		if _, err := os.Stat(path.Join(parentDir, "cue.mod")); !errors.Is(err, os.ErrNotExist) {
			found = true
			break // found it!
		}

		parentDir = filepath.Dir(parentDir)

		if parentDir == fmt.Sprintf("%s%s", filepath.VolumeName(parentDir), string(os.PathSeparator)) {
			// reached the root
			parentDir = cwd // reset to working directory
			break
		}
	}

	return parentDir, found
}

func CueModInit(ctx context.Context, parentDir, module string) error {
	lg := log.Ctx(ctx)

	absParentDir, err := filepath.Abs(parentDir)
	if err != nil {
		return err
	}

	modDir := path.Join(absParentDir, "cue.mod")
	if err := os.MkdirAll(modDir, 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	modFile := path.Join(modDir, "module.cue")
	if _, err := os.Stat(modFile); err != nil {
		statErr, ok := err.(*os.PathError)
		if !ok {
			return statErr
		}

		lg.Debug().Str("mod", parentDir).Msg("initializing cue.mod")
		contents := fmt.Sprintf(`module: "%s"`, module)
		if err := os.WriteFile(modFile, []byte(contents), 0600); err != nil {
			return err
		}
	}

	if err := os.Mkdir(path.Join(modDir, "pkg"), 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	mainFile := path.Join(absParentDir, "main.cue")
	if _, err := os.Stat(mainFile); err != nil {
		statErr, ok := err.(*os.PathError)
		if !ok {
			return statErr
		}

		lg.Debug().Str("mod", parentDir).Msg("initializing main.cue")
		contents := `package main
import (
	"github.com/xxf098/dagflow"
	"github.com/xxf098/dagflow/core"
)

dagflow.#Plan & {
	actions: {

		mkdir: core.#Mkdir & {
			path:  "./hello"
		}
		
	}
}		
`
		if err := os.WriteFile(mainFile, []byte(contents), 0600); err != nil {
			return err
		}
	}

	return nil
}
