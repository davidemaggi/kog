package k8s

import (
	"errors"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
)

var kubeconfig string

func GetConfig() (kubeConfig *api.Config, err error) {

	kubeConfigPath, err := FindKubeConfig()
	if err != nil {
		log.Fatal(err)
	}

	kubeConfig, err = clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//log.Printf("current context is %s", kubeConfig.CurrentContext)
	return kubeConfig, nil
	//err = clientcmd.WriteToFile(*kubeConfig, "copy.yaml")
	//if err != nil {
	//	log.Fatal(err)
	//}

}

func SaveConfig(filePath string, kubeConfig *api.Config, backup bool) (msg string, err error) {

	if backup {
		oldConfig, _ := GetConfig()
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

func GetContexts(verbose bool) (ctxs []string, err error) {
	ctxs = make([]string, 0)

	kubeConfig, err := GetConfig()

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return ctxs, err
	}

	for s, _ := range kubeConfig.Contexts {
		ctxs = append(ctxs, s)
	}
	return ctxs, nil

}

func GetCurrentContext(verbose bool) (ctx string, err error) {

	kubeConfig, err := GetConfig()
	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return "", err
	}

	return kubeConfig.CurrentContext, nil
}

func SetContext(ctx string, verbose bool) (err error) {

	kubeConfig, err := GetConfig()

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return err
	}

	kubeConfig.CurrentContext = ctx
	path, _ := FindKubeConfig()
	s, err := SaveConfig(path, kubeConfig, true)
	_ = s
	return nil

}

func SetNameSpace(ns string, verbose bool) (err error) {

	kubeConfig, err := GetConfig()

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

func GetNamespaces(verbose bool) (namespaces []string, err error) {
	kubeConfigPath, err := FindKubeConfig()

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return namespaces, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
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
