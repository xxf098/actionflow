package model

import (
	"fmt"
	"regexp"
	"strings"
)

// Step is the structure of one step in a job
type Step struct {
	ID                 string            `yaml:"id"`
	If                 bool              `yaml:"if"`
	Name               string            `yaml:"name"`
	Uses               string            `yaml:"uses"`
	Run                string            `yaml:"run"`
	WorkingDirectory   string            `yaml:"working-directory"`
	Shell              string            `yaml:"shell"`
	Env                map[string]string `yaml:"env"`
	With               map[string]string `yaml:"with"`
	RawContinueOnError string            `yaml:"continue-on-error"`
	TimeoutMinutes     string            `yaml:"timeout-minutes"`
}

func NewStep(uses string, with map[string]string) Step {
	return Step{
		Uses: uses,
		With: with,
	}
}

// String gets the name of step
func (s *Step) String() string {
	if s.Name != "" {
		return s.Name
	} else if s.Uses != "" {
		return s.Uses
	} else if s.Run != "" {
		return s.Run
	}
	return s.ID
}

// Environments returns string-based key=value map for a step
// Note: all keys are uppercase
func (s *Step) Environment() map[string]string {
	env := map[string]string{}
	for k, v := range env {
		delete(env, k)
		env[strings.ToUpper(k)] = v
	}
	return env
}

// GetEnv gets the env for a step
// getenv from input
func (s *Step) GetEnv() map[string]string {
	env := s.Environment()

	for k, v := range s.With {
		envKey := regexp.MustCompile("[^A-Z0-9-]").ReplaceAllString(strings.ToUpper(k), "_")
		envKey = fmt.Sprintf("INPUT_%s", strings.ToUpper(envKey))
		env[envKey] = v
	}
	return env
}

// StepType describes what type of step we are about to run
type StepType int

const (
	// StepTypeRun is all steps that have a `run` attribute run with docker exec
	StepTypeRun StepType = iota

	// StepTypeUsesDockerURL is all steps that have a `uses` that is of the form `docker://...`
	StepTypeUsesDockerURL

	// StepTypeUsesActionLocal is all steps that have a `uses` that is a local action in a subdirectory
	StepTypeUsesActionLocal

	// StepTypeUsesActionRemote is all steps that have a `uses` that is a reference to a github repo  run with node index.js
	StepTypeUsesActionRemote

	// StepTypeReusableWorkflowLocal is all steps that have a `uses` that is a local workflow in the .github/workflows directory
	StepTypeReusableWorkflowLocal

	// StepTypeReusableWorkflowRemote is all steps that have a `uses` that references a workflow file in a github repo
	StepTypeReusableWorkflowRemote

	// StepTypeInvalid is for steps that have invalid step action
	StepTypeInvalid
)

func (s StepType) String() string {
	switch s {
	case StepTypeInvalid:
		return "invalid"
	case StepTypeRun:
		return "run"
	case StepTypeUsesActionLocal:
		return "local-action"
	case StepTypeUsesActionRemote:
		return "remote-action"
	case StepTypeUsesDockerURL:
		return "docker"
	case StepTypeReusableWorkflowLocal:
		return "local-reusable-workflow"
	case StepTypeReusableWorkflowRemote:
		return "remote-reusable-workflow"
	}
	return "unknown"
}

// Type returns the type of the step
func (s *Step) Type() StepType {
	if s.Run == "" && s.Uses == "" {
		return StepTypeInvalid
	}

	if s.Run != "" {
		if s.Uses != "" {
			return StepTypeInvalid
		}
		return StepTypeRun
	} else if strings.HasPrefix(s.Uses, "docker://") {
		return StepTypeUsesDockerURL
	} else if strings.HasPrefix(s.Uses, "./.github/workflows") && (strings.HasSuffix(s.Uses, ".yml") || strings.HasSuffix(s.Uses, ".yaml")) {
		return StepTypeReusableWorkflowLocal
	} else if !strings.HasPrefix(s.Uses, "./") && strings.Contains(s.Uses, ".github/workflows") && (strings.Contains(s.Uses, ".yml@") || strings.Contains(s.Uses, ".yaml@")) {
		return StepTypeReusableWorkflowRemote
	} else if strings.HasPrefix(s.Uses, "./") {
		return StepTypeUsesActionLocal
	}
	return StepTypeUsesActionRemote
}
