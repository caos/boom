package helper

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	firstResource = &Resource{
		Kind:       "test1",
		ApiVersion: "api/v1",
		Metadata: &Metadata{
			Name:      "test1",
			Namespace: "test",
		},
	}
	secondResource = &Resource{
		Kind:       "test2",
		ApiVersion: "api/v1",
		Metadata: &Metadata{
			Name:      "test2",
			Namespace: "test",
		},
	}
)

type metadataTest struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}
type resourceTest struct {
	Kind       string        `yaml:"kind"`
	ApiVersion string        `yaml:"apiVersion"`
	Metadata   *metadataTest `yaml:"metadata"`
}

func TestHelper_AddStructToYaml(t *testing.T) {
	root := "/tmp/nonexistent"
	err := os.MkdirAll(root, os.ModePerm)
	assert.NoError(t, err)

	path := "/tmp/nonexistent/test.yaml"
	err = AddStructToYaml(path, firstResource)
	assert.NoError(t, err)

	files := getFiles(root)
	assert.Len(t, files, 1)

	var restTest Resource
	err = YamlToStruct(path, &restTest)
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(firstResource, &restTest))

	err = os.RemoveAll(root)
	assert.NoError(t, err)
}

func TestHelper_AddStringToYaml(t *testing.T) {
	root := "/tmp/nonexistent"
	err := os.MkdirAll(root, os.ModePerm)
	assert.NoError(t, err)

	path := "/tmp/nonexistent/test.yaml"
	err = AddStringToYaml(path, "test")
	assert.NoError(t, err)

	files := getFiles(root)
	assert.Len(t, files, 1)

	restTest, err := YamlToString(path)
	assert.NoError(t, err)

	assert.Equal(t, "test", restTest)

	err = os.RemoveAll(root)
	assert.NoError(t, err)
}

func TestHelper_AddStringObjectToYaml(t *testing.T) {
	root := "/tmp/nonexistent"
	err := os.MkdirAll(root, os.ModePerm)
	assert.NoError(t, err)

	path := "/tmp/nonexistent/test.yaml"
	err = AddStringObjectToYaml(path, "test: test")
	assert.NoError(t, err)

	files := getFiles(root)
	assert.Len(t, files, 1)

	restTest, err := YamlToString(path)
	assert.NoError(t, err)
	assert.Equal(t, "test: test", restTest)

	err = os.RemoveAll(root)
	assert.NoError(t, err)
}

func TestHelper_AddStringBeforePointForKindAndName(t *testing.T) {
	root := "/tmp/nonexistent"
	err := os.MkdirAll(root, os.ModePerm)
	assert.NoError(t, err)

	path := "/tmp/nonexistent/test.yaml"
	err = AddStructToYaml(path, firstResource)
	assert.NoError(t, err)

	firstResourceTest := &resourceTest{
		Kind:       firstResource.Kind,
		ApiVersion: firstResource.ApiVersion,
		Metadata: &metadataTest{
			Name:      firstResource.Metadata.Name,
			Namespace: "test",
		},
	}

	files := getFiles(root)
	assert.Len(t, files, 1)

	var restTest resourceTest
	err = YamlToStruct(path, &restTest)
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(firstResourceTest, &restTest))

	err = os.RemoveAll(root)
	assert.NoError(t, err)
}
