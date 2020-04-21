package core

import (
	"github.com/ydye/provision-service-for-k8s/config"
	"k8s.io/client-go/kubernetes"
)

// Provision
type ProvisionServiceOptions struct {
	config.ProvisionOptions
	KubeClient kubernetes.Interface
}

type Provision interface {
	start() error
}

