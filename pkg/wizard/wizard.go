package wizard

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/davidemaggi/kog/pkg/k8s"
	"log"
)

func SelectContext(verbose bool) (err error) {

	ctxs, err := k8s.GetContexts(verbose)

	if err != nil {

		log.Fatal("Error getting Contexts")
		return
	}
	currentCtx, err := k8s.GetCurrentContext(verbose)
	newCtx := ""
	promptCtx := &survey.Select{
		Message: "Select Context:",
		Options: ctxs,
		Default: currentCtx,
	}
	survey.AskOne(promptCtx, &newCtx)

	k8s.SetContext(newCtx, verbose)

	namespaces, err := k8s.GetNamespaces(verbose)

	if err != nil {

		log.Fatal("Error getting Namespaces")
		return err
	}
	newNs := ""
	promptNs := &survey.Select{
		Message: "Select Namespace:",
		Options: append(namespaces, "dEfAulT3"),
	}
	survey.AskOne(promptNs, &newNs)

	k8s.SetNameSpace(newNs, verbose)

	return nil
}
