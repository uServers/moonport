// Copyright 2021 - 2021, Adolfo García Veytia and the moonport contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/uServers/moonport/pkg/launchpad"
	"github.com/uServers/moonport/pkg/pipeline"
)

func main() {
	p, err := pipeline.NewPipelineFromFile(
		"/home/urbano/Projects/moonport/examples/pipelines/hello-world.yaml",
		&pipeline.Options{
			StepSources: []string{"/home/urbano/Projects/moonport/examples/steps/"},
		},
	)
	if err != nil {
		logrus.Fatal(err)
	}
	lp, err := launchpad.New()
	if err != nil {
		logrus.Fatal(err)
	}
	data, err := lp.Run(p)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Result: %+v", data)
}
