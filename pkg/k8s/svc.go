package k8s

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/davidemaggi/kog/pkg/common"
	"golang.org/x/net/context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

func GetServices(configPath string, verbose bool) (servicesString []string, services []v1.Service, err error) {
	kc, err := getKubeClient(configPath, verbose)

	var svcList, errNs = kc.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})

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

func PortForwardPod(configPath string, pod *v1.Pod, fwdPort int32, localPort int32, verbose bool) (err error) {

	config, err := buildConfigFromPath(configPath, verbose)
	if !common.HandleError(err, "Error Building Config", verbose) {
		return err
	}

	path := ""

	path = fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward",
		pod.Namespace, pod.Name)

	hostIP := strings.TrimLeft(config.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(config)
	if !common.HandleError(err, "Error Setting up port forward", verbose) {
		return err
	}

	var Streams genericclioptions.IOStreams
	// StopCh is the channel used to manage the port forward lifecycle
	var StopCh <-chan struct{}
	// ReadyCh communicates when the tunnel is ready to receive traffic
	var ReadyCh chan struct{}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", localPort, fwdPort)}, StopCh, ReadyCh, Streams.Out, Streams.ErrOut)
	if !common.HandleError(err, "Error Executing port forward", verbose) {
		return err
	}
	return fw.ForwardPorts()
}

func getPodsForSvc(svc *v1.Service, namespace string, k8sClient *kubernetes.Clientset) (*v1.PodList, error) {
	set := labels.Set(svc.Spec.Selector)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	pods, err := k8sClient.CoreV1().Pods(namespace).List(context.TODO(), listOptions)

	return pods, err
}
