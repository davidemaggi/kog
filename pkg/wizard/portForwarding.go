package wizard

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/davidemaggi/kog/pkg/common"
	"github.com/davidemaggi/kog/pkg/k8s"
	v1 "k8s.io/api/core/v1"
)

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

	if !common.HandleError(err, "Error Selecting Context", verbose) {
		return err
	}
	var toForward = []string{}
	var pods []v1.Pod
	var services []v1.Service

	if strings.ToLower(fwdWhat) == "pod" {
		toForward, pods, err = k8s.GetPods(configPath, verbose)
	}
	if strings.ToLower(fwdWhat) == "service" {
		toForward, services, err = k8s.GetServices(configPath, verbose)

	}

	if !common.HandleError(err, "Error getting entities", verbose) {
		return err
	}

	fwdEntityt := ""
	promptEntity := &survey.Select{
		Message: "Entity to Forward:",
		Options: toForward,
	}
	err = survey.AskOne(promptEntity, &fwdEntityt)

	if !common.HandleError(err, "Error Selecting Entity to forward", verbose) {
		return err
	}

	if strings.ToLower(fwdWhat) == "pod" {
		for i := range pods {
			if pods[i].Name == fwdEntityt {

				fwdPort := "0"

				var ports = []string{}

				for ci := range pods[i].Spec.Containers {
					for pi := range pods[i].Spec.Containers[ci].Ports {
						ports = append(ports, fmt.Sprintf("%s: %d",
							pods[i].Spec.Containers[ci].Ports[pi].Protocol, pods[i].Spec.Containers[ci].Ports[pi].ContainerPort))
					}
				}

				promptPort := &survey.Select{
					Message: "Port to forward:",
					Options: ports,
				}
				err = survey.AskOne(promptPort, &fwdPort)

				localport := "0"
				promptlocalPort := &survey.Input{
					Message: "On Local port",
				}
				survey.AskOne(promptlocalPort, &localport)

				portInt, err := strconv.Atoi(fwdPort[5:])
				if !common.HandleError(err, "Error converting Port", verbose) {
					return err
				}
				localportInt, err := strconv.Atoi(localport)
				if !common.HandleError(err, "Error converting Port", verbose) {
					return err
				}

				k8s.PortForwardPod(configPath, &pods[i], int32(portInt), int32(localportInt), false)
			}
		}
	}
	if strings.ToLower(fwdWhat) == "service" {
		for i := range services {
			if services[i].Name == fwdEntityt {

				fwdPort := "0"

				var ports = []string{}

				for pi := range services[i].Spec.Ports {

					ports = append(ports, fmt.Sprintf("%s: %d",
						services[i].Spec.Ports[pi].Protocol, services[i].Spec.Ports[pi].Port))

				}

				promptPort := &survey.Select{
					Message: "Port to forward:",
					Options: ports,
				}
				err = survey.AskOne(promptPort, &fwdPort)

				localport := "0"
				promptlocalPort := &survey.Input{
					Message: "On Local port",
				}
				survey.AskOne(promptlocalPort, &localport)
				portInt, err := strconv.Atoi(fwdPort[5:])
				if !common.HandleError(err, "Error converting Port", verbose) {
					return err
				}
				localportInt, err := strconv.Atoi(localport)
				if !common.HandleError(err, "Error converting Port", verbose) {
					return err
				}

				for pi := range services[i].Spec.Ports {

					if services[i].Spec.Ports[pi].Port == int32(portInt) {
						k8s.PortForwardSvc(
							configPath, &services[i], services[i].Spec.Ports[pi].TargetPort.IntVal, int32(localportInt), false)
					}
				}

			}
		}
	}
	_ = toForward

	return nil
}
