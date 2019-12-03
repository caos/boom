package helper

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
)

func Run(logger logging.Logger, cmd exec.Cmd) error {

	logger.WithFields(map[string]interface{}{
		"cmd": strings.Join(cmd.Args, " "),
	}).Debug("Executing")

	if logger.IsVerbose() {
		cmd.Stdout = os.Stdout
	}

	var buf bytes.Buffer
	cmd.Stderr = &buf
	return errors.Wrap(cmd.Run(), buf.String())
}
