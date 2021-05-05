// Copyright 2021 - 2021, Adolfo García Veytia and the moonport contributors
// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Pipeline abstracts a pipeline
type Pipeline struct {
	Repository     string            `yaml:"repo"`
	ServiceAccount string            `yaml:"service-account"`
	EnvVars        map[string]string `yaml:"env"`
	Stages         map[string]*Stage `yaml:"stages"`
	options        *Options
	impl           Implementation
}

// Options return the pipeline options set
func (p *Pipeline) Options() *Options {
	return p.options
}

type Options struct{}

var defaultPipelineOpts = Options{}

// Stage binds together a series of jobs
type Stage struct {
	Steps []Step
}

// Step abstracts a specific job in the pipeline
type Step struct {
	Image  string
	Bundle string
	Args   []string
}

// NewPipeline return a new Pipeline
func NewPipeline(opts *Options) *Pipeline {
	if opts == nil {
		opts = &defaultPipelineOpts
	}
	return &Pipeline{
		options: opts,
		impl:    &defaultImplementation{},
	}
}

// NewPipelineFromFile reads a pipeline configuration file
func NewPipelineFromFile(path string, opts *Options) (pipeline *Pipeline, err error) {
	pipeline = NewPipeline(opts)
	yamlData, err := os.ReadFile(path)
	if err != nil {
		return pipeline, errors.Wrap(err, "reading yaml file")
	}
	if err = yaml.Unmarshal(yamlData, &pipeline); err != nil {
		return pipeline, errors.Wrap(err, "deciding yaml data")
	}

	return pipeline, nil
}

// PipelineInterface
type Implementation interface {
}

// defaultPipelineImplementation
type defaultImplementation struct {
}
