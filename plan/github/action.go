package github

import (
	"context"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/xxf098/actionflow/plan/github/model"
)

func ReadAction(ctx context.Context, actionDir string) (*model.Action, error) {
	actionPath := path.Join(actionDir, "action.yml")
	f, err := os.Open(actionPath)
	if os.IsNotExist(err) {
		f, err = os.Open(actionPath)
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return model.ReadAction(f)
}

func ActionCacheDir() string {
	var xdgCache string
	var ok bool
	if xdgCache, ok = os.LookupEnv("XDG_CACHE_HOME"); !ok || xdgCache == "" {
		if home, err := homedir.Dir(); err == nil {
			xdgCache = filepath.Join(home, ".cache")
		} else if xdgCache, err = filepath.Abs("."); err != nil {
			log.Fatal(err)
		}
	}
	return filepath.Join(xdgCache, "flow")
}
