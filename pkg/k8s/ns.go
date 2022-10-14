package k8s

import (
	"github.com/davidemaggi/kog/pkg/common"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SetNameSpace(ns string, verbose bool) (err error) {

	_, kubeConfig, err := GetConfig("")

	if !common.HandleError(err, "Error Getting Config", true) {
		return err

	}

	kubeConfig.Contexts[kubeConfig.CurrentContext].Namespace = ns

	path, _ := FindKubeConfig()
	s, err := SaveConfig(path, kubeConfig, true)
	_ = s
	return nil

}

func GetNamespaces(configPath string, verbose bool) (namespaces []string, err error) {

	// create the clientset
	clientset, err := getKubeClient(configPath, verbose)

	if !common.HandleError(err, "Error Getting Client", verbose) {
		return namespaces, err

	}

	var nsList, errNs = clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if !common.HandleError(errNs, "Error Getting Namespaces", verbose) {
		return namespaces, errNs

	}
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.Name)
	}
	return namespaces, nil

}

func GetCurrentNamespace(configPath string, verbose bool) (namespace string, err error) {

	// create the clientset
	_, config, err := GetConfig(configPath)

	if !common.HandleError(err, "Error Getting Config", verbose) {
		return "", err

	}

	curCtx, err := GetCurrentContext(configPath, verbose)

	if !common.HandleError(err, "Error Getting Current Context", verbose) {
		return "", err

	}

	if config.Contexts[curCtx] != nil {
		return config.Contexts[curCtx].Namespace, nil

	}

	return "", nil

}
