package wizard

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/davidemaggi/kog/pkg/k8s"
	"github.com/fatih/color"
	v1 "k8s.io/api/core/v1"
	"log"
	"strings"
)

func SelectContext(configPath string, verbose bool) (err error) {

	ctxs, rawCtxs, err := k8s.GetContexts(configPath, verbose)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		log.Fatal("Error getting contexts")

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

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		log.Fatal("Error getting context")

		return err
	}

	k8s.SetContext(configPath, newCtx, verbose)

	namespaces, err := k8s.GetNamespaces(configPath, verbose)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		log.Fatal("Error getting namespaces")

		return err
	}
	newNs := ""
	promptNs := &survey.Select{
		Message:  "Select Namespace:",
		Options:  namespaces,
		PageSize: 20,
	}
	err = survey.AskOne(promptNs, &newNs)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		log.Fatal("Error getting namespace")

		return err
	}

	k8s.SetNameSpace(newNs, verbose)
	_ = rawCtxs
	return nil
}
func DeleteContext(configPath string, verbose bool) (err error) {

	ctxs, _, err := k8s.GetContexts(configPath, verbose)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		log.Fatal("Error getting contexts")

		return err
	}
	currentCtx, err := k8s.GetCurrentContext(configPath, verbose)
	delCtx := ""
	promptCtx := &survey.Select{
		Message: "Delete Context:",
		Options: ctxs,
		Default: currentCtx,
	}
	err = survey.AskOne(promptCtx, &delCtx)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		log.Fatal("Error getting context to delete")

		return err
	}

	k8s.DeleteContext(configPath, delCtx, verbose)

	return nil
}
func ShowInfo(configPath string, verbose bool, raw bool) (err error) {

	if !raw {
		_, rawCtxs, _ := k8s.GetContexts(configPath, verbose)

		currentCtx, _ := k8s.GetCurrentContext(configPath, verbose)

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
	} else {
		k8s.PrintConfig(configPath)

	}

	return nil
}

func RenameContext(configPath string, verbose bool) (err error) {

	ctxs, _, err := k8s.GetContexts(configPath, verbose)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		log.Fatal("Error getting Contexts")
		return
	}
	currentCtx, err := k8s.GetCurrentContext(configPath, verbose)
	renCtx := ""
	promptCtx := &survey.Select{
		Message: "Rename Context:",
		Options: ctxs,
		Default: currentCtx,
	}
	err = survey.AskOne(promptCtx, &renCtx)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		log.Fatal("Error getting context to rename")

		return err
	}

	alias := ""
	prompt := &survey.Input{
		Message: "With:",
	}
	err = survey.AskOne(prompt, &alias)

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		log.Fatal("Error getting alias")

		return err
	}

	k8s.RenameContext(configPath, renCtx, alias, verbose)

	return nil
}

func PortForwarding(configPath string, verbose bool) (err error) {

	selRes := SelectContext(configPath, verbose)

	if selRes != nil {
		return selRes
	}

	fwdWhat := ""
	promptWhat := &survey.Select{
		Message: "Entity type to Forward:",
		Options: []string{"Service", "Pod"},
		Default: "Service",
	}
	err = survey.AskOne(promptWhat, &fwdWhat)

	if err != nil {
		return err
	}
	var toForward = []string{}
	var pods []v1.Pod
	var services []v1.Service

	if strings.ToLower(fwdWhat) == "pod" {
		toForward, pods, _ = k8s.GetPods(configPath, verbose)
	}
	if strings.ToLower(fwdWhat) == "service" {
		toForward, services, _ = k8s.GetServices(configPath, verbose)

	}
	fwdEntityt := ""
	promptEntity := &survey.Select{
		Message: "Entity to Forward:",
		Options: toForward,
	}
	err = survey.AskOne(promptEntity, &fwdEntityt)

	if strings.ToLower(fwdWhat) == "pod" {
		for i := range pods {
			if pods[i].Name == fwdEntityt {

				k8s.PortForwardPod(configPath, pods[i], 80, 8080, false)
			}
		}
	}
	if strings.ToLower(fwdWhat) == "service" {
		for i := range services {
			if services[i].Name == fwdEntityt {
				k8s.PortForwardSvc(configPath, services[i], 80, 8080, false)
			}
		}
	}
	_ = toForward

	return nil
}
