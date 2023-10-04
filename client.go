// Package rgclient implements client-go style API client to access
// RouteGroups and Kubernetes resources.
package rgclient

import (
	"errors"
	"fmt"

	zclient "github.com/szuecs/routegroup-client/client/clientset/versioned"
	zalandov1 "github.com/szuecs/routegroup-client/client/clientset/versioned/typed/zalando.org/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	LocalAPIServer = "http://127.0.0.1:8001"
)

// Interface is the unified interface for Kubernetes and Zalando
// RouteGroup access
type Interface interface {
	kubernetes.Interface
	ZalandoInterface
}

// ZalandoInterface to access RouteGroups
type ZalandoInterface interface {
	ZalandoV1() zalandov1.ZalandoV1Interface
}

// Clientset is the unified client implementation
type Clientset struct {
	*kubernetes.Clientset
	zclient *zclient.Clientset
}

type Options struct {
	TokenFile string
}

// ZalandoV1 implements ZalandoInterface
func (cs *Clientset) ZalandoV1() zalandov1.ZalandoV1Interface {
	return cs.zclient.ZalandoV1()
}

// NewClientset returns the unified client, but users should in
// general prefer rgclient.CreateUnified().
func NewClientset(config *rest.Config) (*Clientset, error) {
	zcli, err := zclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	k8scli, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Clientset{
		Clientset: k8scli,
		zclient:   zcli,
	}, err
}

// Create returns the Zalandointerface client
func Create() (ZalandoInterface, error) {
	config, err := getRestConfigWithOptions(nil)
	if err != nil {
		return nil, err
	}

	return zclient.NewForConfig(config)
}

// CreateUnified returns the unified client
func CreateUnified() (Interface, error) {
	config, err := getRestConfigWithOptions(nil)
	if err != nil {
		return nil, err
	}

	client, err := NewClientset(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create a unified client: %v", err)
	}

	return client, nil
}

// CreateUnifiedWithOptions returns the unified client that
func CreateUnifiedWithOptions(opts *Options) (Interface, error) {
	config, err := getRestConfigWithOptions(opts)
	if err != nil {
		return nil, err
	}

	client, err := NewClientset(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create a unified client: %v", err)
	}

	return client, nil
}

func getRestConfigWithOptions(opts *Options) (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		if errors.Is(err, rest.ErrNotInCluster) {
			if opts == nil {
				config = &rest.Config{
					Host: LocalAPIServer,
				}
			} else {
				config = &rest.Config{
					Host:            LocalAPIServer,
					BearerTokenFile: opts.TokenFile,
				}
			}
			err = nil
		} else {
			return nil, err
		}
	}
	return config, err
}
