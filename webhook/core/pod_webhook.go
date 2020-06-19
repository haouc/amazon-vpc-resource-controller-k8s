package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	webhookutils "github.com/aws/amazon-vpc-resource-controller-k8s/pkg/utils"
)

const (
	resourceLimit         = "1"
	podEniRequest         = "vpc.amazonaws.com/pod-eni"
	ResourceNameIPAddress = "vpc.amazonaws.com/PrivateIPv4Address"

	// NodeLabelOS is the Kubernetes OS label.
	NodeLabelOS        = "kubernetes.io/os"
	NodeLabelOSBeta    = "beta.kubernetes.io/os"
	NodeLabelOSWindows = "windows"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create,versions=v1,name=mpod.vpc.k8s.aws

// PodAnnotator annotates Pods
type PodAnnotator struct {
	Client      client.Client
	decoder     *admission.Decoder
	CacheHelper *webhookutils.K8sCacheHelper
	Log         logr.Logger
}

// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=serviceaccounts/status,verbs=get

func (a *PodAnnotator) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}
	webhookLog := a.Log.WithValues("Pod name", pod.Name, "Pod namespace", pod.Namespace)
	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// Ignore pod that is scheduled on host network.
	if pod.Spec.HostNetwork {
		webhookLog.Info("Not injecting resources as the pod runs on Host Network.")
		return admission.Allowed("Pod on HostNetwork will not be injected with resources.")
	}

	webhookLog.Info("Requesting Mutating Pod: ",
		"OS", pod.Spec.NodeSelector,
		"Resources Limits", pod.Spec.Containers[0].Resources.Limits)

	if pod.Spec.Containers[0].Resources.Limits == nil {
		pod.Spec.Containers[0].Resources.Limits = make(corev1.ResourceList)
	}

	if pod.Spec.Containers[0].Resources.Requests == nil {
		pod.Spec.Containers[0].Resources.Requests = make(corev1.ResourceList)
	}

	// Attach private ip to Windows pod which is not running on Host Network.
	// Attach eni to non-Windows pod which is not running on Host Network.
	if shouldAttachPrivateIP(pod) {
		webhookLog.Info("The pod is valid to be added with private ipv4 address.")
		pod.Spec.Containers[0].Resources.Limits[ResourceNameIPAddress] = resource.MustParse(resourceLimit)
		pod.Spec.Containers[0].Resources.Requests[ResourceNameIPAddress] = resource.MustParse(resourceLimit)
	} else if sgList := a.CacheHelper.ShouldAddEniLimits(pod); len(sgList) > 0 {
		webhookLog.Info("The pod is valid to be added with eni resources.")
		pod.Spec.Containers[0].Resources.Limits[podEniRequest] = resource.MustParse(resourceLimit)
		pod.Spec.Containers[0].Resources.Requests[podEniRequest] = resource.MustParse(resourceLimit)
	} else {
		return admission.Allowed("Pod will not be injected with resources limits.")
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		webhookLog.Error(err, "Marshalling pod failed:")
		return admission.Errored(http.StatusInternalServerError, err)
	}
	webhookLog.Info("Mutating Pod finished. ",
		"Resources Limits", pod.Spec.Containers[0].Resources.Limits)

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

func shouldAttachPrivateIP(pod *corev1.Pod) bool {
	return hasWindowsNodeSelector(pod) || hasWindowsNodeAffinity(pod)
}

func hasWindowsNodeSelector(pod *corev1.Pod) bool {
	osLabel := pod.Spec.NodeSelector[NodeLabelOS]

	// Version Beta is going to be deprecated soon.
	osLabelBeta := pod.Spec.NodeSelector[NodeLabelOSBeta]

	if osLabel != NodeLabelOSWindows && osLabelBeta != NodeLabelOSWindows {
		return false
	}

	return true
}

func hasWindowsNodeAffinity(pod *corev1.Pod) bool {
	// TODO: implement node affinity for Windows pod
	// Referring to https://t.corp.amazon.com/V167778691
	return false
}

func containerHasCustomizedLimit(pod *corev1.Pod) bool {
	// TODO: implement container limits user input
	// Referring to https://sim.amazon.com/issues/EKS-NW-424
	return false
}

// PodAnnotator implements inject.Client.
// A client will be automatically injected.

// InjectClient injects the client.
func (a *PodAnnotator) InjectClient(c client.Client) error {
	a.Client = c
	return nil
}

// PodAnnotator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (a *PodAnnotator) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
