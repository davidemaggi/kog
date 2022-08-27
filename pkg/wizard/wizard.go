package wizard

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/davidemaggi/kog/pkg/k8s"
	"log"
)

func SelectContext() {

	ctxs, err := k8s.GetContexts()

	if err != nil {
		log.Fatal("Error getting Contexts")
		return
	}
	currentCtx, err := k8s.GetCurrentContext()
	newCtx := ""
	promptCtx := &survey.Select{
		Message: "Select Context:",
		Options: ctxs,
		Default: currentCtx,
	}
	survey.AskOne(promptCtx, &newCtx)

	k8s.SetContexts(newCtx)

	namespaces, err := k8s.GetNamespaces()

	newNs := ""
	promptNs := &survey.Select{
		Message: "Select Namespace:",
		Options: namespaces,
	}
	survey.AskOne(promptNs, &newNs)
	_ = newNs
}
