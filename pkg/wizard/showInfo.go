package wizard

import (
	"fmt"

	"github.com/davidemaggi/kog/pkg/k8s"
	"github.com/fatih/color"
)

func ShowInfo(configPath string, verbose bool, raw bool) (err error) {

	if !raw {

		currentCtx, _ := k8s.GetCurrentContext(configPath, verbose)

		currentNs, _ := k8s.GetCurrentNamespace(configPath, verbose)

		cyan := color.New(color.FgCyan).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		fmt.Printf("Current Context is %s\n", cyan(currentCtx))
		if currentNs != "" {
			fmt.Printf("Current NameSpace is %s\n", cyan(currentNs))

		} else {
			fmt.Printf("Current NameSpace is %s\n", red("not set"))

		}

	} else {
		k8s.PrintConfig(configPath)

	}

	return nil
}
