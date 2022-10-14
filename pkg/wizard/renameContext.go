package wizard

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/davidemaggi/kog/pkg/common"
	"github.com/davidemaggi/kog/pkg/k8s"
)

func RenameContext(configPath string, verbose bool) (err error) {

	ctxs, _, err := k8s.GetContexts(configPath, verbose)

	if !common.HandleError(err, "Error Getting Contexts", verbose) {
		return err

	}
	currentCtx, err := k8s.GetCurrentContext(configPath, verbose)
	renCtx := ""
	promptCtx := &survey.Select{
		Message: "Rename Context:",
		Options: ctxs,
		Default: currentCtx,
	}
	err = survey.AskOne(promptCtx, &renCtx)

	if !common.HandleError(err, "Error Asking Context", verbose) {
		return err

	}

	alias := ""
	prompt := &survey.Input{
		Message: "With:",
	}
	err = survey.AskOne(prompt, &alias)

	if !common.HandleError(err, "Error Getting the new name", verbose) {
		return err

	}

	k8s.RenameContext(configPath, renCtx, alias, verbose)

	return nil
}
