package kubernetes

import (
	"fmt"
	"time"

	"github.com/ydye/provision-service-for-k8s/config"
	apiv1 "k8s.io/api/core/v1"
)

// IsNodeReadyAndSchedulable returns true if the node is ready and schedulable.
func IsNodeReadyAndSchedulable(node *apiv1.Node) bool {
	ready, _, _ := GetReadinessState(node)
	if !ready {
		return false
	}
	// Ignore nodes that are marked unschedulable
	if node.Spec.Unschedulable {
		return false
	}
	return true
}

// IsNodeProvisioned return true if the node is provisioned successfully.
func IsNodeProvisionedAndSuccessed(node *apiv1.Node) bool {
	provisioned, _ := GetProvisionState(node)
	if !provisioned {
		return false
	}
	return true
}

func IsNodeNeededToProvision

func GetReadinessState(node *apiv1.Node) (isNodeReady bool, lastTransitionTime time.Time, err error) {
	canNodeBeReady, readyFound := true, false
	lastTransitionTime = time.Time{}

	for _, cond := range node.Status.Conditions {
		switch cond.Type {
		case apiv1.NodeReady:
			readyFound = true
			if cond.Status == apiv1.ConditionFalse || cond.Status == apiv1.ConditionUnknown {
				canNodeBeReady = false
			}
			if lastTransitionTime.Before(cond.LastTransitionTime.Time) {
				lastTransitionTime = cond.LastTransitionTime.Time
			}
		case apiv1.NodeDiskPressure:
			if cond.Status == apiv1.ConditionTrue {
				canNodeBeReady = false
			}
			if lastTransitionTime.Before(cond.LastTransitionTime.Time) {
				lastTransitionTime = cond.LastTransitionTime.Time
			}
		case apiv1.NodeNetworkUnavailable:
			if cond.Status == apiv1.ConditionTrue {
				canNodeBeReady = false
			}
			if lastTransitionTime.Before(cond.LastTransitionTime.Time) {
				lastTransitionTime = cond.LastTransitionTime.Time
			}
		}
	}
	if !readyFound {
		return false, time.Time{}, fmt.Errorf("readiness information not found")
	}
	return canNodeBeReady, lastTransitionTime, nil
}

func GetProvisionState(node *apiv1.Node) (isNodeReady bool, err error) {
	provisioned := false

	if val, ok := node.Labels["provision"]; ok {
		if val == config.SuccessfulProvision {
			provisioned = true
		}
	}

	return provisioned, nil
}