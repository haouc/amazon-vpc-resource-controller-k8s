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

package k8s

import (
	"context"
	"time"

	"github.com/aws/amazon-vpc-resource-controller-k8s/apis/vpcresources/v1alpha1"
	"github.com/aws/amazon-vpc-resource-controller-k8s/pkg/config"
	rcHealthz "github.com/aws/amazon-vpc-resource-controller-k8s/pkg/healthz"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

type CNINodeCleaner struct {
	Log            logr.Logger
	shutdown       bool
	ctx            context.Context
	k8sClient      client.Client
	leakedCNINodes []string
}

func (e *CNINodeCleaner) SetupWithManager(ctx context.Context, mgr ctrl.Manager, healthzHandler *rcHealthz.HealthzHandler) error {
	e.ctx = ctx
	e.k8sClient = mgr.GetClient()
	healthzHandler.AddControllersHealthCheckers(
		map[string]healthz.Checker{
			"health-cninode-cleaner": rcHealthz.SimplePing("cninode cleanup", e.Log),
		},
	)

	return mgr.Add(e)
}

// StartCNINodeCleaner starts the CNINOde Cleaner routine that cleans up dangling CNINode created by the controller
func (e *CNINodeCleaner) Start(ctx context.Context) error {
	e.Log.Info("starting CNINode clean up routine")
	// Start routine to listen for shut down signal, on receiving the signal it set shutdown to true
	go func() {
		<-ctx.Done()
		e.shutdown = true
	}()
	// Perform CNINode cleanup after fixed time intervals till shut down variable is set to true on receiving the shutdown
	// signal
	for !e.shutdown {
		e.cleanUpLeakedCNINodes()
		time.Sleep(config.CNINodeCleanUpInterval)
	}

	return nil
}

// cleanUpLeakedCNINodes describes all the CNINode objects that have no longer had k8s node objects referred to.
// this should be a very rare case and race condition caused CNINodes being left over when nodes were deleted.
func (e *CNINodeCleaner) cleanUpLeakedCNINodes() {
	e.Log.Info("CNINode cleaner doing routine cleanup...")
	nodes := &v1.NodeList{}
	if err := e.k8sClient.List(e.ctx, nodes); err != nil {
		e.Log.Error(err, "CNINode cleaner listing k8s nodes failed")
	}
	nodeMap := make(map[string]bool)
	for _, node := range nodes.Items {
		nodeMap[node.Name] = true
	}

	// check if the CNINodes added 12 hours ago are still not having k8s node referred to
	for _, nodeName := range e.leakedCNINodes {
		if _, found := nodeMap[nodeName]; !found {
			cniNode := &v1alpha1.CNINode{}
			key := types.NamespacedName{
				Name:      nodeName,
				Namespace: config.KubeDefaultNamespace,
			}

			var err error
			if err = e.k8sClient.Get(e.ctx, key, cniNode); err == nil {
				if err = e.k8sClient.Delete(e.ctx, cniNode); err == nil {
					e.Log.Info("CNINode cleanup deleted the leaked CNINode", "CNINode", cniNode)
				} else {
					e.Log.Error(err, "CNINode cleanup failed", "CNINode", cniNode)
				}
			} else {
				e.Log.Info("CNINode cleaner couldn't find leaked CNINode", "CNINodeName", nodeName)
			}
		}
	}

	e.leakedCNINodes = []string{}

	CNINodes := &v1alpha1.CNINodeList{}
	if err := e.k8sClient.List(e.ctx, CNINodes); err != nil {
		e.Log.Error(err, "CNINode cleaner listing CNINodes failed")
	}

	for _, cniNode := range CNINodes.Items {
		if _, found := nodeMap[cniNode.Name]; !found {
			// add "leaked" CNINodes to candidate list for next cleanup in 12 hours
			e.leakedCNINodes = append(e.leakedCNINodes, cniNode.Name)
			e.Log.Info("CNINode cleaner found possible leaked CNINode and added it to waiting list", "CNINode", cniNode)
		}
	}
}
