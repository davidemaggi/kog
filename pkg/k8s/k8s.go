package k8s

import (
	"errors"
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

func GetConfig(kubeConfigPath string) (kubeConfig *api.Config, err error) {

	if kubeConfigPath == "" {
		kubeConfigPath, err = FindKubeConfig()
		if err != nil {
			log.Fatal(err)
		}
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
		oldConfig, _ := GetConfig("")
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

func GetContexts(configPath string, verbose bool) (ctxs []string, err error) {
	ctxs = make([]string, 0)

	kubeConfig, err := GetConfig(configPath)

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

func GetCurrentContext(configPath string, verbose bool) (ctx string, err error) {

	kubeConfig, err := GetConfig(configPath)
	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return "", err
	}

	return kubeConfig.CurrentContext, nil
}

func SetContext(configPath string, ctx string, verbose bool) (err error) {

	kubeConfig, err := GetConfig(configPath)

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

func SetNameSpace(ns string, verbose bool) (err error) {

	kubeConfig, err := GetConfig("")

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
	oldConfig, err := GetConfig(oldFile)

	result.From = newFile
	result.To = oldFile

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		return nil, result, err

	}
	origConfig := oldConfig.DeepCopy()
	newConfig, err := GetConfig(newFile)

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
	//endregion

	//region Extensions
	/*
		for ext_id, ext := range newConfig.Extensions {
			// index is the index where we are
			// element is the element from someSlice for where we are

			exists := oldConfig.Extensions[ext_id] != nil

			if exists {

				if force {
					oldConfig.Extensions[ext_id] = ext.DeepCopyObject()
				}

			} else {

				oldConfig.Extensions[ext_id] = ext.DeepCopyObject()
			}

		}

	*/
	//endregion

	//region Extensions 2
	/*
		for ext2_id, ext2 := range newConfig.Preferences.Extensions {
			// index is the index where we are
			// element is the element from someSlice for where we are

			exists := oldConfig.Preferences.Extensions[ext2_id] != nil

			if exists {

				if force {
					oldConfig.Preferences.Extensions[ext2_id] = ext2.DeepCopyObject()
				}

			} else {

				oldConfig.Preferences.Extensions[ext2_id] = ext2.DeepCopyObject()
			}

		}
	*/

	//endregion

	_ = origConfig
	result.IsOk = true

	return oldConfig, result, nil
}
