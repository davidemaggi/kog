package k8s

import (
	"errors"
	"fmt"
	"github.com/davidemaggi/kog/structs"
	"golang.org/x/net/context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/client-go/util/homedir"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var kubeconfig string

func GetConfig(kubeConfigPath string) (path string, kubeConfig *api.Config, err error) {

	if kubeConfigPath == "" {
		kubeConfigPath, err = FindKubeConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	kubeConfig, err = clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		log.Fatal(err)
		return kubeConfigPath, nil, err
	}

	//log.Printf("current context is %s", kubeConfig.CurrentContext)
	return kubeConfigPath, kubeConfig, nil
	//err = clientcmd.WriteToFile(*kubeConfig, "copy.yaml")
	//if err != nil {
	//	log.Fatal(err)
	//}

}

func PrintConfig(kubeConfigPath string) (err error) {

	if kubeConfigPath == "" {
		kubeConfigPath, err = FindKubeConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	b, err := os.ReadFile(kubeConfigPath)
	if err != nil {
		fmt.Print(err)
	}

	str := string(b) // convert content to a 'string'

	fmt.Println(str) // print the content as a 'string'

	//log.Printf("current context is %s", kubeConfig.CurrentContext)
	return nil
	//err = clientcmd.WriteToFile(*kubeConfig, "copy.yaml")
	//if err != nil {
	//	log.Fatal(err)
	//}

}

func SaveConfig(filePath string, kubeConfig *api.Config, backup bool) (msg string, err error) {

	if backup {
		_, oldConfig, _ := GetConfig("")
		err = clientcmd.WriteToFile(*oldConfig, filePath+".backup")
		if err != nil {
			log.Fatal(err)
			return "Error during backup", err
		}
	}

	err = clientcmd.WriteToFile(*kubeConfig, filePath)
	if err != nil {
		log.Fatal(err)
		return "Something happened", err
	}
	return "ok", nil
}

func FindKubeConfig() (string, error) {
	env := os.Getenv("KUBECONFIG")
	if env != "" {
		return env, nil
	}

	if home := homedir.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config"), nil
	}

	return "", errors.New("Error Retrieving KubeConfig")
}

func GetContexts(configPath string, verbose bool) (ctxs []string, rawCtxs map[string]*api.Context, err error) {
	ctxs = make([]string, 0)

	_, kubeConfig, err := GetConfig(configPath)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return ctxs, nil, err
	}

	for s, _ := range kubeConfig.Contexts {
		ctxs = append(ctxs, s)
	}
	return ctxs, kubeConfig.Contexts, nil

}

func GetCurrentContext(configPath string, verbose bool) (ctx string, err error) {

	_, kubeConfig, err := GetConfig(configPath)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return "", err
	}

	return kubeConfig.CurrentContext, nil
}

func SetContext(configPath string, ctx string, verbose bool) (err error) {

	_, kubeConfig, err := GetConfig(configPath)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return err
	}

	kubeConfig.CurrentContext = ctx
	if configPath == "" {
		configPath, _ = FindKubeConfig()
	}

	s, err := SaveConfig(configPath, kubeConfig, true)
	_ = s
	return nil

}

func DeleteContext(configPath string, ctx string, verbose bool) (err error) {

	_, kubeConfig, err := GetConfig(configPath)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return err
	}

	if ctx == kubeConfig.CurrentContext {

		//You are deleting the current default
		kubeConfig.CurrentContext = ""

	}

	delete(kubeConfig.Contexts, ctx)

	if configPath == "" {
		configPath, _ = FindKubeConfig()
	}

	s, err := SaveConfig(configPath, kubeConfig, true)

	_ = s
	return nil

}

func RenameContext(configPath string, ctx string, newName string, verbose bool) (err error) {

	_, kubeConfig, err := GetConfig(configPath)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return err
	}

	if ctx == kubeConfig.CurrentContext {

		//You are renaming the current default
		kubeConfig.CurrentContext = newName

	}

	kubeConfig.Contexts[newName] = kubeConfig.Contexts[ctx]
	delete(kubeConfig.Contexts, ctx)

	if configPath == "" {
		configPath, _ = FindKubeConfig()
	}

	s, err := SaveConfig(configPath, kubeConfig, true)

	_ = s
	return nil

}

func SetNameSpace(ns string, verbose bool) (err error) {

	_, kubeConfig, err := GetConfig("")

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return err
	}

	kubeConfig.Contexts[kubeConfig.CurrentContext].Namespace = ns

	path, _ := FindKubeConfig()
	s, err := SaveConfig(path, kubeConfig, true)
	_ = s
	return nil

}

func GetNamespaces(configPath string, verbose bool) (namespaces []string, err error) {
	if configPath == "" {
		configPath, err = FindKubeConfig()
	}

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return namespaces, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return namespaces, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return namespaces, err

	}

	var nsList, errNs = clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if errNs != nil {

		if verbose {
			log.Fatal(errNs)
		}
		return namespaces, errNs

	}
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.Name)
	}
	return namespaces, nil

}

