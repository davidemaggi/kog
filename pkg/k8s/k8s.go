package k8s

import (
	"errors"
	"fmt"
	"github.com/davidemaggi/kog/structs"
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

func GetPods(configPath string, verbose bool) (pods []string, err error) {
	if configPath == "" {
		configPath, err = FindKubeConfig()
	}

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return pods, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return pods, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return pods, err

	}

	var podsList, errNs = clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if errNs != nil {

		if verbose {
			log.Fatal(errNs)
		}
		return pods, errNs

	}
	for _, pod := range podsList.Items {
		pods = append(pods, pod.Name)
	}
	return pods, nil

}

func GetServices(configPath string, verbose bool) (services []string, err error) {
	if configPath == "" {
		configPath, err = FindKubeConfig()
	}

	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return services, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return services, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return services, err

	}

	var svcList, errNs = clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if errNs != nil {

		if verbose {
			log.Fatal(errNs)
		}
		return services, errNs

	}
	for _, svc := range svcList.Items {
		services = append(services, svc.Name)
	}
	return services, nil

}
