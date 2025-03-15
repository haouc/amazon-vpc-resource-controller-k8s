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

package core

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/amazon-vpc-resource-controller-k8s/controllers/apps"
	"github.com/aws/amazon-vpc-resource-controller-k8s/pkg/condition"
	"github.com/aws/amazon-vpc-resource-controller-k8s/pkg/config"
	rcHealthz "github.com/aws/amazon-vpc-resource-controller-k8s/pkg/healthz"

	"github.com/go-logr/logr"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/validate-v1-pod,mutating=false,matchPolicy=Equivalent,failurePolicy=ignore,groups="",resources=pods,verbs=create;update,versions=v1,name=vpod.vpc.k8s.aws,sideEffects=None,admissionReviewVersions=v1

// AnnotationValidator validates the resource allocated to the Pod via annotations. The WebHook
// prevents unauthorized user from modifying/removing these Annotations.
type AnnotationValidator struct {
	decoder       admission.Decoder
	Condition     condition.Conditions
	Log           logr.Logger
	Checker       healthz.Checker
	sgpController *apps.SGPReconciler
}

func NewAnnotationValidator(condition condition.Conditions, log logr.Logger, d admission.Decoder, healthzHandler *rcHealthz.HealthzHandler, sgpController *apps.SGPReconciler) *AnnotationValidator {
	annotationValidator := &AnnotationValidator{
		Condition:     condition,
		Log:           log,
		decoder:       d,
		sgpController: sgpController,
	}

	// add health check on subpath for pod annotation validating webhook
	healthzHandler.AddControllersHealthCheckers(
		map[string]healthz.Checker{
			"health-annotation-validating-webhook": rcHealthz.SimplePing("pod annotation validating webhook", log),
		},
	)

	return annotationValidator
}

// We are allowing multiple usernames to annotate the Windows/SGP Pod, eventually we will
// only allow user based authentication and optionally a service account based authentication
// for users wanting to run the controller for Windows IPAM on AWS Kubernetes.
const validUserInfo = "system:serviceaccount:kube-system:vpc-resource-controller"
const newValidUserInfo = "system:serviceaccount:kube-system:eks-vpc-resource-controller"
const vpcControllerUserName = "eks:vpc-resource-controller"

// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch

func (a *AnnotationValidator) Handle(_ context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}
	if req.Operation == admissionv1.Create || req.Operation == admissionv1.Update {
		if err := a.decoder.DecodeRaw(req.Object, pod); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		if WhichPod(pod) == Linux && !a.sgpController.GetSGPEnabledFlag() {
			return admission.Allowed("Security groups for pod is not enabled for linux pods")
		}
	}

	var response admission.Response

	switch req.Operation {
	case admissionv1.Create:
		response = a.handleCreate(pod)
	case admissionv1.Update:
		oldPod := &corev1.Pod{}
		if err := a.decoder.DecodeRaw(req.OldObject, oldPod); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
		response = a.handleUpdate(req.UserInfo.Username, pod, oldPod)
	default:
		response = admission.Allowed("")
	}

	a.Log.V(1).Info("annotation validating webhook response",
		"response", response)

	return response
}

func (a *AnnotationValidator) handleCreate(pod *corev1.Pod) admission.Response {
	// The annotation is added by vpc-resource-controller which will come as an update event
	// so we should block all request on create event
	for _, annotationKey := range a.getAnnotationKeysToBeValidated() {
		if val, ok := pod.Annotations[annotationKey]; ok {
			a.Log.Info("blocking request", "event", "create",
				"annotation key", annotationKey, "annotation value", val)
			return admission.Denied(
				fmt.Sprintf("pod cannot be created with %s annotation", annotationKey))
		}
	}
	return admission.Allowed("")
}

func (a *AnnotationValidator) handleUpdate(userName string, pod, oldPod *corev1.Pod) admission.Response {
	logger := a.Log.WithValues("name", pod.Name, "namespace", pod.Namespace, "uid", pod.UID)

	// Block any update on Fargate SGP Annotation Key. The Fargate Security Group Annotation is
	// added by the mutating WebHook on Create Event.
	if pod.Annotations[FargatePodSGAnnotationKey] !=
		oldPod.Annotations[FargatePodSGAnnotationKey] {
		logger.Info("denying annotation", "username", userName,
			"annotation key", FargatePodSGAnnotationKey)
		return admission.Denied("annotation is not set by mutating webhook")
	}

	// This will block any update on the specific annotation from non vpc resource controller
	// service accounts
	for _, annotationKey := range a.getAnnotationKeysToBeValidated() {
		if pod.Annotations[annotationKey] != oldPod.Annotations[annotationKey] {
			// Checking for two users, as the Service Account used by controller was changed
			// after first release.
			if (userName != validUserInfo) && (userName != newValidUserInfo) &&
				(userName != vpcControllerUserName) {
				logger.Info("denying annotation", "username", userName,
					"annotation key", annotationKey)
				return admission.Denied("annotation is not set by vpc-resource-controller")
			}
		}
	}
	return admission.Allowed("")
}

// getAnnotationKeysToBeValidated returns the list of
func (a *AnnotationValidator) getAnnotationKeysToBeValidated() []string {
	// Pod ENI annotation is validated by default
	annotationsToValidate := []string{config.ResourceNamePodENI}
	if a.Condition.IsWindowsIPAMEnabled() {
		// Windows IPv4 Annotation is validated if feature is enabled, as the older controller could
		// be installed on Customer Data Plane and new controller should not block it's annotations
		annotationsToValidate = append(annotationsToValidate, config.ResourceNameIPAddress)
	}
	return annotationsToValidate
}
