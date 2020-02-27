package main

import (
	"flag"
	"fmt"

	"github.com/caos/boom/internal/clientgo"
	"github.com/caos/boom/internal/labels"
	"github.com/caos/orbiter/mntr"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

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

	resourceInfoList, err := clientgo.GetGroupVersionsResources([]string{})

	resources, err := clientgo.ListResources(monitor, resourceInfoList, labels.GetGlobalLabels())
	if err != nil {
		panic(err)
	}
	for _, resource := range resources {
		fmt.Println(resource)
	}
}
