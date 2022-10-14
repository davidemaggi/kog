package k8s

import (
	"github.com/davidemaggi/kog/pkg/common"
	"k8s.io/client-go/kubernetes"
)

var kubeconfig string

func getKubeClient(configPath string, verbose bool) (client *kubernetes.Clientset, err error) {

	config, err := buildConfigFromPath(configPath, verbose)

	if !common.HandleError(err, "Error Building Config", true) {
		return nil, err

	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)

	return clientset, err

}
