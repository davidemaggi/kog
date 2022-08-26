package k8s

import (
	"errors"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
)

var kubeconfig string

func GetConfig() {

	kubeConfigPath, err := FindKubeConfig()
	if err != nil {
		log.Fatal(err)
	}

	kubeConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("current context is %s", kubeConfig.CurrentContext)

	//err = clientcmd.WriteToFile(*kubeConfig, "copy.yaml")
	//if err != nil {
	//	log.Fatal(err)
	//}

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
