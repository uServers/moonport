// Copyright 2021 - 2021, Adolfo García Veytia and the moonport contributors
// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uServers/moonport/pkg/metadata"
	"gopkg.in/yaml.v2"
)

// Pipeline abstracts a pipeline
type Pipeline struct {
	Metadata metadata.Metadata `yaml:"metadata"`
	Spec     struct {
		Repository     string            `yaml:"repo"`
		ServiceAccount string            `yaml:"service-account"`
		EnvVars        map[string]string `yaml:"env"`
		Stages         map[string]*Stage `yaml:"stages"`
	} `yaml:"spec"`
	stepCatalog map[string]*Step
	options     *Options
	impl        Implementation
}

// Options return the pipeline options set
func (p *Pipeline) Options() *Options {
	return p.options
}

type Options struct {
	StepSources []string
}

var defaultPipelineOpts = Options{}

// Stage binds together a series of jobs
type Stage struct {
	Comment string     `yaml:"comment"`
	Steps   []StepSpec `yaml:"steps"`
}

type StepSpec struct {
	StepLabel string `yaml:"step"`
}

// Step abstracts a job template defined in a file
type Step struct {
	Metadata metadata.Metadata `yaml:"metadata"`
	Spec     struct {
		Image  string   `yaml:"image"`
		Bundle string   `yaml:"bundle"`
		Args   []string `yaml:"args"`
	} `yaml:"spec"`
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
func NewPipelineFromFile(path string, opts *Options) (p *Pipeline, err error) {
	p = NewPipeline(opts)
	yamlData, err := os.ReadFile(path)
	if err != nil {
		return p, errors.Wrap(err, "reading yaml file")
	}
	// Unmarshall
	if err = yaml.Unmarshal(yamlData, &p); err != nil {
		return p, errors.Wrap(err, "deciding yaml data")
	}

	logrus.Infof("Successfully loaded pipeline %s (%d stages)", p.Metadata.Name, len(p.Spec.Stages))
	return p, nil
}

// Validate checks if the pipeline can run
func (p *Pipeline) Validate() error {
	// Check the steps
	if err := p.loadStepCatalog(); err != nil {
		return errors.Wrap(err, "loading step catalog")
	}
	if len(p.stepCatalog) == 0 {
		return errors.New("step catalog is empty, this means that you pipeline does not know hoy to do anything")
	}
	for name, stage := range p.Spec.Stages {
		i := 0
		logrus.Infof("Verifying stage %s (%d steps)", name, len(stage.Steps))
		for _, stepspec := range stage.Steps {
			if stepspec.StepLabel == "" {
				return errors.New(fmt.Sprintf("step #%d from stage %s has no label", i, name))
			}
			i++
		}
	}
	return nil
}

func (p *Pipeline) GetStep(label string) *Step {
	if step, ok := p.stepCatalog[label]; ok {
		return step
	} else {
		return nil
	}
}

func (p *Pipeline) loadStepCatalog() error {
	steps, err := ReadStepsDirectories(p.options.StepSources)
	if err != nil {
		return errors.Wrap(err, "reading step drectories")
	}

	p.stepCatalog = steps
	return nil
}

func ReadStepsDirectories(paths []string) (steps map[string]*Step, err error) {
	steps = map[string]*Step{}
	for _, path := range paths {
		dirSteps, err := ReadStepsDirectory(path)
		if err != nil {
			return nil, errors.Wrapf(err, "reading files from %s", path)
		}
		for n, s := range dirSteps {
			if _, ok := steps[n]; ok {
				logrus.Warn("Duplicate step '%s' ", n)
			}
			steps[n] = s
		}
	}
	return steps, nil
}

// ReadStepsDirectory reads all the yaml files in a directory
// to look for step definitions
func ReadStepsDirectory(path string) (map[string]*Step, error) {
	steps := map[string]*Step{}
	fileList := []string{}
	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			fileList = append(fileList, path)
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "reading steps dir")
	}

	for _, filename := range fileList {
		filesteps, err := ReadStepsFromFile(filename)
		if err != nil {
			return nil, errors.Wrapf(err, "reading steps from %s", filename)
		}
		for name, step := range filesteps {
			steps[name] = step
		}
	}
	return steps, nil
}

// ReadStepsFromFile reads a yaml file and returns the step definitions found in it
func ReadStepsFromFile(path string) (map[string]*Step, error) {
	// Open the file
	yamlReader, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "opening path")
	}

	steps := map[string]*Step{}

	decoder := yaml.NewDecoder(yamlReader)
	i := 0
	for {
		step := Step{}
		if err := decoder.Decode(&step); err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.Wrap(err, "decoding step yaml code")
		}

		if step.Metadata.Name != "" {
			steps[step.Metadata.Name] = &step
		} else {
			logrus.Warning("Ignoring step #%d from %s as it does not have a name")
		}
		i++
	}

	return steps, nil
}

// PipelineInterface
type Implementation interface {
}

// defaultPipelineImplementation
type defaultImplementation struct {
}
