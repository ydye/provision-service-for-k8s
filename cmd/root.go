package cmd

import (
	"fmt"
	"github.com/ydye/provision-service-for-k8s/pkg/config"
	"github.com/ydye/provision-service-for-k8s/pkg/core"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile              string
	kubeConfigFile       string
	ansiblePlaybooksPath string
	ignoredLable         map[string]string
	period               time.Duration
	provisonNodeBulk     int

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
			interrupt := make(chan struct{})
			opts := core.ProvisionServiceOptions{
				ProvisionOptions: provisionOptions,
				KubeClient:       kubeClient,
			}
			provision := core.NewDefaultProvison(opts, interrupt)
			go wait.Until(func() {
				if err := provision.RunOnce(time.Now()); err != nil {
					klog.Fatalf("Error occurs when provisioning. %v", err)
					provision.ExitCleanUp()
				}
			}, period, interrupt)
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
	RootCmd.PersistentFlags().StringVar(&kubeConfigFile, "kube-config-file", "", "The path of kubeConfig file")
	RootCmd.PersistentFlags().StringToStringVarP(&ignoredLable, "ignored-label", "i", nil, "Ignore the node with the label and value")
	RootCmd.PersistentFlags().IntVar(&provisonNodeBulk, "provision-node-bulk", 1, "Max nodes numbers to provision in one round")
	RootCmd.PersistentFlags().StringVar(&ansiblePlaybooksPath, "ansible-playbooks-path", "~/provision-openpai/playbooks",
		"The Path to store ansible-playbooks which is used when provison a node.")
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
		kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			klog.Fatal("Failed to build kubeConfig: %v", err)
		}
		return kubeconfig
	} else {
		kubeconfig, err := rest.InClusterConfig()
		if err != nil {
			klog.Fatal("Failed to build kubeConfig: %v", err)
		}
		return kubeconfig
	}
}

func createKubeClient(kubeConfig *rest.Config) kubernetes.Interface {
	return kubernetes.NewForConfigOrDie(kubeConfig)
}

func createProvisionOptions() config.ProvisionOptions {
	return config.ProvisionOptions{
		KubeConfigPath:       kubeConfigFile,
		AnsiblePlaybooksPath: ansiblePlaybooksPath,
		Period:               period,
		IgnoredLable:         ignoredLable,
		ProvisionNodeBulk:    provisonNodeBulk,
	}
}
