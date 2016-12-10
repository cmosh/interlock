package utils

import (
	ctypes "github.com/docker/engine-api/types/container"
	"github.com/cmosh/interlock/ext"
)

func IPHash(config *ctypes.Config) bool {
	ipHash := false

	if _, ok := config.Labels[ext.InterlockIPHashLabel]; ok {
		ipHash = true
	}

	return ipHash
}
