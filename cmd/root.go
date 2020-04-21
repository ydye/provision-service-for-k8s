package cmd

import (
	"fmt"
	"github.com/ydye/provision-service-for-k8s/pkg/core"
	"log"
	"os"
	"time"
	"path"
	"path/filepath"

	"github.com/ydye/provision-service-for-k8s/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	cfgFile        string
	kubeConfigFile string
	period         time.Duration

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
			provisionOptions := createProvisionOptions()
			kubeClient := createKubeClient(getKubeConfig())

			opts := core.ProvisionServiceOptions{
				ProvisionOptions: provisionOptions,
				KubeClient: kubeClient,
			}
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
	RootCmd.PersistentFlags().DurationVar(&period, "period", 60*time.Second, "How often to find the new joined node")
	RootCmd.PersistentFlags().StringVar(&kubeConfigFile, "kubeConfigFile", "", "The path of kubeConfig file")
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

func getKubeConfig() *rest.Config {
	if kubeConfigFile != "" {
		klog.V(1).Infof("Using kubeconfig file: %s", kubeConfigFile)
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			klog.Fatal("Failed to build kubeConfig: %v", err)
		}
		return config
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			klog.Fatal("Failed to build kubeConfig: %v", err)
		}
		return config
	}
}

func createKubeClient(kubeConfig *rest.config) kubernetes.Interface {
	return kubernetes.NewForConfigOrDie(kubeConfig)
}

func createProvisionOptions() config.ProvisionOptions {
	return config.ProvisionOptions{
		KubeConfigPath: kubeConfigFile,
	}
}