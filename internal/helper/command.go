package helper

import (
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

	kubectlLogger := logger.WithFields(map[string]interface{}{
		"cmd": command,
	})

	kubectlLogger.Debug("Executing")

	out, err := cmd.CombinedOutput()
	kubectlLogger.Debug(string(out))

	return errors.Wrapf(err, "Error while executing command: Response: %s", string(out))
}
