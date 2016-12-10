package utils

import (
	"strings"

	"github.com/cmosh/interlock/ext"
	ctypes "github.com/docker/engine-api/types/container"
)

func BackendOptions(config *ctypes.Config) []string {
	options := []string{}

	for l, v := range config.Labels {
		// this is for labels like interlock.backend_option.1=foo
		if strings.Index(l, ext.InterlockBackendOptionLabel) > -1 {
			options = append(options, v)
		}
	}

	return options
}
