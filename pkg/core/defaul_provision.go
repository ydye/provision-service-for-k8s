package core

import "time"

type defaultProvision struct {
	opts              ProvisionServiceOptions
	startTime         time.Time
	lastProvisionTime time.Time
}

func NewDefaultProvison(opts ProvisionServiceOptions) *defaultProvision {
	return &defaultProvision{
		opts:      opts,
		startTime: time.Now(),
	}
}






