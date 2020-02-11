package kubectl

import (
	"strings"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Version struct {
	BuildDate    string `yaml:"buildDate"`
	Compiler     string `yaml:"compiler"`
	GitCommit    string `yaml:"gitCommit"`
	GitTreeState string `yaml:"gitTreeState"`
	GitVersion   string `yaml:"gitVersion"`
	GoVersion    string `yaml:"goVersion"`
	Major        string `yaml:"major"`
	Minor        string `yaml:"minor"`
	Platform     string `yaml:"platform"`
}

type Versions struct {
	ClientVersion *Version `yaml:"clientVersion"`
	ServerVersion *Version `yaml:"serverVersion"`
}

type KubectlVersion struct {
	kubectl *Kubectl
}

func NewVersion() *KubectlVersion {
	return &KubectlVersion{kubectl: New("version").AddParameter("-o", "yaml")}
}

func (k *KubectlVersion) GetKubeVersion(logger logging.Logger) (string, error) {
	cmd := k.kubectl.Build()

	kubectlLogger := logger.WithFields(map[string]interface{}{
		"cmd":   cmd,
		"logId": "CMD-sN18gqW3pTG8rUR",
	})
	kubectlLogger.Debug("Executing")

	out, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.Wrap(err, "Error while executing command")
		kubectlLogger.Error(err)
		return "", err
	}

	versions := &Versions{}
	err = yaml.Unmarshal(out, versions)
	if err != nil {
		err = errors.Wrap(err, "Error while unmarshaling output")
		kubectlLogger.Error(err)
		return "", err
	}

	parts := strings.Split(versions.ServerVersion.GitVersion, "-")
	return parts[0], nil
}
