package preapply

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/caos/boom/api/v1beta1"
	"github.com/caos/orbiter/logging"
)

func (p *PreApply) GetYaml(logger logging.Logger, spec *v1beta1.ToolsetSpec) interface{} {
	files, err := ioutil.ReadDir(spec.PreApply.Folder)
	if err != nil {
		return nil
	}

	filesContent := ""
	for _, f := range files {
		filePath := filepath.Join(spec.PreApply.Folder, f.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil
		}

		if filesContent == "" {
			filesContent = string(data)
		} else {
			filesContent = strings.Join([]string{filesContent, "---", string(data)}, "\n")
		}
	}

	return filesContent
}
