package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	kube_flag "k8s.io/component-base/cli/flag"
)

// NewProvisionCommand creates a *cobra.Command object with default parameters
func NewProvisionCommand() *cobra.Command {
	cleanFlagSet := pflag.NewFlagSet("provision", pflag.ContinueOnError)
	cleanFlagSet.SetNormalizeFunc(kube_flag.WordSepNormalizeFunc)

}
