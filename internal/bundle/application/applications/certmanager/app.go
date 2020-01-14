package certmanager

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/caos/orbiter/logging"

	"github.com/caos/boom/api/v1beta1"
	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/bundle/application/applications/certmanager/helm"
	"github.com/caos/boom/internal/bundle/application/resources/service"
	"github.com/caos/boom/internal/name"
	"github.com/caos/boom/internal/templator/helm/chart"
	"github.com/caos/boom/internal/helper"
)

const (
	applicationName name.Application = "cert-manager"
)

func GetName() name.Application {
	return applicationName
}

type CertManager struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.CertManager
}

func New(logger logging.Logger) *CertManager {
	c := &CertManager{
		logger: logger,
	}

	return c
}
func (c *CertManager) GetName() name.Application {
	return applicationName
}

func Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.CertManager.Deploy
}

func (c *CertManager) Initial() bool {
	return c.spec == nil
}

func (c *CertManager) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.CertManager, c.spec)
}

func (c *CertManager) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	c.spec = toolsetCRDSpec.CertManager
}

func (c *CertManager) GetNamespace() string {
	return "caos-system"
}

func (c *CertManager) HelmPreApplySteps(spec *v1beta1.ToolsetSpec) ([]interface{}, error) {

	crdDirectoryPath := filepath.Join("..", "..", "tools")

	crds, err := getCrds(crdDirectoryPath)
	if err != nil {
		return nil, err
	}
	ret := make([]interface{}, len(crds))
	for n, crd := range crds {
		ret[n] = crd
	}

	svc := getService(c.GetNamespace())
	ret = append(ret, svc)

	return ret, nil
}

func getService(namespace string) *service.Service {
	name := "cert-manager"
	svc := service.New(&service.Config{
		Name:       name,
		Namespace:  namespace,
		Labels:     map[string]string{"app": "cert-manager"},
		Protocol:   "TCP",
		Port:       9402,
		TargetPort: 9402,
		Selector:   map[string]string{"app": "cert-manager"},
	})

	return svc
}

func getCrds(toolsDirectoryPath string) ([]string, error) {
	chartsPath := filepath.Join(toolsDirectoryPath, "charts", applicationName.String(), "crds")

	ret := make([]string, 0)
	var files []string
	err := filepath.Walk(chartsPath, func(path string, info os.FileInfo, err error) error {
		if path != chartsPath {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileStr, err := helper.YamlToString(file)
		if err != nil {
			return nil, err
		}
		ret = append(ret, fileStr)
	}

	return ret, nil
}

func (c *CertManager) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	spec := toolset.CertManager
	values := helm.DefaultValues(c.GetImageTags())

	if spec.ReplicaCount != 0 {
		values.ReplicaCount = spec.ReplicaCount
	}

	values.Global.LeaderElection = &helm.LeaderElection{
		Namespace: c.GetNamespace(),
	}

	return values
}

func (c *CertManager) GetChartInfo() *chart.Chart {
	return helm.GetChartInfo()
}

func (c *CertManager) GetImageTags() map[string]string {
	return helm.GetImageTags()
}
