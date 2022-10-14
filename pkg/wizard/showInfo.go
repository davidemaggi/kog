package wizard

import (
	"fmt"

	"github.com/davidemaggi/kog/pkg/k8s"
	"github.com/fatih/color"
)

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
