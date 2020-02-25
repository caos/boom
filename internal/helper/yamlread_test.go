package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type teststruct struct {
	Test string `yaml:"test"`
}

func TestHelper_YamlToString(t *testing.T) {
	var str string
	str, err := YamlToString("testfiles/struct.yaml")
	assert.NoError(t, err)

	assert.Equal(t, "test: test", str)
}

func TestHelper_YamlToString_nonexistent(t *testing.T) {
	var str string
	str, err := YamlToString("testfiles/nonexistent.yaml")
	assert.Error(t, err)
	assert.Empty(t, str)
}

func TestHelper_YamlToStruct(t *testing.T) {
	var teststruct teststruct
	err := YamlToStruct("testfiles/struct.yaml", &teststruct)
	assert.NoError(t, err)

	assert.NotNil(t, teststruct)
	assert.Equal(t, "test", teststruct.Test)
}

func TestHelper_YamlToStruct_nonexistent(t *testing.T) {
	var teststruct teststruct
	err := YamlToStruct("testfiles/nonexistent.yaml", &teststruct)
	assert.Error(t, err)
}
