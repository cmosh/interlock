package utils

import (
	"github.com/cmosh/interlock/ext"
	ctypes "github.com/docker/engine-api/types/container"
)

func Domain(config *ctypes.Config) string {
	domain := config.Domainname

	if v, ok := config.Labels[ext.InterlockDomainLabel]; ok {
		domain = v
	}

	return domain
}
