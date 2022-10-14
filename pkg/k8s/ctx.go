package k8s

import (
	"github.com/davidemaggi/kog/pkg/common"
	"k8s.io/client-go/tools/clientcmd/api"
)

func GetContexts(configPath string, verbose bool) (ctxs []string, rawCtxs map[string]*api.Context, err error) {
	ctxs = make([]string, 0)

	_, kubeConfig, err := GetConfig(configPath)

	if !common.HandleError(err, "Error Getting Config", verbose) {
		return ctxs, nil, err
	}

	for s := range kubeConfig.Contexts {
		ctxs = append(ctxs, s)
	}
	return ctxs, kubeConfig.Contexts, nil

}

func GetCurrentContext(configPath string, verbose bool) (ctx string, err error) {

	_, kubeConfig, err := GetConfig(configPath)

	if !common.HandleError(err, "Error Getting Config", verbose) {
		return "", err
	}
	return kubeConfig.CurrentContext, nil
}

func SetContext(configPath string, ctx string, verbose bool) (err error) {

	_, kubeConfig, err := GetConfig(configPath)

	if !common.HandleError(err, "Error Getting Config", verbose) {
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

func DeleteContext(configPath string, ctx string, force bool, verbose bool) (err error) {

	_, kubeConfig, err := GetConfig(configPath)

	if !common.HandleError(err, "Error Getting Config", verbose) {
		return err
	}

	if ctx == kubeConfig.CurrentContext {

		//You are deleting the current default
		kubeConfig.CurrentContext = ""

	}

	toRemove := kubeConfig.Contexts[ctx]

	delete(kubeConfig.Contexts, ctx)

	if force {
		delete(kubeConfig.AuthInfos, toRemove.AuthInfo)
		delete(kubeConfig.Clusters, toRemove.Cluster)

	}

	if configPath == "" {
		configPath, _ = FindKubeConfig()
	}

	s, err := SaveConfig(configPath, kubeConfig, true)

	_ = s
	return nil

}

func RenameContext(configPath string, ctx string, newName string, verbose bool) (err error) {

	_, kubeConfig, err := GetConfig(configPath)

	if !common.HandleError(err, "Error Getting Config", verbose) {
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
