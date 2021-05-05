// Copyright 2021 - 2021, Adolfo García Veytia and the moonport contributors
// SPDX-License-Identifier: Apache-2.0

package backend

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uServers/moonport/pkg/pipeline"
)

// NewBuildBackend arma un backend nuevo
func NewBuildBackend(moniker string) (bg *BuildBackend, err error) {
	bg = &BuildBackend{}
	var d BuildBackendDriver

	// Choose the backend
	switch moniker {
	case "gcb":
		d, err = NewDriverCloudBuild()
	default:
		return nil, errors.New("unknown build driver")
	}

	if err != nil {
		return bg, errors.Wrapf(err, "creating %s backend driver", moniker)
	}
	bg.driver = d
	return bg, nil
}

type JobData struct {
	JobID string
}

type BuildBackend struct {
	driver BuildBackendDriver
}

func (b *BuildBackend) RunPipeline(p *pipeline.Pipeline) (*JobData, error) {
	logrus.Info("Launching pipeline using backend")
	return b.driver.CreatePipeline(context.Background(), p)
}

type BuildBackendDriver interface {
	CreatePipeline(context.Context, *pipeline.Pipeline) (*JobData, error)
}
