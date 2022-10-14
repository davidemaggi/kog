package k8s

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/davidemaggi/kog/pkg/common"
	"github.com/davidemaggi/kog/structs"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

func GetConfig(kubeConfigPath string) (path string, kubeConfig *api.Config, err error) {

	if kubeConfigPath == "" {
		kubeConfigPath, err = FindKubeConfig()
		if !common.HandleError(err, "Cannot find hube config path", true) {
			return kubeConfigPath, nil, err

		}
	}

	kubeConfig, err = clientcmd.LoadFromFile(kubeConfigPath)

	if !common.HandleError(err, "Error Loading KubeConfig", true) {
		return kubeConfigPath, nil, err

	}

	return kubeConfigPath, kubeConfig, nil

}

func PrintConfig(kubeConfigPath string) (err error) {

	if kubeConfigPath == "" {
		kubeConfigPath, err = FindKubeConfig()
		if !common.HandleError(err, "Error retrieving Kube Config", true) {
			return nil

		}
	}

	b, err := os.ReadFile(kubeConfigPath)
	if err != nil {
		fmt.Print(err)
	}

	str := string(b) // convert content to a 'string'

	fmt.Println(str) // print the content as a 'string'

	return nil

}

func SaveConfig(filePath string, kubeConfig *api.Config, backup bool) (msg string, err error) {

	if backup {
		_, oldConfig, _ := GetConfig("")
		err = clientcmd.WriteToFile(*oldConfig, filePath+".backup")

		if !common.HandleError(err, "Error during backup", true) {
			return "Error during backup", err

		}
	}

	err = clientcmd.WriteToFile(*kubeConfig, filePath)

	if !common.HandleError(err, "Error Saving Backup", true) {
		return "Error Saving Backup", err

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

func MergeConfigs(newFile string, oldFile string, force bool, verbose bool) (kubeConfig *api.Config, result structs.MergeResult, err error) {

	result = structs.New_MergeResult()
	oldPath, oldConfig, err := GetConfig(oldFile)

	result.From = newFile
	result.To = oldPath

	if !common.HandleError(err, "Error Retrieving kube config", verbose) {
		return nil, result, err

	}
	origConfig := oldConfig.DeepCopy()
	_, newConfig, err := GetConfig(newFile)

	if !common.HandleError(err, "Error Loading kube config", verbose) {
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

func buildConfigFromPath(configPath string, verbose bool) (config *rest.Config, err error) {

	if configPath == "" {
		configPath, err = FindKubeConfig()
	}
	if !common.HandleError(err, "Error Retrieving kube config", verbose) {
		return nil, err
	}

	if configPath == "" {
		configPath, err = FindKubeConfig()
	}

	if !common.HandleError(err, "Error Loafing kube config", verbose) {
		return nil, err
	}
	config, err = clientcmd.BuildConfigFromFlags("", configPath)
	if !common.HandleError(err, "Error Building kube config", verbose) {
		return nil, err
	}

	return config, err

}
