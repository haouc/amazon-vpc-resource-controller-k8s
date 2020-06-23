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
	testSA                *corev1.ServiceAccount
	testPod               *corev1.Pod
	testScheme            *runtime.Scheme
	testClient            client.Client
	testSecurityGroupsOne []string
	testSecurityGroupsTwo []string
	helper                *K8sCacheHelper
	name                  string
	namespace             string
	saName                string
)

func init() {
	name = "test"
	namespace = "test_namespace"
	saName = "test_sa"
	testSA = NewServiceAccount(saName, namespace)
	testPod = NewPod(name, saName, namespace)
	testScheme = runtime.NewScheme()
	clientgoscheme.AddToScheme(testScheme)
	vpcresourcesv1beta1.AddToScheme(testScheme)

	testSecurityGroupsOne = []string{"sg-00001", "sg-00002"}
	testSecurityGroupsTwo = []string{"sg-00003", "sg-00004"}
	testClient = fake.NewFakeClientWithScheme(
		testScheme,
		NewPod(name, saName, namespace),
		NewPodNotForENI(name+"_NoENI", saName, namespace),
		NewPodForMultiENI(name+"_ENIs", saName, namespace),
		NewServiceAccount(saName, namespace),
		NewSecurityGroupPolicyOne(name+"_1", namespace, testSecurityGroupsOne),
		NewSecurityGroupPolicyTwo(name+"_2", namespace, append(testSecurityGroupsOne, testSecurityGroupsTwo...)),
	)

	helper = &K8sCacheHelper{
		Client: testClient,
		Log:    ctrl.Log.WithName("testLog"),
	}
}
