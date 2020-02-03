package helper

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var InConfig = true

func GetClusterConfig() (*rest.Config, error) {
	if InConfig {
		return getInClusterConfig()
	} else {
		return getOutClusterConfig()
	}
}

func getOutClusterConfig() (*rest.Config, error) {
	kubeconfig := filepath.Join(homeDir(), ".kube", "config")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	return config, errors.Wrap(err, "Error while creating out-cluster config")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getInClusterConfig() (*rest.Config, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	return config, errors.Wrap(err, "Error while creating in-cluster config")
}

func GetClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, errors.Wrap(err, "Error while creating clientset")
}
