package kubectl

import (
	"strings"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
)

type KubectlLabel struct {
	kubectl *Kubectl
}

func NewLabel(resultFilePath string) *KubectlLabel {
	return &KubectlLabel{kubectl: New("label").AddFlag("--overwrite").AddParameter("-f", resultFilePath)}
}
func (k *KubectlLabel) AddParameter(key, value string) *KubectlLabel {
	k.kubectl.AddParameter(key, value)
	return k
}

func (k *KubectlLabel) Apply(logger logging.Logger, labels map[string]string) error {
	for key, value := range labels {
		k.kubectl.AddFlag(strings.Join([]string{key, value}, "="))
	}

	cmd := k.kubectl.Build()

	kubectlLogger := logger.WithFields(map[string]interface{}{
		"cmd": cmd,
	})
	kubectlLogger.Debug("Executing")

	out, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.Wrap(err, "Error while executing command")
		kubectlLogger.Debug(string(out))
		kubectlLogger.Error(err)
	}

	return err
}
