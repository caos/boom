package clientgo

import (
	"path/filepath"
	"strings"

	"github.com/caos/orbiter/logging"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

type ResourceInfo struct {
	Group      string
	Version    string
	Resource   string
	Namespaced bool
}

type Resource struct {
	Group     string
	Version   string
	Resource  string
	Kind      string
	Name      string
	Namespace string
	Labels    map[string]string
}

func GetResource(group, version, resource, namespace, name string) (*Resource, error) {
	res := schema.GroupVersionResource{Group: group, Version: version, Resource: resource}
	conf, err := getClusterConfig()
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

func DeleteResource(resource *Resource) error {
	res := schema.GroupVersionResource{Group: resource.Group, Version: resource.Version, Resource: resource.Resource}
	conf, err := getClusterConfig()
	if err != nil {
		return err
	}

	client, err := dynamic.NewForConfig(conf)
	if err != nil {
		return err
	}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	clientRes := client.Resource(res)
	if resource.Namespace != "" {
		err = clientRes.Namespace(resource.Namespace).Delete(resource.Name, deleteOptions)
	} else {
		err = clientRes.Delete(resource.Name, deleteOptions)
	}

	return errors.Wrapf(err, "Error while deleting %s", resource.Name)
}

func getGroupVersionsResources() ([]*ResourceInfo, error) {
	conf, err := getClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := discovery.NewDiscoveryClientForConfig(conf)
	if err != nil {
		return nil, err
	}

	apiGroups, err := client.ServerGroups()
	if err != nil {
		return nil, err
	}

	resourceInfoList := make([]*ResourceInfo, 0)
	for _, apiGroup := range apiGroups.Groups {
		for _, version := range apiGroup.Versions {
			groupVersion := filepath.Join(apiGroup.Name, version.Version)
			apiResources, err := client.ServerResourcesForGroupVersion(groupVersion)
			if err != nil {
				return nil, err
			}

			for _, apiResource := range apiResources.APIResources {

				resourceInfo := &ResourceInfo{
					Group:      apiGroup.Name,
					Version:    version.Version,
					Resource:   apiResource.Name,
					Namespaced: apiResource.Namespaced,
				}
				parts := strings.Split(resourceInfo.Resource, "/")
				if len(parts) == 1 && resourceInfo.Resource != "componentstatuses" {
					resourceInfoList = append(resourceInfoList, resourceInfo)
				}
			}
		}
	}
	return resourceInfoList, nil
}

func ListResources(logger logging.Logger, labels map[string]string) ([]*Resource, error) {
	resourceInfoList, err := getGroupVersionsResources()
	if err != nil {
		return nil, err
	}

	conf, err := getClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := dynamic.NewForConfig(conf)
	if err != nil {
		return nil, err
	}

	labelSelector := ""
	for k, v := range labels {
		if labelSelector == "" {
			labelSelector = strings.Join([]string{k, v}, "=")
		} else {
			keyValue := strings.Join([]string{k, v}, "=")
			labelSelector = strings.Join([]string{labelSelector, keyValue}, ", ")
		}
	}

	resourceList := make([]*Resource, 0)
	for _, resourceInfo := range resourceInfoList {
		gvr := schema.GroupVersionResource{Group: resourceInfo.Group, Version: resourceInfo.Version, Resource: resourceInfo.Resource}
		list, err := client.Resource(gvr).List(metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			continue
		}
		for _, item := range list.Items {
			name, found, err := unstructured.NestedString(item.Object, "metadata", "name")
			if err != nil || !found {
				return nil, err
			}

			kind, _, err := unstructured.NestedString(item.Object, "kind")
			if err != nil {
				return nil, err
			}

			namespace, _, err := unstructured.NestedString(item.Object, "metadata", "namespace")
			if err != nil {
				return nil, err
			}

			labels, _, err := unstructured.NestedMap(item.Object, "metadata", "labels")
			if err != nil {
				return nil, err
			}

			labelStrs := make(map[string]string)
			for k, label := range labels {
				labelStrs[k] = label.(string)
			}

			resourceList = append(resourceList, &Resource{
				Group:     resourceInfo.Group,
				Version:   resourceInfo.Version,
				Resource:  resourceInfo.Resource,
				Name:      name,
				Kind:      kind,
				Namespace: namespace,
				Labels:    labelStrs,
			})
		}
	}

	return resourceList, nil
}
