package wizard

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/davidemaggi/kog/pkg/common"
	"github.com/davidemaggi/kog/pkg/k8s"
)

func DeleteContext(configPath string, verbose bool) (err error) {

	ctxs, _, err := k8s.GetContexts(configPath, verbose)

	if !common.HandleError(err, "Error Getting Contexts", verbose) {
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

	if !common.HandleError(err, "Error Asking Context", verbose) {
		return err

	}

	k8s.DeleteContext(configPath, delCtx, verbose)

	return nil
}
