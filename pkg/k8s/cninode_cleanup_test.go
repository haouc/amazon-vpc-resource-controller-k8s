package k8s

import (
	"context"
	"testing"

	"github.com/aws/amazon-vpc-resource-controller-k8s/apis/vpcresources/v1alpha1"
	"github.com/aws/amazon-vpc-resource-controller-k8s/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	mockNodeList = &v1.NodeList{
		Items: []v1.Node{*mockNode},
	}

	mockCNINodeList = &v1alpha1.CNINodeList{
		Items: []v1alpha1.CNINode{*mockCNINode},
	}
)

func TestCNINodeCleaner_cleanUpLeakedCNINodes_No_Leak(t *testing.T) {
	ctrl := gomock.NewController(t)
	_, k8sClient, _ := getMockK8sWrapperWithClient(ctrl, []runtime.Object{mockNodeList, mockCNINodeList})

	cleaner := &CNINodeCleaner{
		Log:            zap.New(),
		k8sClient:      k8sClient,
		ctx:            context.Background(),
		leakedCNINodes: []string{"test-node-leaked"},
	}

	cleaner.cleanUpLeakedCNINodes()
	assert.True(t, len(cleaner.leakedCNINodes) == 0)
}

func TestCNINodeCleaner_cleanUpLeakedCNINodes_Has_Leak(t *testing.T) {
	ctrl := gomock.NewController(t)
	leakedCNINode := &v1alpha1.CNINode{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-node-leaked",
			Namespace: config.KubeDefaultNamespace,
		},
	}
	mockCNINodeLeakedList := &v1alpha1.CNINodeList{
		Items: []v1alpha1.CNINode{
			*mockCNINode,
			*leakedCNINode,
		},
	}
	_, k8sClient, _ := getMockK8sWrapperWithClient(ctrl, []runtime.Object{mockNodeList, mockCNINodeLeakedList})

	cleaner := &CNINodeCleaner{
		Log:       zap.New(),
		k8sClient: k8sClient,
		ctx:       context.Background(),
	}

	// the leaked CNINode list is empty and adding the found leaked CNINode to the list but not cleanup yet
	cleaner.cleanUpLeakedCNINodes()
	var wantedCNINode v1alpha1.CNINode
	err := k8sClient.Get(cleaner.ctx, types.NamespacedName{Name: "test-node-leaked", Namespace: config.KubeDefaultNamespace}, &wantedCNINode)
	assert.True(t, err == nil)
	assert.True(t, len(cleaner.leakedCNINodes) == 1)
	var leakedNode v1.Node
	err = k8sClient.Get(cleaner.ctx, types.NamespacedName{Name: "test-node-leaked"}, &leakedNode)
	assert.True(t, apierrors.IsNotFound(err))

	// simulate the next run in 12 hours, the cached leaked CNINode will be cleaned if no node exists
	cleaner.cleanUpLeakedCNINodes()
	err = k8sClient.Get(cleaner.ctx, types.NamespacedName{Name: "test-node-leaked", Namespace: config.KubeDefaultNamespace}, &wantedCNINode)
	assert.True(t, apierrors.IsNotFound(err))
	assert.True(t, len(cleaner.leakedCNINodes) == 0)
}

func TestCNINodeCleaner_cleanUpLeakedCNINodes_Has_Leak_False_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	leakedCNINode := &v1alpha1.CNINode{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-node-leaked",
			Namespace: config.KubeDefaultNamespace,
		},
	}
	mockCNINodeLeakedList := &v1alpha1.CNINodeList{
		Items: []v1alpha1.CNINode{
			*mockCNINode,
			*leakedCNINode,
		},
	}
	_, k8sClient, _ := getMockK8sWrapperWithClient(ctrl, []runtime.Object{mockNodeList, mockCNINodeLeakedList})

	cleaner := &CNINodeCleaner{
		Log:       zap.New(),
		k8sClient: k8sClient,
		ctx:       context.Background(),
	}

	// the leaked CNINode list is empty and adding the found leaked CNINode to the list but not cleanup yet
	cleaner.cleanUpLeakedCNINodes()
	var wantedCNINode v1alpha1.CNINode
	err := k8sClient.Get(cleaner.ctx, types.NamespacedName{Name: "test-node-leaked", Namespace: config.KubeDefaultNamespace}, &wantedCNINode)
	assert.True(t, err == nil)
	assert.True(t, len(cleaner.leakedCNINodes) == 1)
	var leakedNode v1.Node
	err = k8sClient.Get(cleaner.ctx, types.NamespacedName{Name: "test-node-leaked"}, &leakedNode)
	assert.True(t, apierrors.IsNotFound(err))

	// for some reason the node was not available
	node := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-node-leaked",
		},
	}
	err = k8sClient.Create(cleaner.ctx, node)
	assert.NoError(t, err)
	// simulate false alarm
	// in the next run in 12 hours, the cached leaked CNINode will not be cleaned
	cleaner.cleanUpLeakedCNINodes()
	err = k8sClient.Get(cleaner.ctx, types.NamespacedName{Name: "test-node-leaked", Namespace: config.KubeDefaultNamespace}, &wantedCNINode)
	assert.NoError(t, err)
	assert.True(t, len(cleaner.leakedCNINodes) == 0)
}
