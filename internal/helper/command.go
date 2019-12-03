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

	var command string
	for _, arg := range cmd.Args {
		if strings.Contains(arg, " ") {
			command += " \\\"" + arg + "\\\""
			continue
		}
		command += " " + arg
	}
	command = command[1:]

	logger.WithFields(map[string]interface{}{
		"cmd": command,
	}).Debug("Executing")

	if logger.IsVerbose() {
		cmd.Stdout = os.Stdout
	}

	var buf bytes.Buffer
	cmd.Stderr = &buf
	return errors.Wrap(cmd.Run(), buf.String())
}
