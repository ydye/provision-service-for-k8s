package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	RootCmd = &cobra.Command{
		Use:   "provision",
		Short: "Provision necessary service on a new arrived kubernetes node.",
		Long: `Provision service has 2 stage. 

At stage 1, provision service will install service or software on the host 
through ansible-playbooks. At this stage, node will be labeled with 
provision=running.

At stage 2, provision service will add a necessary label to the new node so
that some agent services will be assigned to the new node through kubernetes.
At this stage, node will be labeled with provision=running. And the configured 
label will be added into this node.

At each stage, if the provision task failed, a taint will be added into the node
to prevent pod to be scheduled into the node. And the node will be labeled as 
provision=failed.

If all the provision tasks success, a label provision=successful will be added
into the node.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (defailt is ./provision.yaml)")

}

func initConfig() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigName("provision")
	if cfgFile != "" {
		fmt.Println(">>>> cfgFile: ", cfgFile)
		viper.SetConfigFile(cfgFile)
		configDir := path.Dir(cfgFile)
		if configDir != "." && configDir != dir {
			viper.AddConfigPath(configDir)
		}
	}

	viper.AddConfigPath(dir)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}
}
