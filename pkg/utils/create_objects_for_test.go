package utils

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	vpcresourcesv1beta1 "github.com/aws/amazon-vpc-resource-controller-k8s/apis/vpcresources/v1beta1"
)

// NewPod creates a regular pod for test
func NewPod(name string, sa string, namespace string) *corev1.Pod {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels: map[string]string{
				"role":        "db",
				"environment": "qa",
				"app":         "test_app",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Resources: corev1.ResourceRequirements{},
				},
			},
			ServiceAccountName: sa,
		},
		Status: corev1.PodStatus{
			Conditions: []corev1.PodCondition{
				{
					Type:   corev1.PodReady,
					Status: corev1.ConditionTrue,
				},
			},
		},
	}
	return pod
}

// NewPodNotForENI creates a regular pod no need for ENI for test.
func NewPodNotForENI(name string, sa string, namespace string) *corev1.Pod {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels: map[string]string{
				"app": "test_app",
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: sa,
		},
	}
	return pod
}

// NewPodForMultiENI creates a regular pod for ENIs for test.
func NewPodForMultiENI(name string, sa string, namespace string) *corev1.Pod {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels: map[string]string{
				"app":         "vpc-controller",
				"role":        "db",
				"environment": "qa",
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: sa,
		},
	}
	return pod
}

// NewWindowsPod creates a windows pod for test.
// Parameter useSelector can set if using nodeSelector or nodeAffinity for OS type.
func NewWindowsPod(name string, namespace string, useSelector bool) *corev1.Pod {
	var spec corev1.PodSpec
	containers := []corev1.Container{
		{
			Resources: corev1.ResourceRequirements{},
		},
	}

	if useSelector {
		spec = corev1.PodSpec{
			Containers: containers,
			NodeSelector: map[string]string{
				"kubernetes.io/os": "windows",
			},
		}
	} else {
		spec = corev1.PodSpec{
			Containers: containers,
			Affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "kubernetes.io/os",
										Operator: "In",
										Values:   []string{"windows"},
									},
								},
								MatchFields: nil,
							},
						},
					},
				},
			},
		}
	}

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels: map[string]string{
				"role":        "db",
				"environment": "qa",
				"app":         "test_app",
			},
		},
		Spec: spec,
		Status: corev1.PodStatus{
			Conditions: []corev1.PodCondition{
				{
					Type:   corev1.PodReady,
					Status: corev1.ConditionTrue,
				},
			},
		},
	}
	return pod
}

// NewPodWithContainerLimits creates a test pod with containers that already have resources limits.
// The parameter hasLimit can set if resource limit is set.
func NewPodWithContainerLimits(name string, namespace string, hasLimit bool) *corev1.Pod {
	pod := NewPod(name, "", namespace)
	limit := corev1.ResourceList{}
	if hasLimit {
		limit["key"] = resource.MustParse("1")
	}
	pod.Spec.Containers = []corev1.Container{
		corev1.Container{
			Name: "test_container_1",
			Resources: corev1.ResourceRequirements{
				Limits: nil,
			},
		},
		corev1.Container{
			Name: "test_container_1",
			Resources: corev1.ResourceRequirements{
				Limits: limit,
			},
		},
	}
	return pod
}

// NewSecurityGroupPolicyCombined creates a test SGP with both pod and SA selectors.
func NewSecurityGroupPolicyCombined(
	name string, namespace string, securityGroups []string) vpcresourcesv1beta1.SecurityGroupPolicy {
	sgp := vpcresourcesv1beta1.SecurityGroupPolicy{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: vpcresourcesv1beta1.SecurityGroupPolicySpec{
			PodSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"role": "db"},
				MatchExpressions: []metav1.LabelSelectorRequirement{{
					Key:      "environment",
					Operator: "In",
					Values:   []string{"qa", "production"},
				}},
			},
			ServiceAccountSelector: vpcresourcesv1beta1.ServiceAccountSelector{
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"role": "db"},
					MatchExpressions: []metav1.LabelSelectorRequirement{{
						Key:      "environment",
						Operator: "In",
						Values:   []string{"qa", "production"},
					}},
				},
				MatchNames: []string{"test_sa"},
			},
			SecurityGroups: vpcresourcesv1beta1.GroupIds{
				Groups: securityGroups,
			},
		},
	}
	return sgp
}

