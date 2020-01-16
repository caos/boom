package application

import (
	"testing"

	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/boom/internal/name"
	"github.com/golang/mock/gomock"
)

const (
	Name      name.Application = "mockapplication"
	Namespace                  = "mocknamespace"
)

type TestApplication struct {
	app *MockApplication
}

func NewTestApplication(t *testing.T) *TestApplication {

	mockApp := NewMockApplication(gomock.NewController(t))
	mockApp.EXPECT().GetName().AnyTimes().DoAndReturn(Name)
	mockApp.EXPECT().GetNamespace().AnyTimes().DoAndReturn(Namespace)
	app := &TestApplication{
		app: mockApp,
	}

	return app
}

func (a *TestApplication) Application() *MockApplication {
	return a.app
}

func (a *TestApplication) SetDeploy(spec *v1beta1.ToolsetSpec, b bool) *TestApplication {
	a.app.EXPECT().Deploy(spec).AnyTimes().Return(b)
	return a
}

func (a *TestApplication) SetInitial(b bool) *TestApplication {
	a.app.EXPECT().Initial().AnyTimes().Return(b)
	return a
}

func (a *TestApplication) SetChanged(spec *v1beta1.ToolsetSpec, b bool) *TestApplication {
	a.app.EXPECT().Changed(spec).AnyTimes().Return(b)
	return a
}

func (a *TestApplication) SetAppliedSpec(spec *v1beta1.ToolsetSpec) *TestApplication {
	a.app.EXPECT().SetAppliedSpec(spec).AnyTimes()
	return a
}

type TestYAMLApplication struct {
	app *MockYAMLApplication
}

func NewTestYAMLApplication(t *testing.T) *TestYAMLApplication {

	mockApp := NewMockYAMLApplication(gomock.NewController(t))
	mockApp.EXPECT().GetName().AnyTimes().Return(Name)
	mockApp.EXPECT().GetNamespace().AnyTimes().Return(Namespace)
	app := &TestYAMLApplication{
		app: mockApp,
	}

	return app
}

func (a *TestYAMLApplication) Application() *MockYAMLApplication {
	return a.app
}

func (a *TestYAMLApplication) SetDeploy(spec *v1beta1.ToolsetSpec, b bool) *TestYAMLApplication {
	a.app.EXPECT().Deploy(spec).AnyTimes().DoAndReturn(b)
	return a
}

func (a *TestYAMLApplication) SetInitial(b bool) *TestYAMLApplication {
	a.app.EXPECT().Initial().AnyTimes().DoAndReturn(b)
	return a
}

func (a *TestYAMLApplication) SetChanged(spec *v1beta1.ToolsetSpec, b bool) *TestYAMLApplication {
	a.app.EXPECT().Changed(spec).AnyTimes().DoAndReturn(b)
	return a
}

func (a *TestYAMLApplication) AllowSetAppliedSpec(spec *v1beta1.ToolsetSpec) *TestYAMLApplication {
	a.app.EXPECT().SetAppliedSpec(spec).AnyTimes()
	return a
}

func (a *TestYAMLApplication) SetGetYaml(yaml string) *TestYAMLApplication {
	a.app.EXPECT().GetYaml().AnyTimes().Return(yaml)
	return a
}
