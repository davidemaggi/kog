package wizard

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/davidemaggi/kog/pkg/common"
	"github.com/davidemaggi/kog/pkg/k8s"
)

func SelectContext(configPath string, verbose bool) (err error) {

	ctxs, rawCtxs, err := k8s.GetContexts(configPath, verbose)

	if !common.HandleError(err, "Error Getting Contexts", verbose) {
		return err

	}

	currentCtx, err := k8s.GetCurrentContext(configPath, verbose)
	newCtx := ""
	promptCtx := &survey.Select{
		Message:  "Select Context:",
		Options:  ctxs,
		Default:  currentCtx,
		PageSize: 20,
	}
	err = survey.AskOne(promptCtx, &newCtx)

	if !common.HandleError(err, "Error Asking Context", verbose) {
		return err

	}

	k8s.SetContext(configPath, newCtx, verbose)

	namespaces, err := k8s.GetNamespaces(configPath, verbose)

	if !common.HandleError(err, "Error Getting Namespaces", verbose) {
		return err

	}
	newNs := ""
	promptNs := &survey.Select{
		Message:  "Select Namespace:",
		Options:  namespaces,
		PageSize: 20,
	}
	err = survey.AskOne(promptNs, &newNs)

	if !common.HandleError(err, "Error Asking Namespace Config", verbose) {
		return err

	}

	k8s.SetNameSpace(newNs, verbose)
	_ = rawCtxs
	return nil
}
