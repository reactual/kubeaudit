package cmd

import (
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
)

var path = "../fixtures/"

func runTest(t *testing.T, file string, function interface{}, errCode int, argStr ...string) (results []Result) {
	assert := assert.New(t)
	file = filepath.Join(path, file)
	var image imgFlags
	var limits limitFlags
	switch function.(type) {
	case (func(imgFlags, Items) []Result):
		if len(argStr) != 1 {
			log.Fatal("Incorrect number of images specified")
		}
		image = imgFlags{img: argStr[0]}
		image.splitImageString()
	case (func(limitFlags, Items) []Result):
		if len(argStr) == 2 {
			limits = limitFlags{cpuArg: argStr[0], memoryArg: argStr[1]}
		} else if len(argStr) == 0 {
			limits = limitFlags{cpuArg: "", memoryArg: ""}
		} else {
			log.Fatal("Incorrect number of images specified")
		}
		limits.parseLimitFlags()
	}

	resources, err := getKubeResourcesManifest(file)
	assert.Nil(err)

	for _, resource := range resources {
		var currentResults []Result
		switch f := function.(type) {
		case (func(Items) []Result):
			currentResults = f(resource)
		case (func(imgFlags, Items) []Result):
			currentResults = f(image, resource)
		case (func(limitFlags, Items) []Result):
			currentResults = f(limits, resource)
		default:
			log.Fatal("Invalid function provided")
		}
		for _, currentResult := range currentResults {
			results = append(results, currentResult)
		}
	}
	var errors []int
	for _, result := range results {
		for _, occurrence := range result.Occurrences {
			errors = append(errors, occurrence.id)
		}
	}

	if errCode != 0 {
		assert.Contains(errors, errCode)
	}
	return
}

func runTestInNamespace(t *testing.T, namespace string, file string, function interface{}, errCode ...int) {
	rootConfig.namespace = namespace
	var results = runTest(t, file, function, 0)
	var errors []int
	for _, result := range results {
		for _, occurrence := range result.Occurrences {
			errors = append(errors, occurrence.id)
			assert.Contains(t, errCode, occurrence.id)
		}
	}
	assert.Equal(t, len(errCode), len(errors))
	rootConfig.namespace = apiv1.NamespaceAll
}
