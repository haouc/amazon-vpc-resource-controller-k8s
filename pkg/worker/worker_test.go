/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package worker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	resourceName                    = "vpc.amazonaws.com/pod-eni"
	workerCount                     = 1
	mockTimeToProcessWorkerFunc     = time.Duration(20)
	bufferTimeBwWorkerFuncExecution = time.Duration(3)
	maxRequeue                      = 3
)

func GetMockWorkerPool(ctx context.Context) Worker {
	log := zap.New(zap.UseDevMode(true)).WithValues("worker resource Id", resourceName)
	return NewDefaultWorkerPool(resourceName, workerCount, maxRequeue, log, ctx)
}

func MockWorkerFunc(job interface{}) (result ctrl.Result, err error) {
	v := job.(*int)
	*v++
	time.Sleep(time.Millisecond * mockTimeToProcessWorkerFunc)

	return ctrl.Result{}, nil
}

func TestNewWorkerPool(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w := GetMockWorkerPool(ctx)
	assert.NotNil(t, w)
}

// TestWorker_SubmitJob verifies that two different jobs are executed successfully.
func TestWorker_SubmitJob(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w := GetMockWorkerPool(ctx)
	err := w.StartWorkerPool(MockWorkerFunc)
	assert.NoError(t, err)

	// Count to verify job executed
	var jobCount = 2
	var job1 = 0
	var job2 = 0

	// Submit two jobs
	w.SubmitJob(&job1)
	w.SubmitJob(&job2)

	// Wait till the job complete. If the test is flaky, increase the buffer sleep time.
	time.Sleep(time.Millisecond * (mockTimeToProcessWorkerFunc + bufferTimeBwWorkerFuncExecution) * time.Duration(jobCount))

	// Verify job completed.
	assert.Equal(t, job1, 1)
	assert.Equal(t, job2, 1)
}

func TestWorker_SubmitJob_RequeueOnError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workerFunc := func(job interface{}) (result ctrl.Result, err error) {
		invoked := job.(*int)
		*invoked++

		return ctrl.Result{}, fmt.Errorf("error")
	}

	w := GetMockWorkerPool(ctx)
	err := w.StartWorkerPool(workerFunc)
	assert.NoError(t, err)

	var invoked = 0
	w.SubmitJob(&invoked)

	time.Sleep((mockTimeToProcessWorkerFunc + bufferTimeBwWorkerFuncExecution) * time.Millisecond * time.Duration(maxRequeue))

	// expected invocation = max requeue + the first invocation
	assert.Equal(t, maxRequeue+1, invoked)
}
