package loggingoperator

import (
	"reflect"

	"github.com/caos/orbiter/logging"

	toolsetsv1beta1 "github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/app/name"
)

const (
	applicationName name.Application = "logging-operator"
)

func GetName() name.Application {
	return applicationName
}

type LoggingOperator struct {
	logger logging.Logger
	spec   *toolsetsv1beta1.LoggingOperator
}

func New(logger logging.Logger) *LoggingOperator {
	lo := &LoggingOperator{
		logger: logger,
	}

	return lo
}
func (l *LoggingOperator) GetName() name.Application {
	return applicationName
}

func Deploy(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return toolsetCRDSpec.LoggingOperator.Deploy
}

func (l *LoggingOperator) Initial() bool {
	return l.spec == nil
}

func (l *LoggingOperator) Changed(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) bool {
	return !reflect.DeepEqual(toolsetCRDSpec.LoggingOperator, l.spec)
}

func (l *LoggingOperator) SetAppliedSpec(toolsetCRDSpec *toolsetsv1beta1.ToolsetSpec) {
	l.spec = toolsetCRDSpec.LoggingOperator
}

func (l *LoggingOperator) GetNamespace() string {
	return "caos-system"
}

func (l *LoggingOperator) SpecToHelmValues(toolset *toolsetsv1beta1.ToolsetSpec) interface{} {
	// spec := toolset.LoggingOperator
	values := defaultValues(l.GetImageTags())

	// if spec.ReplicaCount != 0 {
	// 	values.ReplicaCount = spec.ReplicaCount
	// }

	return values
}

func defaultValues(imageTags map[string]string) *Values {
	return &Values{
		FullnameOverride: "logging-operator",
		ReplicaCount:     1,
		Image: Image{
			Repository: "banzaicloud/logging-operator",
			Tag:        imageTags["banzaicloud/logging-operator"],
			PullPolicy: "IfNotPresent",
		},
		HTTP: HTTP{
			Port: 8080,
			Service: Service{
				Type: "ClusterIP",
			},
		},
		RBAC: RBAC{
			Enabled: true,
			PSP: PSP{
				Enabled: true,
			},
		},
	}
}
