package task

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"cuelang.org/go/cue"
	"github.com/xxf098/actionflow/plan/github"
)

func init() {
	Register("Step", func() Task { return &stepTask{} })
}

// run github step
type stepTask struct {
}

func (t *stepTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	usesValue := v.Lookup("uses")
	uses, err := usesValue.String()
	if err != nil {
		return nil, err
	}
	withValue := v.Lookup("with")
	envs := []string{}
	if withValue.Exists() {
		ik := withValue.IncompleteKind()
		if !(ik.IsAnyOf(cue.StructKind) && v.IsConcrete()) {
			return nil, errors.New("")
		}
		iter, _ := withValue.Fields()
		for iter.Next() {
			label := iter.Label()
			if v, err := iter.Value().String(); err == nil {
				envs = append(envs, fmt.Sprintf("%s=%s", label, v))
				continue
			} else if v, err := iter.Value().Bool(); err == nil {
				envs = append(envs, fmt.Sprintf("%s=%t", label, v))
				continue
			} else if v, err := iter.Value().Int64(); err == nil {
				envs = append(envs, fmt.Sprintf("%s=%d", label, v))
				continue
			} else if v, err := iter.Value().Float64(); err == nil {
				envs = append(envs, fmt.Sprintf("%s=%f", label, v))
				continue
			}
			return nil, fmt.Errorf("wrong field %s", label)
		}
	}
	_, err = t.clone(ctx, uses, envs)
	if err != nil {
		return nil, err
	}
	// read

	return nil, nil
}

func (t *stepTask) Name() string {
	return "Step"
}

// FIXME: git pull
func (t *stepTask) clone(ctx context.Context, uses string, envs []string) (string, error) {
	remoteAction := newRemoteAction(uses)
	actionDir := fmt.Sprintf("%s/%s", github.ActionCacheDir(), strings.ReplaceAll(uses, "/", "-"))
	cloneURL := fmt.Sprintf("%s.git", remoteAction.CloneURL())
	// git ref
	// git --git-dir=./setup-go/.git checkout v3
	cmd := exec.CommandContext(ctx, "git", "clone", cloneURL, actionDir)
	// cmd.Env = append(cmd.Env, envs...)
	return actionDir, cmd.Run()
}

type remoteAction struct {
	URL  string
	Org  string
	Repo string
	Path string
	Ref  string
}

func (ra *remoteAction) CloneURL() string {
	return fmt.Sprintf("https://%s/%s/%s", ra.URL, ra.Org, ra.Repo)
}

func newRemoteAction(action string) *remoteAction {
	// GitHub's document[^] describes:
	// > We strongly recommend that you include the version of
	// > the action you are using by specifying a Git ref, SHA, or Docker tag number.
	// Actually, the workflow stops if there is the uses directive that hasn't @ref.
	// [^]: https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions
	r := regexp.MustCompile(`^([^/@]+)/([^/@]+)(/([^@]*))?(@(.*))?$`)
	matches := r.FindStringSubmatch(action)
	if len(matches) < 7 || matches[6] == "" {
		return nil
	}
	return &remoteAction{
		Org:  matches[1],
		Repo: matches[2],
		Path: matches[4],
		Ref:  matches[6],
		URL:  "github.com",
	}
}
