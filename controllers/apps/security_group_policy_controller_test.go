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
	"testing"

	"github.com/aws/amazon-vpc-resource-controller-k8s/apis/vpcresources/v1beta1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestSGPController_SetFlag(t *testing.T) {
	s1 := &v1beta1.SecurityGroupPolicy{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-1",
			Namespace: "test-namespace-1",
		},
	}
	s2 := &v1beta1.SecurityGroupPolicy{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-2",
			Namespace: "test-namespace-2",
		},
	}
	testScheme := runtime.NewScheme()
	v1beta1.AddToScheme(testScheme)
	controller := &SGPReconciler{
		Client: fake.NewClientBuilder().WithScheme(testScheme).WithObjects(s1, s2).Build(),
	}
	_, err := controller.Reconcile(context.TODO(), reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      "test-3",
		Namespace: "test-namespace-3",
	}})
	assert.NoError(t, err, "reconcile should succeed")
	assert.True(t, controller.GetSGPEnabledFlag())

	controller.Client.DeleteAllOf(context.TODO(), &v1beta1.SecurityGroupPolicy{}, &client.DeleteAllOfOptions{
		ListOptions: client.ListOptions{Namespace: v1.NamespaceAll},
	})
	list := &v1beta1.SecurityGroupPolicyList{}
	err = controller.Client.List(context.TODO(), list, &client.ListOptions{Namespace: v1.NamespaceAll})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(list.Items))
	assert.True(t, controller.GetSGPEnabledFlag())
}