func MergeConfigs(newFile string, oldFile string, force bool, verbose bool) (kubeConfig *api.Config, result structs.MergeResult, err error) {

	result = structs.New_MergeResult()
	oldPath, oldConfig, err := GetConfig(oldFile)

	result.From = newFile
	result.To = oldPath

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return nil, result, err

	}
	origConfig := oldConfig.DeepCopy()
	_, newConfig, err := GetConfig(newFile)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return oldConfig, result, err

	}
	//region Context
	for ctx_id, ctx := range newConfig.Contexts {
		// index is the index where we are
		// element is the element from someSlice for where we are

		exists := oldConfig.Contexts[ctx_id] != nil

		if exists {

			if force {
				result.DoneSomething = true
				oldConfig.Contexts[ctx_id] = ctx.DeepCopy()
				structs.AddAction(result, structs.Modified, structs.Context)
			}

		} else {
			result.DoneSomething = true
			oldConfig.Contexts[ctx_id] = ctx.DeepCopy()
			structs.AddAction(result, structs.Added, structs.Context)

		}

	}
	//endregion

	//region Cluster
	for cl_id, cl := range newConfig.Clusters {
		// index is the index where we are
		// element is the element from someSlice for where we are

		exists := oldConfig.Clusters[cl_id] != nil

		if exists {

			if force {
				result.DoneSomething = true
				oldConfig.Clusters[cl_id] = cl.DeepCopy()
				structs.AddAction(result, structs.Modified, structs.Cluster)

			}

		} else {
			result.DoneSomething = true
			oldConfig.Clusters[cl_id] = cl.DeepCopy()
			structs.AddAction(result, structs.Added, structs.Cluster)

		}

	}
	//endregion

	//region User
	for usr_id, usr := range newConfig.AuthInfos {
		// index is the index where we are
		// element is the element from someSlice for where we are

		exists := oldConfig.AuthInfos[usr_id] != nil

		if exists {

			if force {
				result.DoneSomething = true
				oldConfig.AuthInfos[usr_id] = usr.DeepCopy()
				structs.AddAction(result, structs.Modified, structs.User)

			}

		} else {
			result.DoneSomething = true
			oldConfig.AuthInfos[usr_id] = usr.DeepCopy()
			structs.AddAction(result, structs.Added, structs.User)

		}

	}

	_ = origConfig
	result.IsOk = true

	return oldConfig, result, nil
}

func GetPods(configPath string, verbose bool) (podsString []string, pods []v1.Pod, err error) {
	if configPath == "" {
		configPath, err = FindKubeConfig()
	}

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return podsString, nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return podsString, nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return podsString, nil, err

	}

	var podsList, errNs = clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if errNs != nil {

		if verbose {
			log.Fatal(errNs)
		}
		return podsString, nil, errNs

	}
	for _, pod := range podsList.Items {
		podsString = append(podsString, pod.Name)
	}
	return podsString, podsList.Items, nil

}

func GetServices(configPath string, verbose bool) (servicesString []string, services []v1.Service, err error) {
	if configPath == "" {
		configPath, err = FindKubeConfig()
	}

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return servicesString, nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return servicesString, nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return servicesString, nil, err

	}

	var svcList, errNs = clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if errNs != nil {

		if verbose {
			log.Fatal(errNs)
		}
		return servicesString, nil, errNs

	}
	for _, svc := range svcList.Items {
		servicesString = append(servicesString, svc.Name)

	}
	return servicesString, svcList.Items, nil

}

func PortForwardSvc(configPath string, svc v1.Service, fwdPort int, localPort int, verbose bool) (err error) {

	if configPath == "" {
		configPath, err = FindKubeConfig()
	}

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return err
	}
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return err
	}

	path := ""

	path = fmt.Sprintf("/api/v1/namespaces/%s/services/%s/portforward",
		svc.Namespace, svc.Name)

	hostIP := strings.TrimLeft(config.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return err
	}

	var Streams genericclioptions.IOStreams
	// StopCh is the channel used to manage the port forward lifecycle
	var StopCh <-chan struct{}
	// ReadyCh communicates when the tunnel is ready to receive traffic
	var ReadyCh chan struct{}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", localPort, fwdPort)}, StopCh, ReadyCh, Streams.Out, Streams.ErrOut)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}

func PortForwardPod(configPath string, pod v1.Pod, fwdPort int, localPort int, verbose bool) (err error) {

	if configPath == "" {
		configPath, err = FindKubeConfig()
	}

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return err
	}
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return err
	}

	path := ""

	path = fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward",
		pod.Namespace, pod.Name)

	hostIP := strings.TrimLeft(config.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		return err
	}

	var Streams genericclioptions.IOStreams
	// StopCh is the channel used to manage the port forward lifecycle
	var StopCh <-chan struct{}
	// ReadyCh communicates when the tunnel is ready to receive traffic
	var ReadyCh chan struct{}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", localPort, fwdPort)}, StopCh, ReadyCh, Streams.Out, Streams.ErrOut)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}
