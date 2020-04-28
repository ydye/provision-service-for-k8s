package config

import "time"

// ProvisionOptions contain various option to customize how provision works.
type ProvisionOptions struct {
	// Path to kube configuration if available
	KubeConfigPath       string
	AnsiblePlaybooksPath string
	Period               time.Duration
	IgnoredLable         map[string]string
	ProvisionNodeBulk    int
}
