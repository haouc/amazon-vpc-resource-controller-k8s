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

package apps

import (
	"context"
	"sync"

	"github.com/aws/amazon-vpc-resource-controller-k8s/apis/vpcresources/v1beta1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SGPReconciler struct {
	Log            logr.Logger
	Client         client.Client
	SGPEnabledFlag bool
	lock           sync.RWMutex
}

func (s *SGPReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// TODO: simplify this once getting event on SGP, set the flag to true.
	sgpList := &v1beta1.SecurityGroupPolicyList{}
	if err := s.Client.List(ctx, sgpList, &client.ListOptions{Namespace: v1.NamespaceAll}); err != nil {
		return ctrl.Result{}, err
	}
	if !s.GetSGPEnabledFlag() && len(sgpList.Items) > 0 {
		s.Log.Info("Set SGP flag to true")
		s.setSGPEnabledFlag(true)
	}
	return ctrl.Result{}, nil
}

func (s *SGPReconciler) setSGPEnabledFlag(enabled bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.SGPEnabledFlag = enabled
}

func (s *SGPReconciler) GetSGPEnabledFlag() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.SGPEnabledFlag
}

func (s *SGPReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.SecurityGroupPolicy{}).
		Complete(s)
}
