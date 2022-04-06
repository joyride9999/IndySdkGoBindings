package indySDK

import "github.com/joyride9999/IndySdkGoBindings/mod"

// SetRuntimeConfig set libindy runtime configuration
func SetRuntimeConfig(config mod.Config) error {
	channel := mod.SetRuntimeConfig(config)
	result := <-channel
	return result.Error
}
