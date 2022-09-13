package merge

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/davidemaggi/kog/pkg/k8s"
	"github.com/davidemaggi/kog/structs"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"strconv"
)

func MergeConfig(path2Merge string, ConfigPath string, Force bool, Verbose bool) {

	newConf, mergeResult, err := k8s.MergeConfigs(path2Merge, ConfigPath, Force, Verbose)

	if err != nil {
		if Verbose {
			log.Fatal(err)
		}
	}
	if !mergeResult.IsOk {
		fmt.Println(mergeResult.Msg)
	} else {

		if mergeResult.DoneSomething {

			printResult(mergeResult)

			confirmed := false
			promptConf := &survey.Confirm{
				Message: "Is this Ok?",
			}
			survey.AskOne(promptConf, &confirmed)

			if confirmed {
				if ConfigPath == "" {

					tmpPath, _ := k8s.FindKubeConfig()

					k8s.SaveConfig(tmpPath, newConf, true)

				} else {
					k8s.SaveConfig(ConfigPath, newConf, true)

				}
			} else {
				fmt.Println("That's Ok, we are still friends ❤  ")

			}
		} else {
			fmt.Println("Nothing to do here")

		}
	}

}

func printResult(result structs.MergeResult) {
	data := [][]string{}
	totAdd := 0
	totMod := 0
	totRem := 0

	for _, detail := range result.Details {

		str := ""

		switch detail.ObjType {
		case structs.Context:
			str = "Context"
			break
		case structs.User:
			str = "User"
			break
		case structs.Cluster:
			str = "Cluster"
			break
		}
		totAdd += detail.Added
		totMod += detail.Modified
		totRem += detail.Removed
		data = append(data, []string{str, strconv.Itoa(detail.Added), strconv.Itoa(detail.Modified), strconv.Itoa(detail.Removed)})

	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Item", "Added", "Modified", "Removed"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgHiYellowColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.FgRedColor, tablewriter.Bold})
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgHiYellowColor},
		tablewriter.Colors{tablewriter.FgRedColor})
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})

	//table.SetFooterColor(
	//	tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
	//	tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.Bold},
	//	tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold},
	//	tablewriter.Colors{tablewriter.FgRedColor, tablewriter.Bold})

	for _, v := range data {
		table.Append(v)
	}
	data = append(data, []string{"", "", "", ""})
	//table.SetFooter([]string{"Total", strconv.Itoa(totAdd), strconv.Itoa(totMod), strconv.Itoa(totRem)}) // Add Footer
	fmt.Println(result.From + " --> " + result.To)
	table.Render() // Send output

}
