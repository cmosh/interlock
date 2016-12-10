package utils

import (
	"github.com/cmosh/interlock/ext"
	ctypes "github.com/docker/engine-api/types/container"
)

const (
	DefaultBalanceAlgorithm = "roundrobin"
)

func BalanceAlgorithm(config *ctypes.Config) string {
	algo := DefaultBalanceAlgorithm

	if v, ok := config.Labels[ext.InterlockBalanceAlgorithmLabel]; ok {
		algo = v
	}

	return algo
}
