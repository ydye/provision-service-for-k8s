package main

import (
	"github.com/ydye/provision-service-for-k8s/pkg/version"
	kube_flag "k8s.io/component-base/cli/flag"
	"k8s.io/klog"
)

func run() {

}

func main() {
	klog.InitFlags(nil)
	kube_flag.InitFlags()
	klog.V(1).Infof("OpenPAI provision %s", version.OpenPAIProvisionVersion)

}
