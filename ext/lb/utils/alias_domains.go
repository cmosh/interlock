package utils

import (
	"strings"

	"github.com/cmosh/interlock/ext"
	ctypes "github.com/docker/engine-api/types/container"
)

func AliasDomains(config *ctypes.Config) []string {
	aliasDomains := []string{}

	for l, v := range config.Labels {
		// this is for labels like interlock.alias_domain.1=foo.local
		if strings.Index(l, ext.InterlockAliasDomainLabel) > -1 {
			aliasDomains = append(aliasDomains, v)
		}
	}

	return aliasDomains
}
