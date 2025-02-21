package utils

import (
	"strings"

	ctypes "github.com/docker/engine-api/types/container"
	"github.com/cmosh/interlock/ext"
)

func WebsocketEndpoints(config *ctypes.Config) []string {
	websocketEndpoints := []string{}

	for l, v := range config.Labels {
		// this is for labels like interlock.websocket_endpoint.1=foo
		if strings.Index(l, ext.InterlockWebsocketEndpointLabel) > -1 {
			websocketEndpoints = append(websocketEndpoints, v)
		}
	}

	return websocketEndpoints
}
