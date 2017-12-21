package cmd

import (
	"github.com/spf13/cobra"

	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

func checkPrivileged(container Container, result *Result) {
	if container.SecurityContext == nil {
		occ := Occurrence{id: ErrorSecurityContextNIL, kind: Error, message: "SecurityContext not set, please set it!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if container.SecurityContext.Privileged == nil {
		// TODO find out what this exactly means
		occ := Occurrence{id: ErrorPrivilegedNIL, kind: Warn, message: "Privileged defaults to false, which results in non privileged, which is okay."}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
	if *container.SecurityContext.Privileged == true {
		occ := Occurrence{id: ErrorPrivilegedTrue, kind: Error, message: "Privileged set to true! Please change it to false!"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
}

func auditPrivileged(resources []k8sRuntime.Object) (results []Result) {
	for _, resource := range resources {
		for _, container := range getContainers(resource) {
			result := newResultFromResource(resource)
			checkPrivileged(container, &result)
			if len(result.Occurrences) > 0 {
				results = append(results, result)
				break
			}
		}
	}
	for _, result := range results {
		result.Print()
	}
	return
}

// runAsNonRootCmd represents the runAsNonRoot command
var privileged = &cobra.Command{
	Use:   "priv",
	Short: "Audit containers running as root",
	Long: `This command determines which containers in a kubernetes cluster
are running as privileged.

A PASS is given when a container runs in a non-privileged mode
A FAIL is generated when a container runs in a privileged mode

Example usage:
kubeaudit privileged`,
	Run: runAudit(auditPrivileged),
}

func init() {
	RootCmd.AddCommand(privileged)
}
