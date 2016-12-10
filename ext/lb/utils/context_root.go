package utils

import (
	"github.com/cmosh/interlock/ext"
	ctypes "github.com/docker/engine-api/types/container"
)

func ContextRoot(config *ctypes.Config) string {
	if v, ok := config.Labels[ext.InterlockContextRootLabel]; ok {
		return v
	}

	return ""
}

func ContextRootRewrite(config *ctypes.Config) bool {
	if _, ok := config.Labels[ext.InterlockContextRootRewriteLabel]; ok {
		return true
	}

	return false
}
