// Copyright 2021 - 2021, Adolfo García Veytia and the moonport contributors
// SPDX-License-Identifier: Apache-2.0

package launchpad

import (
	"github.com/uServers/moonport/pkg/backend"
	"github.com/uServers/moonport/pkg/pipeline"
)

type defaultImplementation struct {
}

// Run gets a pipeline and executes it in a backend
func (di *defaultImplementation) Run(
	be *backend.BuildBackend, p *pipeline.Pipeline, opts *Options,
) (*backend.JobData, error) {
	return be.RunPipeline(p)
}
