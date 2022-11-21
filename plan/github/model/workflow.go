package model

import "gopkg.in/yaml.v3"

// Step is the structure of one step in a job
type Step struct {
	ID                 string            `yaml:"id"`
	If                 yaml.Node         `yaml:"if"`
	Name               string            `yaml:"name"`
	Uses               string            `yaml:"uses"`
	Run                string            `yaml:"run"`
	WorkingDirectory   string            `yaml:"working-directory"`
	Shell              string            `yaml:"shell"`
	Env                yaml.Node         `yaml:"env"`
	With               map[string]string `yaml:"with"`
	RawContinueOnError string            `yaml:"continue-on-error"`
	TimeoutMinutes     string            `yaml:"timeout-minutes"`
}
