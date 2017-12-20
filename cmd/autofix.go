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

func runAllAudits(resources []Items) (fixedResources Items) {
	log.Error(len(resources))
	for _, resource := range resources {
		for _, function := range getAuditFunctions() {
			results := getResults([]Items{resource}, function)
			if len(results) > 0 {
				log.Error(results)
				log.Error(resource)
				fixStuff(resource, results)
			}
		}
	}
	return fixedResources
}

func fixStuff(resource Items, results []Result) {
	fmt.Println("fixing stuff")
}

func autofix(*cobra.Command, []string) {
	items, err := getKubeResourcesManifest(rootConfig.manifest)
	if err != nil {
		log.Error(err)
	}
	runAllAudits(items)
	//err = writeManifestFile(rootConfig.manifest, decoded)
	//if err != nil {
	//	return
	//}
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
