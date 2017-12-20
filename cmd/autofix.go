package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getAuditFunctions() []interface{} {
	return []interface{}{
		auditAllowPrivilegeEscalation, auditReadOnlyRootFS, auditRunAsNonRoot,
		auditAutomountServiceAccountToken, auditPrivileged, auditCapabilities,
	}
}

func runAllAudits() {
	resources, err := getResources()
	if err != nil {
		log.Error(err)
		return
	}

	for _, function := range getAuditFunctions() {
		for _, result := range getResults(resources, function) {
			fixStuff(result)
		}
	}
}

func fixStuff(result Result) {
	fmt.Println("fixing stuff")
}

func autofix(*cobra.Command, []string) {
	decoded, err := readManifestFile(rootConfig.manifest)
	if err != nil {
		return
	}
	for _, resource := range decoded {
		fmt.Println(resource)
	}
	err = writeManifestFile(rootConfig.manifest, decoded)
	if err != nil {
		return
	}
}

var autofixCmd = &cobra.Command{
	Use:   "autofix",
	Short: "Automagically fixes a manifest to be secure",
	Long: `"autofix" will examine a manifest file and automagically fill in the
blanks for you leave your yaml file more secure than it found it.

Example usage:
kubeaudit autofix -f /path/to/yaml`,
	Run: autofix,
}

func init() {
	RootCmd.AddCommand(autofixCmd)
}
