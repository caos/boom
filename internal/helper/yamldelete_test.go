package helper

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelper_DeleteKindFromYaml_first(t *testing.T) {
	root := "/tmp/nonexistent"
	err := os.MkdirAll(root, os.ModePerm)
	assert.NoError(t, err)

	path := "/tmp/nonexistent/test.yaml"
	err = AddStructToYaml(path, firstResource)
	assert.NoError(t, err)
	err = AddStructToYaml(path, secondResource)
	assert.NoError(t, err)

	files := getFiles(root)
	assert.Len(t, files, 1)

	err = DeleteKindFromYaml(path, "test1")
	assert.NoError(t, err)

	var restTest Resource
	err = YamlToStruct(path, &restTest)
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(secondResource, &restTest))

	err = os.RemoveAll(root)
	assert.NoError(t, err)
}

func TestHelper_DeleteKindFromYaml_second(t *testing.T) {
	root := "/tmp/nonexistent"
	err := os.MkdirAll(root, os.ModePerm)
	assert.NoError(t, err)

	path := "/tmp/nonexistent/test.yaml"
	err = AddStructToYaml(path, firstResource)
	assert.NoError(t, err)
	err = AddStructToYaml(path, secondResource)
	assert.NoError(t, err)

	files := getFiles(root)
	assert.Len(t, files, 1)

	err = DeleteKindFromYaml(path, "test2")
	assert.NoError(t, err)

	var restTest Resource
	err = YamlToStruct(path, &restTest)
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(firstResource, &restTest))

	err = os.RemoveAll(root)
	assert.NoError(t, err)
}

func TestHelper_DeleteKindFromYaml_both(t *testing.T) {
	root := "/tmp/nonexistent"
	err := os.MkdirAll(root, os.ModePerm)
	assert.NoError(t, err)

	path := "/tmp/nonexistent/test.yaml"
	err = AddStructToYaml(path, firstResource)
	assert.NoError(t, err)
	err = AddStructToYaml(path, secondResource)
	assert.NoError(t, err)

	files := getFiles(root)
	assert.Len(t, files, 1)

	err = DeleteKindFromYaml(path, "test1")
	assert.NoError(t, err)

	err = DeleteKindFromYaml(path, "test2")
	assert.NoError(t, err)

	var restTest Resource
	err = YamlToStruct(path, &restTest)
	assert.NoError(t, err)

	assert.Empty(t, &restTest)

	err = os.RemoveAll(root)
	assert.NoError(t, err)
}

func TestHelper_DeleteFirstResourceFromYaml_first(t *testing.T) {
	root := "/tmp/nonexistent"
	err := os.MkdirAll(root, os.ModePerm)
	assert.NoError(t, err)

	path := "/tmp/nonexistent/test.yaml"
	err = AddStructToYaml(path, firstResource)
	assert.NoError(t, err)
	err = AddStructToYaml(path, secondResource)
	assert.NoError(t, err)

	files := getFiles(root)
	assert.Len(t, files, 1)

	err = DeleteFirstResourceFromYaml(path, "api/v1", "test1", "test1")
	assert.NoError(t, err)

	var restTest Resource
	err = YamlToStruct(path, &restTest)
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(secondResource, &restTest))

	err = os.RemoveAll(root)
	assert.NoError(t, err)
}

func TestHelper_DeleteFirstResourceFromYaml_second(t *testing.T) {
	root := "/tmp/nonexistent"
	err := os.MkdirAll(root, os.ModePerm)
	assert.NoError(t, err)

	path := "/tmp/nonexistent/test.yaml"
	err = AddStructToYaml(path, firstResource)
	assert.NoError(t, err)
	err = AddStructToYaml(path, secondResource)
	assert.NoError(t, err)

	files := getFiles(root)
	assert.Len(t, files, 1)

	err = DeleteFirstResourceFromYaml(path, "api/v1", "test2", "test2")
	assert.NoError(t, err)

	var restTest Resource
	err = YamlToStruct(path, &restTest)
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(firstResource, &restTest))

	err = os.RemoveAll(root)
	assert.NoError(t, err)
}

func TestHelper_DeleteFirstResourceFromYaml_both(t *testing.T) {
	root := "/tmp/nonexistent"
	err := os.MkdirAll(root, os.ModePerm)
	assert.NoError(t, err)

	path := "/tmp/nonexistent/test.yaml"
	err = AddStructToYaml(path, firstResource)
	assert.NoError(t, err)
	err = AddStructToYaml(path, secondResource)
	assert.NoError(t, err)

	files := getFiles(root)
	assert.Len(t, files, 1)

	err = DeleteFirstResourceFromYaml(path, "api/v1", "test1", "test1")
	assert.NoError(t, err)

	err = DeleteFirstResourceFromYaml(path, "api/v1", "test2", "test2")
	assert.NoError(t, err)

	var restTest Resource
	err = YamlToStruct(path, &restTest)
	assert.NoError(t, err)

	assert.Empty(t, &restTest)

	err = os.RemoveAll(root)
	assert.NoError(t, err)
}
