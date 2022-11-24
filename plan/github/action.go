package github

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/xxf098/actionflow/plan/github/model"
)

type actionStep interface {
	step

	getActionModel() *model.Action
	// getCompositeRunContext(context.Context) *RunContext
	// getCompositeSteps() *compositeSteps
}

func readActionImpl(ctx context.Context, actionDir string) (*model.Action, error) {
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

func actionCacheDir() string {
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

func runActionImpl(ctx context.Context, step actionStep, actionDir string, remoteAction *remoteAction) error {
	stepModel := step.getStepModel() // workflow.yml

	action := step.getActionModel() // action.yml
	actionPath := path.Join(actionDir, action.Runs.Main)
	cmd := exec.CommandContext(ctx, "node", actionPath)
	envs := setupActionEnv(ctx, step)
	cmd.Env = append(cmd.Env, envs...)
	log.Printf("type=%v actionDir=%s actionPath=%s\n envs=%s", stepModel.Type(), actionDir, actionPath, strings.Join(envs, " "))
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		log.Println(outBuf.String())
		err = fmt.Errorf("%s: %s", err.Error(), errBuf.String())
	}
	// log.Println(outBuf.String())
	return err
}

func setupActionEnv(ctx context.Context, step actionStep) []string {

	// populateEnvsFromInput(ctx, step.getEnv(), step.getActionModel())
	envs := []string{}
	envs = setupGithubEnv(envs)
	stepEnvs := step.getEnv()
	populateEnvsFromInput(ctx, stepEnvs, step.getActionModel())
	for k, v := range *stepEnvs {
		// FIXME
		if strings.HasPrefix(v, "${{") && strings.HasSuffix(v, "}}") {
			continue
		}
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	envs = setupRunnerEnv(envs)
	return envs
}

func setupGithubEnv(envs []string) []string {
	envs = append(envs, "CI=true")
	envs = append(envs, "GITHUB_ACTIONS=true")
	envs = append(envs, "GITHUB_SERVER_URL=https://github.com")
	envs = append(envs, "GITHUB_API_URL=https://api.github.com")
	envs = append(envs, "GITHUB_GRAPHQL_URL=https://api.github.com/graphql")
	keys := []string{"GITHUB_ENV", "GITHUB_PATH", "GITHUB_WORKFLOW", "GITHUB_RUN_ID", "GITHUB_RUN_NUMBER", "GITHUB_REPOSITORY", "GITHUB_EVENT_NAME",
		"GITHUB_EVENT_PATH", "GITHUB_WORKSPACE", "GITHUB_SHA", "GITHUB_REF", "GITHUB_REF_NAME", "GITHUB_REF_TYPE", "GITHUB_TOKEN", "GITHUB_JOB",
		"GITHUB_REPOSITORY_OWNER", "GITHUB_RETENTION_DAYS", "RUNNER_PERFLOG", "RUNNER_TRACKING_ID", "XDG_CONFIG_HOME"}
	for _, key := range keys {
		envs = append(envs, fmt.Sprintf("%s=%s", key, os.Getenv(key)))
	}
	return envs
}

func setupRunnerEnv(envs []string) []string {
	if !checkEnv(envs, "RUNNER_TOOL_CACHE") {
		envs = append(envs, fmt.Sprintf("%s=%s", "RUNNER_TOOL_CACHE", "/opt/hostedtoolcache"))
	}
	if !checkEnv(envs, "RUNNER_OS") {
		envs = append(envs, fmt.Sprintf("%s=%s", "RUNNER_OS", "Linux"))
	}
	if !checkEnv(envs, "RUNNER_ARCH") {
		goarch := runtime.GOARCH
		if goarch == "amd64" {
			goarch = "x64"
		}
		if goarch == "386" {
			goarch = "x86"
		}
		envs = append(envs, fmt.Sprintf("%s=%s", "RUNNER_ARCH", goarch))
	}
	if !checkEnv(envs, "RUNNER_TEMP") {
		envs = append(envs, fmt.Sprintf("%s=%s", "RUNNER_TEMP", "/tmp"))
	}
	return envs
}

func checkEnv(envs []string, key string) bool {
	for _, v := range envs {
		if strings.HasPrefix(v, key+"=") {
			return true
		}
	}
	return false
}

// setup input
// FIXME eval
func populateEnvsFromInput(ctx context.Context, env *map[string]string, action *model.Action) {
	for inputID, input := range action.Inputs {
		envKey := regexp.MustCompile("[^A-Z0-9-]").ReplaceAllString(strings.ToUpper(inputID), "_")
		envKey = fmt.Sprintf("INPUT_%s", envKey)
		if _, ok := (*env)[envKey]; !ok {
			(*env)[envKey] = input.Default
		}
	}
}
