package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/labels"
	logcontext "github.com/caos/orbiter/logging/context"
	"github.com/caos/orbiter/logging/stdlib"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {

	var localMode bool

	flag.BoolVar(&localMode, "local-mode", false, "Disable the controller manager and only use the operator to handle gitcrds")
	flag.Parse()

	if localMode {
		clientgo.InConfig = false
	}

	logger := logcontext.Add(stdlib.New(os.Stdout))

	resources, err := clientgo.ListResources(logger, labels.GetGlobalLabels())
	if err != nil {
		panic(err)
	}
	for _, resource := range resources {
		fmt.Println(resource)
	}
}
