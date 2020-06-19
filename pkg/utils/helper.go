package utils

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	vpcresourcesv1beta1 "github.com/aws/amazon-vpc-resource-controller-k8s/apis/vpcresources/v1beta1"
)

// Include checks if a string existing in a string slice and returns true or false.
func Include(target string, values []string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

// RemoveDuplicatedSg removes duplicated items from a string slice.
// It returns a no duplicates string slice.
func RemoveDuplicatedSg(list []string) []string {
	set := make(map[string]bool)
	var processedList []string
	for _, sg := range list {
		if _, ok := set[sg]; !ok {
			processedList = append(processedList, sg)
			set[sg] = true
		}
	}
	return processedList
}

// ShouldAddEniLimits decide if the testPod can be mutated to inject ENI annotation for security groups.
// The function returns security group list and true or false for mutating testPod.
func (kch *K8sCacheHelper) ShouldAddEniLimits(pod *corev1.Pod) []string {
	helperLog := kch.Log.WithValues("Pod name", pod.Name, "Pod namespace", pod.Namespace)

	// Build SGP list from cache.
	ctx := context.Background()
	sgpList := &vpcresourcesv1beta1.SecurityGroupPolicyList{}
	if err := kch.Client.List(ctx, sgpList, &client.ListOptions{Namespace: pod.Namespace}); err != nil {
		panic(err)
	}

	sa := &corev1.ServiceAccount{}
	key := types.NamespacedName{
		Namespace: pod.Namespace,
		Name:      pod.Spec.ServiceAccountName}

	// Get metadata of SA associated with Pod from cache
	if err := kch.Client.Get(ctx, key, sa); err != nil {
		panic(err)
	}

	if sgList := canGetSecurityGroups(sgpList, pod, sa, helperLog); len(sgList) == 0 {
		return nil
	} else {
		helperLog.Info("Pod got security groups from matched SGP.",
			"Security Groups", sgList)
		return sgList
	}
}

func canGetSecurityGroups(
	sgpList *vpcresourcesv1beta1.SecurityGroupPolicyList,
	pod *corev1.Pod,
	sa *corev1.ServiceAccount,
	logger logr.Logger) []string {
	var sgList []string
	for _, sgp := range sgpList.Items {
		hasPodSelector := sgp.Spec.PodSelector != nil
		hasSASelector := sgp.Spec.ServiceAccountSelector.MatchNames != nil ||
			sgp.Spec.ServiceAccountSelector.LabelSelector != nil
		if !hasPodSelector && !hasSASelector {
			logger.Info("Found an invalid SecurityGroupPolicy due to both of podSelector and saSelector are null.",
				"Invalid SGP", types.NamespacedName{Name: sgp.Name, Namespace: sgp.Namespace})
			continue
		}

		podMatched, saMatched := false, false
		podSelector, _ := metav1.LabelSelectorAsSelector(sgp.Spec.PodSelector)
		if podSelector.Matches(labels.Set(pod.Labels)) {
			podMatched = true
		}

		saSelector, _ := metav1.LabelSelectorAsSelector(sgp.Spec.ServiceAccountSelector.LabelSelector)
		if Include(sa.Name, sgp.Spec.ServiceAccountSelector.MatchNames) &&
			saSelector.Matches(labels.Set(sa.Labels)) {
			saMatched = true
		}

		if ((hasPodSelector && podMatched) && (hasSASelector && saMatched)) ||
			(!hasPodSelector && (hasSASelector && saMatched)) ||
			(!hasSASelector && (hasPodSelector && podMatched)) {
			if sgp.Spec.SecurityGroups.Groups != nil {
				sgList = append(sgList, sgp.Spec.SecurityGroups.Groups...)
			}
		}
	}

	sgList = RemoveDuplicatedSg(sgList)
	return sgList
}