// NewSecurityGroupPolicyPodSelector creates a test SGP with only pod selector.
func NewSecurityGroupPolicyPodSelector(
	name string, namespace string, securityGroups []string) vpcresourcesv1beta1.SecurityGroupPolicy {
	sgp := vpcresourcesv1beta1.SecurityGroupPolicy{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: vpcresourcesv1beta1.SecurityGroupPolicySpec{
			PodSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"role": "db"},
				MatchExpressions: []metav1.LabelSelectorRequirement{{
					Key:      "environment",
					Operator: "In",
					Values:   []string{"qa", "production"},
				}},
			},
			SecurityGroups: vpcresourcesv1beta1.GroupIds{
				Groups: securityGroups,
			},
		},
	}
	return sgp
}

// NewSecurityGroupPolicyEmptyPodSelector creates a test SGP with only empty pod selector.
func NewSecurityGroupPolicyEmptyPodSelector(name string, namespace string, securityGroups []string) vpcresourcesv1beta1.SecurityGroupPolicy {
	sgp := vpcresourcesv1beta1.SecurityGroupPolicy{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: vpcresourcesv1beta1.SecurityGroupPolicySpec{
			PodSelector: &metav1.LabelSelector{
				MatchLabels:      map[string]string{},
				MatchExpressions: []metav1.LabelSelectorRequirement{},
			},
			SecurityGroups: vpcresourcesv1beta1.GroupIds{
				Groups: securityGroups,
			},
		},
	}
	return sgp
}

// NewSecurityGroupPolicySaSelector creates a test SGP with only SA selector.
func NewSecurityGroupPolicySaSelector(name string, namespace string, securityGroups []string) vpcresourcesv1beta1.SecurityGroupPolicy {
	sgp := vpcresourcesv1beta1.SecurityGroupPolicy{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: vpcresourcesv1beta1.SecurityGroupPolicySpec{
			ServiceAccountSelector: vpcresourcesv1beta1.ServiceAccountSelector{
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"role": "db"},
					MatchExpressions: []metav1.LabelSelectorRequirement{{
						Key:      "environment",
						Operator: "In",
						Values:   []string{"qa", "production"},
					}},
				},
				MatchNames: []string{"test_sa"},
			},
			SecurityGroups: vpcresourcesv1beta1.GroupIds{
				Groups: securityGroups,
			},
		},
	}
	return sgp
}

// NewSecurityGroupPolicyEmptySaSelector creates a test SGP with only empty SA selector.
func NewSecurityGroupPolicyEmptySaSelector(name string, namespace string, securityGroups []string) vpcresourcesv1beta1.SecurityGroupPolicy {
	sgp := vpcresourcesv1beta1.SecurityGroupPolicy{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: vpcresourcesv1beta1.SecurityGroupPolicySpec{
			ServiceAccountSelector: vpcresourcesv1beta1.ServiceAccountSelector{
				LabelSelector: &metav1.LabelSelector{
					MatchLabels:      map[string]string{},
					MatchExpressions: []metav1.LabelSelectorRequirement{},
				},
				MatchNames: []string{"test_sa"},
			},
			SecurityGroups: vpcresourcesv1beta1.GroupIds{
				Groups: securityGroups,
			},
		},
	}
	return sgp
}

// NewServiceAccount creates a test service account.
func NewServiceAccount(name string, namespace string) *corev1.ServiceAccount {
	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"role":        "db",
				"environment": "qa",
			},
		},
		Secrets:                      nil,
		ImagePullSecrets:             nil,
		AutomountServiceAccountToken: nil,
	}
	return sa
}

// NewSecurityGroupPolicyOne creates a test security group policy's pointer.
func NewSecurityGroupPolicyOne(name string, namespace string, securityGroups []string) *vpcresourcesv1beta1.SecurityGroupPolicy {
	sgp := &vpcresourcesv1beta1.SecurityGroupPolicy{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: vpcresourcesv1beta1.SecurityGroupPolicySpec{
			PodSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"role": "db"},
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "environment",
						Operator: "In",
						Values:   []string{"qa", "production"},
					},
				},
			},
			SecurityGroups: vpcresourcesv1beta1.GroupIds{
				Groups: securityGroups,
			},
		},
	}
	return sgp
}

// NewSecurityGroupPolicyTwo creates a test security group policy's pointer.
func NewSecurityGroupPolicyTwo(name string, namespace string, securityGroups []string) *vpcresourcesv1beta1.SecurityGroupPolicy {
	sgp := &vpcresourcesv1beta1.SecurityGroupPolicy{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: vpcresourcesv1beta1.SecurityGroupPolicySpec{
			PodSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "vpc-controller"},
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "environment",
						Operator: "In",
						Values:   []string{"qa", "production"},
					},
				},
			},
			SecurityGroups: vpcresourcesv1beta1.GroupIds{
				Groups: securityGroups,
			},
		},
	}
	return sgp
}
