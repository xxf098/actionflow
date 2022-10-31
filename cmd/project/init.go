package project

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func Init(ctx context.Context, parentDir, module string) error {
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

	return nil
}
