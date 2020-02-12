package clientgo

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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

func GetSecret(name, namespace string) (*v1.Secret, error) {
	conf, err := GetClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := GetClientSet(conf)
	if err != nil {
		return nil, err
	}

	secret, err := clientset.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, errors.New("Secret not found")
	}

	return secret, nil
}

type Resource struct {
	Kind      string
	Name      string
	Namespace string
	Labels    map[string]string
}

func GetResource(group, version, resource, namespace, name string) (*Resource, error) {
	res := schema.GroupVersionResource{Group: group, Version: version, Resource: resource}

	conf, err := GetClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := dynamic.NewForConfig(conf)
	if err != nil {
		return nil, err
	}

	result, err := clientset.Resource(res).Namespace(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &Resource{
		Kind:      result.GetKind(),
		Name:      result.GetName(),
		Namespace: result.GetNamespace(),
		Labels:    result.GetLabels(),
	}, nil
}
