package core

import (
	"github.com/ydye/provision-service-for-k8s/pkg/utils/errors"
	kube_util "github.com/ydye/provision-service-for-k8s/pkg/utils/kubernetes"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	"time"
)

type defaultProvision struct {
	ProvisionServiceOptions
	startTime         time.Time
	lastProvisionTime time.Time
	interrupt         chan struct{}
}

func NewDefaultProvison(opts ProvisionServiceOptions, interrupt chan struct{}) *defaultProvision {
	return &defaultProvision{
		ProvisionServiceOptions: opts,
		startTime:               time.Now(),
		interrupt:               interrupt,
	}
}

func (a *defaultProvision) ExitCleanUp() {
	close(a.interrupt)
}

func (a *defaultProvision) RunOnce(currentTime time.Time) errors.ProvisionError {
	targetNodeList, fdNodeErr := a.findNodeToProvision()
	if fdNodeErr != nil {
		klog.Errof("Failed to get node to provisions: %v", fdNodeErr)
		return fdNodeErr
	}
	if len(targetNodeList) == 0 {
		klog.Info("There is no node needed to be provision")
		return nil
	}

	for _, v := range targetNodeList {
		klog.Infof("Node: ", v)
	}
	return nil
}

func (a *defaultProvision) findNodeToProvision() ([]*apiv1.Node, errors.ProvisionError) {
	targetNodeToProvision, err := a.ListRegistry.UnprovisionedNodeLister().List()
	if err != nil {
		klog.Errorf("Failed to list all nodes: %v", err)
		return nil, errors.ToProvisionError(errors.ApiCallError, err)
	}
	filteredTargetNodeToProvision := kube_util.FilterOutNodesWithIgnoredLabel(a.IgnoredLable, targetNodeToProvision)
	return filteredTargetNodeToProvision, nil
}