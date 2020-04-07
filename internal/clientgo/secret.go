package clientgo

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ErrNotFound struct{}

func (e ErrNotFound) Error() string {
	return "Not found"
}

func GetSecret(name, namespace string) (s *v1.Secret, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error getting secret %s in namespace %s: %w", name, namespace, err)
		}
	}()

	conf, err := getClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := getClientSet(conf)
	if err != nil {
		return nil, err
	}

	secret, err := clientset.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, ErrNotFound{}
	}

	return secret, nil
}
