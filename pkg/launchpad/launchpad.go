// Copyright 2021 - 2021, Adolfo García Veytia and the moonport contributors
// SPDX-License-Identifier: Apache-2.0

package launchpad

import (
	"github.com/pkg/errors"
	"github.com/uServers/moonport/pkg/backend"
	"github.com/uServers/moonport/pkg/pipeline"
)

// New Creates a new Launchpad
func New() (*LaunchPad, error) {
	b, err := backend.NewBuildBackend("gcb")
	if err != nil {
		return nil, errors.Wrap(err, "creating build backend")
	}
	return &LaunchPad{
		backend: b,
		impl:    &defaultImplementation{},
	}, nil
}

type Options struct {
	StepSources []string
}

// Options returns the launchpad options
func (lp *LaunchPad) Options() *Options {
	return lp.options
}

type LaunchPad struct {
	Steps   map[string]*pipeline.Step
	options *Options
	backend *backend.BuildBackend
	impl    Implementation
}

// Run submnits a pipeline
func (lp *LaunchPad) Run(p *pipeline.Pipeline) (*backend.JobData, error) {
	return lp.impl.Run(lp.backend, p, lp.options)
}

// Implementation
type Implementation interface {
	Run(*backend.BuildBackend, *pipeline.Pipeline, *Options) (*backend.JobData, error)
}
