package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	k8sRuntime "k8s.io/apimachinery/pkg/runtime"
)

var imgConfig imgFlags

type imgFlags struct {
	img  string
	name string
	tag  string
}

func (image *imgFlags) splitImageString() {
	tokens := strings.Split(image.img, ":")
	if len(tokens) > 0 {
		image.name = tokens[0]
	}
	if len(tokens) > 1 {
		image.tag = tokens[1]
	}
}

func checkImage(container Container, image imgFlags, result *Result) {
	image.splitImageString()
	contImage := imgFlags{img: container.Image}
	contImage.splitImageString()
	result.ImageName = contImage.name
	result.ImageTag = contImage.tag

	if len(contImage.tag) == 0 {
		occ := Occurrence{id: ErrorImageTagMissing, kind: Warn, message: "Image tag was missing"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if contImage.name == image.name && contImage.tag != image.tag {
		occ := Occurrence{id: ErrorImageTagIncorrect, kind: Error, message: "Image tag was incorrect"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}

	if contImage.name == image.name && contImage.tag == image.tag {
		occ := Occurrence{id: InfoImageCorrect, kind: Info, message: "Image tag was correct"}
		result.Occurrences = append(result.Occurrences, occ)
		return
	}
}

func auditImages(image imgFlags, resource k8sRuntime.Object) (results []Result) {
	for _, container := range getContainers(resource) {
		result := newResultFromResource(resource)
		checkImage(container, image, &result)
		if len(result.Occurrences) > 0 {
			results = append(results, result)
			break
		}
	}
	for _, result := range results {
		result.Print()
	}
	return
}

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Audit container images",
	Long: `This command audits a container against a given image:tag
An INFO log is given when a container has a matching image:tag
An ERROR log is generated when a container does not match the image:tag
This command is also a root command, check kubeaudit sc --help
Example usage:
kubeaudit image --image gcr.io/google_containers/echoserver:1.7
kubeaudit image -i gcr.io/google_containers/echoserver:1.7`,
	Run: runAudit(auditImages),
}

func init() {
	RootCmd.AddCommand(imageCmd)
	imageCmd.Flags().StringVarP(&imgConfig.img, "image", "i", "", "image to check against")
}
