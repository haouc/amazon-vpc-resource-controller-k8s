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

package k8s

import (
	"context"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	prometheusRegistered = false

	annotatePodRequestCallCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "annotate_pod_request_call_count",
			Help: "The number of request to annotate pod object",
		},
		[]string{"annotate_key"},
	)

	annotatePodRequestErrCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "annotate_pod_request_err_count",
			Help: "The number of request that failed to annotate the pod",
		},
		[]string{"annotate_key"},
	)

	advertiseResourceRequestCallCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "advertise_resource_request_call_count",
			Help: "The number of request to advertise extended resource",
		},
		[]string{"resource_name"},
	)

	advertiseResourceRequestErrCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "advertise_resource_request_err_count",
			Help: "The number of request that failed to advertise extended resource",
		},
		[]string{"resource_name"},
	)
)

func prometheusRegister() {
	prometheus.MustRegister(annotatePodRequestErrCount)
	prometheus.MustRegister(annotatePodRequestCallCount)
	prometheus.MustRegister(advertiseResourceRequestErrCount)
	prometheus.MustRegister(advertiseResourceRequestCallCount)

	prometheusRegistered = true
}

// K8sWrapper represents an interface with all the common operations on K8s objects
type K8sWrapper interface {
	GetPod(namespace string, name string) (*v1.Pod, error)
	AnnotatePod(podNamespace string, podName string, key string, val string) error
	AdvertiseCapacityIfNotSet(nodeName string, resourceName string, capacity int) error
}

// k8sWrapper is the wrapper object with the client
type k8sWrapper struct {
	client client.Client
}

// NewK8sWrapper returns a new K8sWrapper
func NewK8sWrapper(client client.Client) K8sWrapper {
	if !prometheusRegistered {
		prometheusRegister()
	}
	return &k8sWrapper{client: client}
}

// AnnotatePod annotates the pod with the provided key and value
func (k *k8sWrapper) AnnotatePod(podNamespace string, podName string, key string, val string) error {
	annotatePodRequestCallCount.WithLabelValues(key).Inc()
	ctx := context.Background()

	request := types.NamespacedName{
		Namespace: podNamespace,
		Name:      podName,
	}

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		// Get the latest copy of the pod from cache
		pod := &v1.Pod{}
		if err := k.client.Get(ctx, request, pod); err != nil {
			return err
		}
		newPod := pod.DeepCopy()
		newPod.Annotations[key] = val

		return k.client.Patch(ctx, newPod, client.MergeFrom(pod))
	})

	if err != nil {
		annotatePodRequestErrCount.WithLabelValues(key).Inc()
		return err
	}

	return nil
}

// GetPod returns the pod object using the client cache
func (k *k8sWrapper) GetPod(namespace string, name string) (*v1.Pod, error) {
	ctx := context.Background()

	pod := &v1.Pod{}
	if err := k.client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, pod); err != nil {
		return nil, err
	}
	return pod, nil
}

// AdvertiseCapacity advertises the resource capacity for the given resource
func (k *k8sWrapper) AdvertiseCapacityIfNotSet(nodeName string, resourceName string, capacity int) error {
	ctx := context.Background()

	request := types.NamespacedName{
		Name: nodeName,
	}

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		node := &v1.Node{}
		if err := k.client.Get(ctx, request, node); err != nil {
			return err
		}

		existingCapacity := node.Status.Capacity[v1.ResourceName(resourceName)]
		if !existingCapacity.IsZero() && existingCapacity.Value() == int64(capacity) {
			return nil
		}

		// Capacity doesn't match the expected capacity, need to advertise again
		advertiseResourceRequestCallCount.WithLabelValues(resourceName).Inc()

		newNode := node.DeepCopy()
		newNode.Status.Capacity[v1.ResourceName(resourceName)] = resource.MustParse(strconv.Itoa(capacity))

		return k.client.Status().Patch(ctx, newNode, client.MergeFrom(node))
	})

	if err != nil {
		advertiseResourceRequestErrCount.WithLabelValues(resourceName).Inc()
		return err
	}

	return nil
}
