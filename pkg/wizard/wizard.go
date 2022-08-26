package wizard

import (
	"fmt"
	"github.com/davidemaggi/kog/pkg/k8s"
	"github.com/manifoldco/promptui"
)

func SelectContext() (result string) {

	var ctxs = make([]string, 0)

	var config, err = k8s.GetConfig()
	for s, _ := range config.Contexts {
		ctxs = append(ctxs, s)
	}

	prompt := promptui.Select{
		Label: "Select Context",
		Items: ctxs,
	}

	_, result, err = prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)

	return result
}
