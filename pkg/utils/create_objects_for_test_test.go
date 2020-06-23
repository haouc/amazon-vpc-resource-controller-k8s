package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestNewPod tests new test testPod creation for unit test.
func TestNewPod(t *testing.T) {
	pod := NewPod(name, saName, namespace)
	assert.True(t,
		name == pod.Name &&
			namespace == pod.Namespace &&
			saName == pod.Spec.ServiceAccountName)
}

// TestNewWindowsPod tests test new windows testPod creation by node selector.
func TestNewWindowsPod(t *testing.T) {
	pod := NewWindowsPod(name, namespace, true)
	assert.True(t,
		pod.Name == name &&
			pod.Namespace == namespace &&
			pod.Spec.NodeSelector["kubernetes.io/os"] == "windows")
}

// TestNewWindowsPod2 tests test new windows testPod creation by node affinity.
func TestNewWindowsPod2(t *testing.T) {
	pod := NewWindowsPod(name, namespace, false)
	exp := &pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.
		NodeSelectorTerms[0].MatchExpressions[0]
	assert.True(t, exp.Key == "kubernetes.io/os" &&
		exp.Values[0] == "windows" &&
		exp.Operator == "In")
}

// TestNewServiceAccount tests new test service account.
func TestNewServiceAccount(t *testing.T) {
	serviceAccount := NewServiceAccount(name, namespace)
	assert.True(t,
		serviceAccount.Name == name &&
			serviceAccount.Namespace == namespace &&
			serviceAccount.Labels["role"] == "db" &&
			serviceAccount.Labels["environment"] == "qa")
}

// TestNewSecurityGroupPolicy tests new test Security Group Policy creation.
func TestNewSecurityGroupPolicy(t *testing.T) {
	sgp := NewSecurityGroupPolicyOne(name, namespace, testSecurityGroupsOne)
	assert.True(t,
		sgp.Name == name &&
			sgp.Namespace == namespace &&
			sgp.Spec.SecurityGroups.Groups[0] == testSecurityGroupsOne[0])
}

// TestNewSecurityGroupPolicyEmptyPodSelector tests new SGP with empty testPod selector.
func TestNewSecurityGroupPolicyEmptyPodSelector(t *testing.T) {
	sgp := NewSecurityGroupPolicyEmptyPodSelector(name, namespace, testSecurityGroupsOne)
	ps, _ := metav1.LabelSelectorAsSelector(sgp.Spec.PodSelector)
	assert.True(t,
		sgp.Name == name &&
			sgp.Namespace == namespace &&
			sgp.Spec.SecurityGroups.Groups[0] == testSecurityGroupsOne[0] &&
			ps.Empty())
}

// TestNewSecurityGroupPolicyEmptySaSelector tests new SGP with empty SA selector.
func TestNewSecurityGroupPolicyEmptySaSelector(t *testing.T) {
	sgp := NewSecurityGroupPolicyEmptySaSelector(name, namespace, testSecurityGroupsOne)
	ls, _ := metav1.LabelSelectorAsSelector(sgp.Spec.ServiceAccountSelector.LabelSelector)
	assert.True(t,
		sgp.Name == name &&
			sgp.Namespace == namespace &&
			sgp.Spec.SecurityGroups.Groups[0] == testSecurityGroupsOne[0] &&
			ls.Empty())
}

// TestNewSecurityGroupPolicyCombined tests new SGP with both SA and testPod selector.
func TestNewSecurityGroupPolicyCombined(t *testing.T) {
	sgp := NewSecurityGroupPolicyCombined(name, namespace, testSecurityGroupsOne)
	ps, _ := metav1.LabelSelectorAsSelector(sgp.Spec.PodSelector)
	ls, _ := metav1.LabelSelectorAsSelector(sgp.Spec.ServiceAccountSelector.LabelSelector)
	assert.True(t,
		sgp.Name == name &&
			sgp.Namespace == namespace &&
			sgp.Spec.SecurityGroups.Groups[0] == testSecurityGroupsOne[0] &&
			!ls.Empty() && !ps.Empty())
}

// TestNewPodWithContainerLimits tests a test pod created with containers.
func TestNewPodWithContainerLimits(t *testing.T) {
	podWithContainers := NewPodWithContainerLimits(name, namespace, true)
	quantity := resource.MustParse("1")
	assert.True(t, len(podWithContainers.Spec.Containers) == 2)
	assert.True(t, podWithContainers.Spec.Containers[1].Resources.Limits["key"] == quantity)
}

// TestNewSecurityGroupPolicyPodSelector tests SGP with pod selector.
func TestNewSecurityGroupPolicyPodSelector(t *testing.T) {
	podSgp := NewSecurityGroupPolicyPodSelector(name, namespace, testSecurityGroupsOne)
	assert.True(t, podSgp.Spec.PodSelector.MatchLabels["role"] == "db")
	exp := podSgp.Spec.PodSelector.MatchExpressions[0]
	assert.True(t,
		exp.Key == "environment" &&
			exp.Operator == "In" &&
			exp.Values[0] == "qa")
}

// TestNewSecurityGroupPolicySaSelector tests SGP with SA selector.
func TestNewSecurityGroupPolicySaSelector(t *testing.T) {
	saSgp := NewSecurityGroupPolicySaSelector(name, namespace, testSecurityGroupsOne)
	assert.True(t, saSgp.Spec.ServiceAccountSelector.MatchNames[0] == testSA.Name)
	assert.True(t, saSgp.Spec.ServiceAccountSelector.MatchLabels["role"] == "db")
	exp := saSgp.Spec.ServiceAccountSelector.LabelSelector.MatchExpressions[0]
	assert.True(t,
		exp.Key == "environment" &&
			exp.Operator == "In" &&
			exp.Values[0] == "qa")
}
