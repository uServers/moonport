// Copyright 2021 - 2021, Adolfo García Veytia and the moonport contributors
// SPDX-License-Identifier: Apache-2.0

package launchpad

import (
	"github.com/pkg/errors"
	"github.com/uServers/moonport/pkg/backend"
	"github.com/uServers/moonport/pkg/pipeline"
)

// Creates a new Launchpad
func New() (*LaunchPad, error) {
	b, err := backend.NewBuildBackend("gcb")
	if err != nil {
		return nil, errors.Wrap(err, "creating build backend")
	}
	return &LaunchPad{
		backend: b,
	}, nil
}

type LaunchPad struct {
	backend *backend.BuildBackend
}

func (lp *LaunchPad) Run(p *pipeline.Pipeline) (*backend.JobData, error) {
	if lp.backend == nil {
		return nil, errors.New("Launchpad does not have a valid backend")
	}
	return lp.backend.RunPipeline(p)
}
