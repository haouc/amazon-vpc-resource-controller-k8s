// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package manifest

import (
	"github.com/aws/aws-sdk-go/aws"
	batchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobBuilder struct {
	namespace              string
	name                   string
	parallelism            int
	os                     string
	container              v1.Container
	labels                 map[string]string
	terminationGracePeriod int
}

func NewWindowsJob() *JobBuilder {
	return &JobBuilder{
		namespace:              "windows-test",
		name:                   "windows-job",
		os:                     "windows",
		terminationGracePeriod: 0,
		labels:                 map[string]string{},
	}
}

func (j *JobBuilder) Name(name string) *JobBuilder {
	j.name = name
	return j
}

func (j *JobBuilder) Namespace(namespace string) *JobBuilder {
	j.namespace = namespace
	return j
}

func (j *JobBuilder) OS(os string) *JobBuilder {
	j.os = os
	return j
}

func (j *JobBuilder) Container(container v1.Container) *JobBuilder {
	j.container = container
	return j
}

func (j *JobBuilder) PodLabels(labelKey string, labelVal string) *JobBuilder {
	j.labels[labelKey] = labelVal
	return j
}

func (j *JobBuilder) TerminationGracePeriod(terminationGracePeriod int) *JobBuilder {
	j.terminationGracePeriod = terminationGracePeriod
	return j
}

func (j *JobBuilder) Parallelism(parallelism int) *JobBuilder {
	j.parallelism = parallelism
	return j
}

func (j *JobBuilder) Build() *batchV1.Job {
	return &batchV1.Job{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      j.name,
			Namespace: j.namespace,
		},
		Spec: batchV1.JobSpec{
			Parallelism: aws.Int32(int32(j.parallelism)),
			Template: v1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: j.labels,
				},
				Spec: v1.PodSpec{
					Containers:                    []v1.Container{j.container},
					TerminationGracePeriodSeconds: aws.Int64(int64(j.terminationGracePeriod)),
					NodeSelector:                  map[string]string{"kubernetes.io/os": j.os},
					RestartPolicy:                 v1.RestartPolicyNever,
				},
			},
		},
	}
}
