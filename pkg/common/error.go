package common

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

func HandleError(err error, msg string, verbose bool) (ret bool) {

	red := color.New(color.FgRed).SprintFunc()

	if err != nil {
		if verbose {
			log.Fatal(err)
		}
		fmt.Printf("An Error has occured \"%s\": %s \n", red(msg), red(err))

		return false
	}

	return true
}
