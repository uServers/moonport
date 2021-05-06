// Copyright 2021 - 2021, Adolfo García Veytia and the moonport contributors
// SPDX-License-Identifier: Apache-2.0

package backend

import (
	"context"
	"fmt"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uServers/moonport/pkg/pipeline"
	cloudbuildpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

func NewDriverCloudBuild() (*DriverCloudBuild, error) {
	// Create the new driver
	driver := DriverCloudBuild{}
	cb, err := cloudbuild.NewClient(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "getting new Cloud Build client")
	}

	driver.impl = &defBackendCloudBuildImpl{}
	driver.client = cb
	return &driver, nil
}

type DriverCloudBuild struct {
	impl   DriverCloudBuildImplementation
	client *cloudbuild.Client
	// opts *opts
}

type DriverCloudBuildImplementation interface {
	CreateBuild(context.Context, *pipeline.Pipeline, *cloudbuild.Client) (*JobData, error)
}

// var client
func (d *DriverCloudBuild) CreatePipeline(
	ctx context.Context, p *pipeline.Pipeline) (data *JobData, err error,
) {
	return d.impl.CreateBuild(context.Background(), p, d.client)
}

type defBackendCloudBuildImpl struct{}

func (impl *defBackendCloudBuildImpl) CreateBuild(
	ctx context.Context, p *pipeline.Pipeline, client *cloudbuild.Client,
) (data *JobData, err error) {
	// Verify the pipeline before launching
	if err := p.Validate(); err != nil {
		return nil, errors.Wrap(err, "while validating pipeline before launch")
	}
	// Launch the build
	logrus.Info("Building Cloud Build request")

	// Translate the env vars to a string slice
	env := []string{}
	for n, val := range p.Spec.EnvVars {
		env = append(env, fmt.Sprintf("%s=%s", n, val))
	}
	steps := []*cloudbuildpb.BuildStep{}
	// FIXME: Cicla y dale
	for name, stage := range p.Spec.Stages {
		logrus.Info("Building stage " + name)
		for _, stepspec := range stage.Steps {
			step := p.GetStep(stepspec.StepLabel)
			if step == nil {
				return nil, errors.New("unable to find step " + stepspec.StepLabel)
			}
			steps = append(steps, &cloudbuildpb.BuildStep{
				Name:       step.Spec.Image,
				Env:        env,
				Entrypoint: "bash",
				Args:       []string{"-c", "echo hello world"},
				/*
					Dir:        "",
					Id:         "",
					WaitFor:    []string{},

					SecretEnv:  []string{},
					Volumes:    []*cloudbuildpb.Volume{},
					Timing:     &cloudbuildpb.TimeSpan{},
					PullTiming: &cloudbuildpb.TimeSpan{},
					Timeout:    &durationpb.Duration{},
					Status:     0,
				*/
			})
		}
	}

	req := &cloudbuildpb.CreateBuildRequest{
		//Parent:    "",
		ProjectId: "ulabs-cloud-tests",
		Build: &cloudbuildpb.Build{
			Name:  p.Metadata.Name,
			Steps: steps,
			/*
				Options:          &cloudbuildpb.BuildOptions{
					SourceProvenanceHash:  []cloudbuildpb.Hash_HashType{},
					RequestedVerifyOption: 0,
					MachineType:           0,
					DiskSizeGb:            0,
					SubstitutionOption:    0,
					DynamicSubstitutions:  false,
					LogStreamingOption:    0,
					WorkerPool:            "",
					Logging:               0,
					Env:                   []string{},
					SecretEnv:             []string{},
					Volumes:               []*cloudbuildpb.Volume{},
				},
			*/
			/*
				Name:             "test",
				Id:               "",
				ProjectId:        "",
				Status:           0,
				StatusDetail:     "",
				Source:           &cloudbuildpb.Source{},
				Steps:            []*cloudbuildpb.BuildStep{},
				Results:          &cloudbuildpb.Results{},
				CreateTime:       &timestamppb.Timestamp{},
				StartTime:        &timestamppb.Timestamp{},
				FinishTime:       &timestamppb.Timestamp{},
				Timeout:          &durationpb.Duration{},
				Images:           []string{},
				QueueTtl:         &durationpb.Duration{},
				Artifacts:        &cloudbuildpb.Artifacts{},
				LogsBucket:       "",
				SourceProvenance: &cloudbuildpb.SourceProvenance{},
				BuildTriggerId:   "",

				LogUrl:           "",
				Substitutions:    map[string]string{},
				Tags:             []string{},
				Secrets:          []*cloudbuildpb.Secret{},
				Timing:           map[string]*cloudbuildpb.TimeSpan{},
				ServiceAccount:   "",
				AvailableSecrets: &cloudbuildpb.Secrets{},
			*/
		},
	}

	op, err := client.CreateBuild(ctx, req)
	if err != nil {
		return data, errors.Wrap(err, "creating Cloud Build operation")
	}
	md, err := op.Metadata()
	if err != nil {
		return nil, errors.Wrap(err, "getting build job metadata")
	}

	logrus.Infof("Successfully launched build %s, ", md.Build.Id)

	// Wait for the job to be created
	/*
		logrus.Info("Launching starting cloud build job")
		resp, err := op.Wait(ctx)
		if err != nil {
			return data, errors.Wrap(err, "creating GCB job")
		}
		// TODO: Use resp.
		_ = resp
	*/

	return &JobData{
		JobID: md.Build.Id,
	}, nil
}
