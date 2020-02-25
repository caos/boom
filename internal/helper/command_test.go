package helper

import (
	"os"
	"os/exec"
	"testing"

	"github.com/caos/orbiter/logging"
	logcontext "github.com/caos/orbiter/logging/context"
	"github.com/caos/orbiter/logging/kubebuilder"
	"github.com/caos/orbiter/logging/stdlib"
	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"
)

func newLogger() logging.Logger {
	logger := logcontext.Add(stdlib.New(os.Stdout))
	ctrl.SetLogger(kubebuilder.New(logger))

	return logger
}

func TestHelper_Run(t *testing.T) {
	logger := newLogger()

	cmd := exec.Command("echo", "first")
	err := Run(logger, *cmd)
	assert.NoError(t, err)
}

func TestHelper_Run_MoreArgs(t *testing.T) {
	logger := newLogger()

	cmd := exec.Command("echo", "first", "second")
	err := Run(logger, *cmd)
	assert.NoError(t, err)
}

func TestHelper_Run_UnknowCommand(t *testing.T) {
	logger := newLogger()

	cmd := exec.Command("unknowncommand", "first")
	err := Run(logger, *cmd)
	assert.Error(t, err)
}

func TestHelper_Run_ErrorCommand(t *testing.T) {
	logger := newLogger()

	cmd := exec.Command("ls", "/unknownfolder/unknownsubfolder")
	err := Run(logger, *cmd)
	assert.Error(t, err)
}
