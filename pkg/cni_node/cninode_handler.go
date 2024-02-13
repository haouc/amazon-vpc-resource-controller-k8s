package cni_node

import (
	"github.com/aws/amazon-vpc-resource-controller-k8s/apis/vpcresources/v1alpha1"
	"github.com/aws/amazon-vpc-resource-controller-k8s/pkg/aws/vpc"
	"github.com/aws/amazon-vpc-resource-controller-k8s/pkg/k8s"
	"github.com/go-logr/logr"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/types"
)

type CNINodeHandler struct {
	client k8s.K8sWrapper
	logger logr.Logger
}

func New(client k8s.K8sWrapper, logger logr.Logger) *CNINodeHandler {
	return &CNINodeHandler{client: client, logger: logger}
}

func (cnh *CNINodeHandler) ListAllWarmedBranchENIs(nodeName string) []v1alpha1.WarmBranchENI {
	node, err := cnh.client.GetNode(nodeName)
	if err != nil {
		return nil
	}
	if cnd, err := cnh.client.GetCNINode(types.NamespacedName{Name: node.Name}); err != nil {
		return nil
	} else {
		return cnd.Status.BranchENIs
	}
}

func (cnh *CNINodeHandler) AddBranchENIToWarmPool(nodeName string, enis []v1alpha1.WarmBranchENI) error {
	node, err := cnh.client.GetNode(nodeName)
	if err != nil {
		cnh.logger.Error(err, "CNINodeHandler: AddBranchENIToWarmPool failed on getting node", "nodeName", nodeName)
		return err
	}

	if oldCND, err := cnh.client.GetCNINode(types.NamespacedName{Name: node.Name}); err != nil {
		cnh.logger.Error(err, "CNINodeHandler: AddBranchENIToWarmPool failed on getting CNINode", "nodeName", nodeName)
		return err
	} else {
		newCND := oldCND.DeepCopy()
		combinedENIs := combineBranchENIs(newCND.Status.BranchENIs, enis)
		newCND.Status.BranchENIs = combinedENIs
		cnh.logger.Info("adding branch ENIs to CNINode", "enis", combinedENIs, "updatedCNINode", newCND)
		return cnh.client.UpdateCNINode(node, oldCND, newCND)
	}
}

func (cnh *CNINodeHandler) AvailableWarmPoolSpots(nodeName, nodeType string) (int, error) {
	limit := vpc.Limits[nodeType].BranchInterface
	if cnd, err := cnh.client.GetCNINode(types.NamespacedName{Name: nodeName}); err != nil {
		return 0, err
	} else {
		return limit - len(cnd.Status.BranchENIs), nil
	}
}

func (cnh *CNINodeHandler) DeleteBranchENIFromWarmPool(nodeName string, enis []v1alpha1.WarmBranchENI) error {
	node, err := cnh.client.GetNode(nodeName)
	if err != nil {
		return err
	}
	if oldCND, err := cnh.client.GetCNINode(types.NamespacedName{Name: node.Name}); err != nil {
		return err
	} else {
		newCND := oldCND.DeepCopy()
		updatedENIs := deleteBranchENIs(newCND.Status.BranchENIs, enis)
		newCND.Status.BranchENIs = updatedENIs
		return cnh.client.UpdateCNINode(node, oldCND, newCND)
	}
}

func combineBranchENIs(setOne, setTwo []v1alpha1.WarmBranchENI) []v1alpha1.WarmBranchENI {
	eniMap := make(map[string]v1alpha1.WarmBranchENI)
	for _, eni := range setOne {
		eniMap[eni.ID] = eni
	}
	for _, eni := range setTwo {
		eniMap[eni.ID] = eni
	}
	return maps.Values(eniMap)
}

func deleteBranchENIs(ori, del []v1alpha1.WarmBranchENI) []v1alpha1.WarmBranchENI {
	eniMap := make(map[string]v1alpha1.WarmBranchENI)
	for _, eni := range ori {
		eniMap[eni.ID] = eni
	}
	for _, eni := range del {
		delete(eniMap, eni.ID)
	}
	return maps.Values(eniMap)
}
