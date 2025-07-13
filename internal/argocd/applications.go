package argocd

import (
	"context"
	"fmt"

	argov1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	argoclientset "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func GetApplicationsByName(name string, namespace string) (*argov1.Application, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	clientset, err := argoclientset.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("clientset error: %w", err)
	}

	app, err := clientset.ArgoprojV1alpha1().Applications(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get application: %w", err)
	}

	return app, nil
}

func GetAllApplications(namespace string) ([]argov1.Application, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	clientset, err := argoclientset.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("clientset error: %w", err)
	}

	appList, err := clientset.ArgoprojV1alpha1().Applications(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list applications: %w", err)
	}

	return appList.Items, nil
}
