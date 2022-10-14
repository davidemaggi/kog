package common

import "log"

func HandleError(err error, msg string, verbose bool) (ret bool) {
	if err != nil {
		if verbose {
			log.Fatal(err)
		}

		return true
	}

	return false
}
