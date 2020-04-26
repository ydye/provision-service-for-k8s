package core

import (
	kube_util "github.com/ydye/provision-service-for-k8s/pkg/utils/kubernetes"
	"github.com/ydye/provision-service-for-k8s/pkg/config"
	"github.com/ydye/provision-service-for-k8s/pkg/utils/errors"
	"k8s.io/client-go/kubernetes"
	"time"
)

// Provision
type ProvisionServiceOptions struct {
	config.ProvisionOptions
	KubeClient   kubernetes.Interface
	ListRegistry kube_util.ListerRegistry
}

type Provision interface {
	start() error

	RunOnce(currentTime time.Time) errors.ProvisionError

	ExitCleanUp()
}

func NewProvision(opts ProvisionServiceOptions) (Provision, errors.ProvisionError) {
	err := initializeDefaultOptions(&opts)
	if err != nil {
		return nil, errors.ToProvisionError(errors.InternalError, err)
	}
	return
}

func initializeDefaultOptions(opts *ProvisionServiceOptions) error {
	if opts.ListRegistry == nil {
		listerRegistryStopChannel := make(chan struct{})
		opts.ListRegistry = kube_util.NewListerRegistryWithDefaultListers(opts.KubeClient, listerRegistryStopChannel)
	}
	return nil
}



