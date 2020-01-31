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

	logFields := map[string]interface{}{
		"cmd":   command,
		"logId": "CMD-sN18gqW3pTG8rUR",
	}

	logger.WithFields(logFields).Debug("Executing")

	out, err := cmd.Output()
	logger.WithFields(logFields).Debug(string(out))

	return errors.Wrap(err, "Error while executing command")
}
