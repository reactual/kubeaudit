package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getAuditFunctions() []interface{} {
	return []interface{}{
		auditAllowPrivilegeEscalation, auditReadOnlyRootFS, auditRunAsNonRoot,
		auditAutomountServiceAccountToken, auditPrivileged, auditCapabilities,
	}
}

func runAllAudits(resources []Items) (issues map[Items][]Result) {
	for _, resource := range resources {
		var results []Result
		for _, function := range getAuditFunctions() {
			for _, result := range getResults([]Items{resource}, function) {
				results = append(results, result)
			}
		}
		issues[resource] = results
	}
	return issues
}

func (i *Items) fixPrivilegeEscalation() {
	switch kubeType := i.(type) {
	case Pod:
		containers = kubeType.Spec.Containers
	case default:
		kubeType.Spec.Template.Spec.Containers
	}
}

func fixStuff(resource Items, results []Result) (fixed Items) {
	for _, result := range results {
		for _, occurrence := range result.Occurrences {
			switch occurrence.id {
			case ErrorAllowPrivilegeEscalationNIL:
			case ErrorAllowPrivilegeEscalationTrue:
				fixPrivilegeEscalation(resource)
			case ErrorCapabilitiesAdded:
			case ErrorCapabilitiesNIL:
			case ErrorCapabilitiesNoneDropped:
			case ErrorCapabilitiesSomeDropped:
			case ErrorImageTagIncorrect:
			case ErrorImageTagMissing:
			case ErrorPrivilegedNIL:
			case ErrorPrivilegedTrue:
			case ErrorReadOnlyRootFilesystemFalse:
			case ErrorReadOnlyRootFilesystemNIL:
			case ErrorResourcesLimitsNIL:
			case ErrorResourcesLimitsCpuNIL:
			case ErrorResourcesLimitsCpuExceeded:
			case ErrorResourcesLimitsMemoryNIL:
			case ErrorResourcesLimitsMemoryExceeded:
			case ErrorRunAsNonRootFalse:
			case ErrorRunAsNonRootNIL:
			case ErrorSecurityContextNIL:
			case ErrorServiceAccountTokenDeprecated:
			case ErrorServiceAccountTokenNIL:
			case ErrorServiceAccountTokenNILAndNoName:
			case ErrorServiceAccountTokenNoName:
			case ErrorServiceAccountTokenTrueAndNoName:
			}
		}
	}
}

func autofix(*cobra.Command, []string) {
	items, err := getKubeResourcesManifest(rootConfig.manifest)
	if err != nil {
		log.Error(err)
	}
	issues := runAllAudits(items)
	var fixed []Items
	for resource, results := range issues {
		fixed = append(fixed, fixStuff(resource, results))
	}
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
