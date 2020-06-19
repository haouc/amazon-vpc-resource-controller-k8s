package utils

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	vpcresourcesv1beta1 "github.com/aws/amazon-vpc-resource-controller-k8s/apis/vpcresources/v1beta1"
)

var (
	testSA             *corev1.ServiceAccount
	testPod            *corev1.Pod
	TestScheme         *runtime.Scheme
	TestClient         client.Client
	testSecurityGroups []string
	helper             *K8sCacheHelper
	name               string
	namespace          string
	saName             string
)

func init() {
	name = "test"
	namespace = "test_namespace"
	saName = "test_sa"
	testSA = NewServiceAccount(saName, namespace)
	testPod = NewPod(name, saName, namespace)
	TestScheme = runtime.NewScheme()
	clientgoscheme.AddToScheme(TestScheme)
	vpcresourcesv1beta1.AddToScheme(TestScheme)

	testSecurityGroups = []string{"sg-00001", "sg-00002"}
	TestClient = fake.NewFakeClientWithScheme(
		TestScheme,
		NewPod(name, saName, namespace),
		NewServiceAccount(saName, namespace),
		NewSecurityGroupPolicy(name, namespace, testSecurityGroups),
	)

	helper = &K8sCacheHelper{
		Client: TestClient,
		Log:    ctrl.Log.WithName("testLog"),
	}
}
