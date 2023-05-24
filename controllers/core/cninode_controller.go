package controllers

import (
	"context"

	"github.com/aws/amazon-vpc-resource-controller-k8s/apis/vpcresources/v1alpha1"
	rcHealthz "github.com/aws/amazon-vpc-resource-controller-k8s/pkg/healthz"
	"github.com/aws/amazon-vpc-resource-controller-k8s/pkg/k8s"
	"github.com/aws/amazon-vpc-resource-controller-k8s/pkg/node/manager"
	"github.com/go-logr/logr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type CNINodeController struct {
	Log         logr.Logger
	NodeManager manager.Manager
	Context     context.Context
	K8sAPI      k8s.K8sWrapper
}

func (s *CNINodeController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	cniNode, err := s.K8sAPI.GetCNINode(req.Name, req.Namespace)
	if err != nil {
		s.Log.Error(err, "Failed to get CNINode")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if len(cniNode.Spec.Features) == 0 {
		return ctrl.Result{}, nil
	}

	nodeName := cniNode.Name
	s.Log.Info("The CNINode has been updated", "CNINode", nodeName, "Features", cniNode.Spec.Features)

	if node, err := s.K8sAPI.GetNode(nodeName); err != nil {
		// if not found the node, we don't requeue and just wait VPC CNI to send another request
		s.Log.Info("CNINode reconciler didn't find the node and will do exponential retry", "Node", node.Name)
		return ctrl.Result{}, err
	} else {
		// make sure we have the node in cache already
		if _, found := s.NodeManager.GetNode(nodeName); !found {
			s.Log.Info("CNINode controller could not find node, will try again", "NodeName", nodeName)
			return ctrl.Result{Requeue: true}, nil
		} else {
			// add the node into working queue as UPDATE event
			s.Log.Info("CNINode controller found the node and sending the node update")
			if err = s.NodeManager.UpdateNode(nodeName); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

func (s *CNINodeController) SetupWithManager(mgr ctrl.Manager, healthzHandler *rcHealthz.HealthzHandler) error {
	healthzHandler.AddControllersHealthCheckers(
		map[string]healthz.Checker{
			"health-cninode-controller": s.Check(),
		},
	)
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.CNINode{}).WithEventFilter(
		eventFilters(),
	).Complete(s)
}

func eventFilters() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
	}
}

func (s *CNINodeController) Check() healthz.Checker {
	s.Log.Info("CNINode controller's healthz subpath was added")
	// We can revisit this to use PingWithTimeout() instead if we have concerns on this controller.
	return rcHealthz.SimplePing("configmap controller", s.Log)
}
