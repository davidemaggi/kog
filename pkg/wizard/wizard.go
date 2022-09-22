package wizard

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/davidemaggi/kog/pkg/k8s"
	"github.com/fatih/color"
	"log"
)

func SelectContext(configPath string, verbose bool) (err error) {

	ctxs, rawCtxs, err := k8s.GetContexts(configPath, verbose)

	if err != nil {

		log.Fatal("Error getting Contexts")
		return
	}
	currentCtx, err := k8s.GetCurrentContext(configPath, verbose)
	newCtx := ""
	promptCtx := &survey.Select{
		Message:  "Select Context:",
		Options:  ctxs,
		Default:  currentCtx,
		PageSize: 20,
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
		Message:  "Select Namespace:",
		Options:  namespaces,
		PageSize: 20,
	}
	survey.AskOne(promptNs, &newNs)

	k8s.SetNameSpace(newNs, verbose)
	_ = rawCtxs
	return nil
}
func DeleteContext(configPath string, verbose bool) (err error) {

	ctxs, _, err := k8s.GetContexts(configPath, verbose)

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
func ShowInfo(configPath string, verbose bool) (err error) {

	_, rawCtxs, err := k8s.GetContexts(configPath, verbose)

	currentCtx, err := k8s.GetCurrentContext(configPath, verbose)

	cur := rawCtxs[currentCtx]

	if cur != nil {
		cyan := color.New(color.FgCyan).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		fmt.Printf("Current Context is %s.\n", cyan(currentCtx))
		if cur.Namespace != "" {
			fmt.Printf("Current NameSpace is %s.\n", cyan(cur.Namespace))

		} else {
			fmt.Printf("Current NameSpace is %s.\n", red("not set"))

		}

	}

	return nil
}
