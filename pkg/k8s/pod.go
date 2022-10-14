package k8s

import (
	"github.com/davidemaggi/kog/pkg/common"
	"golang.org/x/net/context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPods(configPath string, verbose bool) (podsString []string, pods []v1.Pod, err error) {

	kc, err := getKubeClient(configPath, verbose)

	var podsList, errNs = kc.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

	if !common.HandleError(errNs, "Error Getting Pods", verbose) {
		return podsString, nil, errNs
	}

	for _, pod := range podsList.Items {
		podsString = append(podsString, pod.Name)
	}
	return podsString, podsList.Items, nil

}
