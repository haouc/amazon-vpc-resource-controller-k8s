package core

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	webhookutils "github.com/aws/amazon-vpc-resource-controller-k8s/pkg/utils"
)

var (
	logger = ctrl.Log.WithName("test")

	testPa = PodAnnotator{
		Client:  webhookutils.TestClient,
		decoder: nil,
		CacheHelper: &webhookutils.K8sCacheHelper{
			Client: webhookutils.TestClient,
			Log:    logger,
		},
		Log: logger,
	}
	ctx = context.Background()
)

// TestAttachPrivateIp tests if pod can be injected with private IP.
func TestAttachPrivateIp(t *testing.T) {
	pod := webhookutils.NewWindowsPod("test", "test_namespace", true)
	ok := shouldAttachPrivateIP(pod)
	assert.True(t, ok)

	pod = webhookutils.NewPod("test", "sa_test", "test_namespace")
	ok = shouldAttachPrivateIP(pod)
	assert.True(t, !ok)
}

// TestAttachPrivateIpByNodeSelector tests if pod is labeled as Windows by NodeSelector.
func TestAttachPrivateIpByNodeSelector(t *testing.T) {
	pod := webhookutils.NewWindowsPod("test", "test_namespace", true)
	ok := hasWindowsNodeSelector(pod)
	assert.True(t, ok)
}

// TestAttachPrivateIpByNodeSelector tests if pod is labeled as Windows by NodeAffinity.
func TestAttachPrivateIpByNodeAffinity(t *testing.T) {
	pod := webhookutils.NewWindowsPod("test", "test_namespace", false)
	ok := hasWindowsNodeAffinity(pod)
	// TODO: implement node affinity for windows pod to enable this test.
	assert.True(t, !ok)
}

// TestCheckContainerLimits tests if pod's container(s) has limits added by user.
func TestCheckContainerLimits(t *testing.T) {
	//pod := webhookutils.NewPodWithContainerLimits("test", "test_namespace", true)
	pod := webhookutils.NewPodWithContainerLimits("test", "test_namespace", true)

	// TODO: implement container user input in limits and/or requests
	hasLimits := containerHasCustomizedLimit(pod)
	assert.True(t, !hasLimits)
}

// TestPodAnnotator_InjectClient tests injecting client to pod annotator.
func TestPodAnnotator_InjectClient(t *testing.T) {
	testPa.InjectClient(webhookutils.TestClient)
	assert.NotNil(t, testPa.Client)
	pod := &corev1.Pod{}
	assert.NoError(t,
		testPa.Client.Get(
			context.Background(),
			types.NamespacedName{Name: "test", Namespace: "test_namespace"},
			pod))
	assert.True(t, pod.Name == "test")
}

// TestPodAnnotator_InjectDecoder tests injecting decoder into pod annotator.
func TestPodAnnotator_InjectDecoder(t *testing.T) {
	var decoder *admission.Decoder
	assert.NoError(t, testPa.InjectDecoder(decoder))
}

// TestPodAnnotator_Handle test webhook mutating requested empty request.
func TestPodAnnotator_Empty_Handle(t *testing.T) {
	resp := testPa.Handle(ctx, admission.Request{})
	assert.True(t, !resp.Allowed && resp.Result.Code == http.StatusBadRequest)
}

// TestPodAnnotator_Handle test webhook mutating requested Linux pod.
func TestPodAnnotator_Handle(t *testing.T) {
	pod := webhookutils.NewPod("test", "test_sa", "test_namespace")
	resp := getResponse(pod)
	assert.True(t, resp.Allowed)

	for _, p := range resp.Patches {
		assert.True(t, p.Operation == "add")
		assert.True(t, p.Path == "/spec/containers/0/resources/limits" ||
			p.Path == "/spec/containers/0/resources/requests")

		pv := p.Value.(map[string]interface{})
		assert.True(t, pv[podEniRequest] == resourceLimit)
	}
}

// TestPodAnnotator_Handle test webhook mutating requested Windows pod.
func TestPodAnnotator_Windows_Handle(t *testing.T) {
	pod := webhookutils.NewWindowsPod("test", "test_namespace", true)
	resp := getResponse(pod)
	assert.True(t, resp.Allowed)

	for _, p := range resp.Patches {
		assert.True(t, p.Operation == "add")
		assert.True(t, p.Path == "/spec/containers/0/resources/limits" ||
			p.Path == "/spec/containers/0/resources/requests")

		pv := p.Value.(map[string]interface{})
		assert.True(t, pv[ResourceNameIPAddress] == resourceLimit)
	}
}

func getResponse(pod *corev1.Pod) admission.Response {
	decoder, _ := admission.NewDecoder(webhookutils.TestScheme)
	testPa.decoder = decoder
	podRaw, _ := json.Marshal(pod)
	req := admission.Request{
		AdmissionRequest: v1beta1.AdmissionRequest{
			Operation: v1beta1.Create,
			Object: runtime.RawExtension{
				Raw: podRaw,
			},
		},
	}
	resp := testPa.Handle(ctx, req)
	return resp
}
