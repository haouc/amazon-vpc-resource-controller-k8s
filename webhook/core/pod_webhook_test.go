package core

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	vpcresourcesv1beta1 "github.com/aws/amazon-vpc-resource-controller-k8s/apis/vpcresources/v1beta1"
	vpcresourceconfig "github.com/aws/amazon-vpc-resource-controller-k8s/pkg/config"
	webhookutils "github.com/aws/amazon-vpc-resource-controller-k8s/pkg/utils"
)

var (
	name      = "test"
	namespace = "test_namespace"
	saName    = "test_sa"
	logger    = ctrl.Log.WithName("test")
	testPa    *PodAnnotator
	handlerPa *PodAnnotator
	ctx       = context.Background()
)

func init() {
	testPa = getPodAnnotator()
	handlerPa = getPodAnnotator()
}

func getPodAnnotator() *PodAnnotator {
	testScheme := runtime.NewScheme()
	clientgoscheme.AddToScheme(testScheme)
	vpcresourcesv1beta1.AddToScheme(testScheme)
	testClient := fake.NewFakeClientWithScheme(
		testScheme,
		NewPod(name, saName, namespace),
		NewServiceAccount(saName, namespace),
		NewSecurityGroupPolicy(name, namespace, []string{"sg-00001"}),
	)
	decoder, _ := admission.NewDecoder(testScheme)
	pa := &PodAnnotator{
		Client:  testClient,
		decoder: decoder,
		CacheHelper: &webhookutils.K8sCacheHelper{
			Client: testClient,
			Log:    logger,
		},
		Log: logger,
	}
	return pa
}

// TestAttachPrivateIP tests if pod can be injected with private IP.
func TestAttachPrivateIP(t *testing.T) {
	pod := NewWindowsPod("test", "test_namespace", true)
	ok := shouldAttachPrivateIP(pod)
	assert.True(t, ok)

	pod = NewPod("test", "sa_test", "test_namespace")
	ok = shouldAttachPrivateIP(pod)
	assert.True(t, !ok)
}

// TestAttachPrivateIPByNodeSelector tests if pod is labeled as Windows by NodeSelector.
func TestAttachPrivateIPByNodeSelector(t *testing.T) {
	pod := NewWindowsPod("test", "test_namespace", true)
	ok := hasWindowsNodeSelector(pod)
	assert.True(t, ok)
}

// TestAttachPrivateIPByNodeSelector tests if pod is labeled as Windows by NodeAffinity.
func TestAttachPrivateIPByNodeAffinity(t *testing.T) {
	pod := NewWindowsPod("test", "test_namespace", false)
	ok := hasWindowsNodeAffinity(pod)
	// TODO: implement node affinity for windows pod to enable this test.
	assert.True(t, !ok)
}

// TestCheckContainerLimits tests if pod's container(s) has limits added by user.
func TestCheckContainerLimits(t *testing.T) {
	//pod := webhookutils.NewPodWithContainerLimits("test", "test_namespace", true)
	pod := NewPodWithContainerLimits("test", "test_namespace", true)

	// TODO: implement container user input in limits and/or requests
	hasLimits := containerHasCustomizedLimit(pod)
	assert.True(t, !hasLimits)
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
	pod := NewPod("test", "test_sa", "test_namespace")
	resp := getResponse(pod)
	assert.True(t, resp.Allowed)

	for _, p := range resp.Patches {
		assert.True(t, p.Operation == "add")
		assert.True(t, p.Path == "/spec/containers/0/resources/limits" ||
			p.Path == "/spec/containers/0/resources/requests")

		pv := p.Value.(map[string]interface{})
		assert.True(t, pv[vpcresourceconfig.ResourceNamePodENI] == resourceLimit)
	}
}

// TestPodAnnotator_Handle test webhook mutating requested Windows pod.
func TestPodAnnotator_Windows_Handle(t *testing.T) {
	pod := NewWindowsPod("test", "test_namespace", true)
	resp := getResponse(pod)
	assert.True(t, resp.Allowed)

	for _, p := range resp.Patches {
		assert.True(t, p.Operation == "add")
		assert.True(t, p.Path == "/spec/containers/0/resources/limits" ||
			p.Path == "/spec/containers/0/resources/requests")

		pv := p.Value.(map[string]interface{})
		assert.True(t, pv[vpcresourceconfig.ResourceNameIPAddress] == resourceLimit)
	}
}

func getResponse(pod *corev1.Pod) admission.Response {
	podRaw, _ := json.Marshal(pod)
	req := admission.Request{
		AdmissionRequest: v1beta1.AdmissionRequest{
			Operation: v1beta1.Create,
			Object: runtime.RawExtension{
				Raw: podRaw,
			},
		},
	}
	resp := handlerPa.Handle(ctx, req)
	return resp
}

func NewPod(name string, sa string, namespace string) *corev1.Pod {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels: map[string]string{
				"role":        "db",
				"environment": "qa",
				"app":         "test_app",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Resources: corev1.ResourceRequirements{},
				},
			},
			ServiceAccountName: sa,
		},
		Status: corev1.PodStatus{
			Conditions: []corev1.PodCondition{
				{
					Type:   corev1.PodReady,
					Status: corev1.ConditionTrue,
				},
			},
		},
	}
	return pod
}

func NewWindowsPod(name string, namespace string, useSelector bool) *corev1.Pod {
	var spec corev1.PodSpec
	containers := []corev1.Container{
		{
			Resources: corev1.ResourceRequirements{},
		},
	}

	if useSelector {
		spec = corev1.PodSpec{
			Containers: containers,
			NodeSelector: map[string]string{
				"kubernetes.io/os": "windows",
			},
		}
	} else {
		spec = corev1.PodSpec{
			Containers: containers,
			Affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "kubernetes.io/os",
										Operator: "In",
										Values:   []string{"windows"},
									},
								},
								MatchFields: nil,
							},
						},
					},
				},
			},
		}
	}

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels: map[string]string{
				"role":        "db",
				"environment": "qa",
				"app":         "test_app",
			},
		},
		Spec: spec,
		Status: corev1.PodStatus{
			Conditions: []corev1.PodCondition{
				{
					Type:   corev1.PodReady,
					Status: corev1.ConditionTrue,
				},
			},
		},
	}
	return pod
}

func NewPodWithContainerLimits(name string, namespace string, hasLimit bool) *corev1.Pod {
	pod := NewPod(name, "", namespace)
	limit := corev1.ResourceList{}
	if hasLimit {
		limit["key"] = resource.MustParse("1")
	}
	pod.Spec.Containers = []corev1.Container{
		{
			Name: "test_container_1",
			Resources: corev1.ResourceRequirements{
				Limits: nil,
			},
		},
		{
			Name: "test_container_1",
			Resources: corev1.ResourceRequirements{
				Limits: limit,
			},
		},
	}
	return pod
}

func NewServiceAccount(name string, namespace string) *corev1.ServiceAccount {
	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"role":        "db",
				"environment": "qa",
			},
		},
	}
	return sa
}

func NewSecurityGroupPolicy(name string, namespace string, securityGroups []string) *vpcresourcesv1beta1.SecurityGroupPolicy {
	sgp := &vpcresourcesv1beta1.SecurityGroupPolicy{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: vpcresourcesv1beta1.SecurityGroupPolicySpec{
			PodSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"role": "db"},
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "environment",
						Operator: "In",
						Values:   []string{"qa", "production"},
					},
				},
			},
			SecurityGroups: vpcresourcesv1beta1.GroupIds{
				Groups: securityGroups,
			},
		},
	}
	return sgp
}
