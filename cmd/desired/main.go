package main

import (
	"flag"
	"fmt"

	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/desired"
	"github.com/caos/boom/internal/name"

	"github.com/caos/orbiter/mntr"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var app name.Application = "loki"

func main() {

	var localMode bool

	verbose := flag.Bool("verbose", false, "Print logs for debugging")
	flag.BoolVar(&localMode, "local-mode", false, "Disable the controller manager and only use the operator to handle gitcrds")
	flag.Parse()

	if localMode {
		clientgo.InConfig = false
	}

	monitor := mntr.Monitor{
		OnInfo:   mntr.LogMessage,
		OnChange: mntr.LogMessage,
		OnError:  mntr.LogError,
	}
	if *verbose {
		monitor = monitor.Verbose()
	}
	resultFilePath := "../../local/tools/loki/caos/results/results.yaml"
	namespace := "caos-system"

	desiredResources, err := desired.Get(monitor, resultFilePath, namespace, app)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, resource := range desiredResources {
		fmt.Println(resource)
	}
}
