package wizard

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/davidemaggi/kog/pkg/k8s"
	"log"
)

func SelectContext(configPath string, verbose bool) (err error) {

	ctxs, err := k8s.GetContexts(configPath, verbose)

	if err != nil {

		log.Fatal("Error getting Contexts")
		return
	}
	currentCtx, err := k8s.GetCurrentContext(configPath, verbose)
	newCtx := ""
	promptCtx := &survey.Select{
		Message: "Select Context:",
		Options: ctxs,
		Default: currentCtx,
	}
	survey.AskOne(promptCtx, &newCtx)

	k8s.SetContext(configPath, newCtx, verbose)

	namespaces, err := k8s.GetNamespaces(configPath, verbose)

	if err != nil {

		log.Fatal("Error getting Namespaces")
		return err
	}
	newNs := ""
	promptNs := &survey.Select{
		Message: "Select Namespace:",
		Options: namespaces,
	}
	survey.AskOne(promptNs, &newNs)

	k8s.SetNameSpace(newNs, verbose)

	return nil
}
func DeleteContext(configPath string, verbose bool) (err error) {

	ctxs, err := k8s.GetContexts(configPath, verbose)

	if err != nil {

		log.Fatal("Error getting Contexts")
		return
	}
	currentCtx, err := k8s.GetCurrentContext(configPath, verbose)
	delCtx := ""
	promptCtx := &survey.Select{
		Message: "Delete Context:",
		Options: ctxs,
		Default: currentCtx,
	}
	survey.AskOne(promptCtx, &delCtx)

	k8s.DeleteContext(configPath, delCtx, verbose)

	return nil
}
