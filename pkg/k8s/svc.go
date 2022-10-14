package k8s

import (
	"github.com/davidemaggi/kog/pkg/common"
	"golang.org/x/net/context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetServices(configPath string, verbose bool) (servicesString []string, services []v1.Service, err error) {
	kc, err := getKubeClient(configPath, verbose)
	ns, err := GetCurrentNamespace(configPath, verbose)

	var svcList, errNs = kc.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{})

	if !common.HandleError(errNs, "Error Getting Services", verbose) {
		return servicesString, nil, errNs
	}

	for _, svc := range svcList.Items {
		servicesString = append(servicesString, svc.Name)

	}
	return servicesString, svcList.Items, nil

}

func PortForwardSvc(configPath string, svc *v1.Service, fwdPort int32, localPort int32, verbose bool) (err error) {

	kc, err := getKubeClient(configPath, verbose)

	if !common.HandleError(err, "Error Getting Client", verbose) {
		return err
	}

	pods, err := getPodsForSvc(svc, svc.Namespace, kc)

	if !common.HandleError(err, "Error Getting Pods for Svc", verbose) {
		return err
	}

	PortForwardPod(configPath, &pods.Items[0], fwdPort, localPort, verbose)

	return nil

}
