package config

// ProvisionOptions contain various option to customize how provision works.
type ProvisionOptions struct {
	// Path to kube configuration if available
	KubeConfigPath string
}
