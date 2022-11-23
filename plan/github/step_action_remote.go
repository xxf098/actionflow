package github

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/xxf098/actionflow/plan/github/model"
)

type StepActionRemote struct {
	Step         *model.Step
	actionDir    string
	action       *model.Action
	remoteAction *remoteAction
	env          map[string]string
}

func newStepActionRemote(step *model.Step) *StepActionRemote {
	return &StepActionRemote{Step: step}
}

func (sar *StepActionRemote) pre(ctx context.Context) error {
	remoteAction := newRemoteAction(sar.Step.Uses)
	if remoteAction == nil {
		return fmt.Errorf("Expected format {org}/{repo}[/path]@ref. Actual '%s' Input string was not in a correct format", sar.Step.Uses)
	}
	actionDir := fmt.Sprintf("%s/%s", actionCacheDir(), strings.ReplaceAll(sar.Step.Uses, "/", "-"))
	cloneURL := fmt.Sprintf("%s.git", remoteAction.CloneURL())
	cmd := exec.CommandContext(ctx, "git", "clone", cloneURL, actionDir)
	if err := cmd.Run(); err != nil {
		return err
	}
	// read action
	action, err := readActionImpl(ctx, actionDir)
	if err != nil {
		return err
	}
	sar.remoteAction = remoteAction
	sar.actionDir = actionDir
	sar.action = action
	sar.env = sar.Step.GetEnv()
	return nil
}

// FIXME: output
func (sar *StepActionRemote) main(ctx context.Context) error {
	return runActionImpl(ctx, sar, sar.actionDir, sar.remoteAction)
}

func (sar *StepActionRemote) post(ctx context.Context) error {
	return nil
}

func (sar *StepActionRemote) getActionModel() *model.Action {
	return sar.action
}

func (sar *StepActionRemote) getStepModel() *model.Step {
	return sar.Step
}

func (sar *StepActionRemote) getEnv() *map[string]string {
	return &sar.env
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
