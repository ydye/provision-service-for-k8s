package main

import (
	"github.com/ydye/provision-service-for-k8s/cmd"
	"github.com/ydye/provision-service-for-k8s/pkg/version"
	"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)
	klog.V(1).Infof("OpenPAI provision %s", version.OpenPAIProvisionVersion)
	cmd.Execute()
}
