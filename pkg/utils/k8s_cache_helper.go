package utils

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

/* This is a shared local cache for storing security group or other objects which can be shared access by controllers
or webhooks.
*/

type K8sCacheHelper struct {
	Client client.Client
	Log    logr.Logger
}

// NewK8sCacheHelper creates and returns a controller-runtime cache operator.
func NewK8sCacheHelper(client client.Client, log logr.Logger) *K8sCacheHelper {
	return &K8sCacheHelper{
		Client: client,
		Log:    log,
	}
}

// GetSecurityGroupsFromPod returns security groups assigned to a testPod based on it's NamespacedName.
func (kch *K8sCacheHelper) GetSecurityGroupsFromPod(podId types.NamespacedName) []string {
	pod := &corev1.Pod{}

	if err := kch.Client.Get(context.Background(), podId, pod); err != nil {
		panic(err)
	} else {
		return kch.ShouldAddENILimits(pod)
	}
}
