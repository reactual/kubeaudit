package cmd

import (
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type capsDropList struct {
	Drop []string `yaml:"capabilitiesToBeDropped"`
}

func recommendedCapabilitiesToBeDropped() (dropList []Capability, err error) {
	yamlFile, err := ioutil.ReadFile("config/capabilities-drop-list.yml")
	if err != nil {
		return
	}
	caps := capsDropList{}
	err = yaml.Unmarshal(yamlFile, &caps)
	if err != nil {
		return
	}
	for _, drop := range caps.Drop {
		dropList = append(dropList, Capability(drop))
	}
	return
}

func capsNotDropped(dropped []Capability) (notDropped []Capability, err error) {
	toBeDropped, err := recommendedCapabilitiesToBeDropped()
	if err != nil {
		return
	}
	for _, toBeDroppedCap := range toBeDropped {
		found := false
		for _, droppedCap := range dropped {
			if toBeDroppedCap == droppedCap {
				found = true
			}
		}
		if found == false {
			notDropped = append(notDropped, toBeDroppedCap)
		}
	}
	return
}

func checkCapabilities(container Container, result *Result) {
	if container.SecurityContext == nil {
		occ := Occurrence{
			id:      ErrorSecurityContextNIL,
			kind:    Error,
			message: "SecurityContext not set, please set it!",
		}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if container.SecurityContext.Capabilities == nil {
		occ := Occurrence{
			id:      ErrorCapabilitiesNIL,
			kind:    Error,
			message: "Capabilities field not defined!",
		}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if container.SecurityContext.Capabilities.Add != nil {
		for _, capAdded := range container.SecurityContext.Capabilities.Add {
			reason := result.Labels["kubeaudit.allow.capability."+strings.ToLower(string(capAdded))]
			if reason == "" {
				occ := Occurrence{
					id:       ErrorCapabilityAdded,
					kind:     Error,
					message:  "Capability added",
					metadata: Metadata{"CapName": string(capAdded)},
				}
				result.Occurrences = append(result.Occurrences, occ)
			} else {
				occ := Occurrence{
					id:      ErrorCapabilityAllowed,
					kind:    Warn,
					message: "Capability allowed",
					metadata: Metadata{
						"CapName": string(capAdded),
						"Reason":  prettifyReason(reason),
					},
				}
				result.Occurrences = append(result.Occurrences, occ)
			}
		}
	}

	capsNotDropped, err := capsNotDropped(container.SecurityContext.Capabilities.Drop)
	if err != nil {
		occ := Occurrence{
			id:      KubeauditInternalError,
			kind:    Error,
			message: "This should not have happened, if you are on kubeaudit master please consider to report: " + err.Error(),
		}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	for _, cap := range capsNotDropped {
		reason := result.Labels["kubeaudit.allow.capability."+strings.ToLower(string(cap))]
		if reason == "" {
			occ := Occurrence{
				id:       ErrorCapabilityNotDropped,
				kind:     Error,
				message:  "Capability not dropped",
				metadata: Metadata{"CapName": string(cap)},
			}
			result.Occurrences = append(result.Occurrences, occ)
		} else {
			occ := Occurrence{
				id:      ErrorCapabilityAllowed,
				kind:    Warn,
				message: "Capability allowed",
				metadata: Metadata{
					"CapName": string(cap),
					"Reason":  prettifyReason(reason),
				},
			}
			result.Occurrences = append(result.Occurrences, occ)
		}
	}
}

func auditCapabilities(items Items) (results []Result) {
	for _, item := range items.Iter() {
		containers, result := containerIter(item)
		for _, container := range containers {
			checkCapabilities(container, result)
			if result != nil && len(result.Occurrences) > 0 {
				results = append(results, *result)
				break
			}
		}
	}
	return
}

var capabilitiesCmd = &cobra.Command{
	Use:   "caps",
	Short: "Audit container for capabilities",
	Run:   runAudit(auditCapabilities),
}

func init() {
	RootCmd.AddCommand(capabilitiesCmd)
}
