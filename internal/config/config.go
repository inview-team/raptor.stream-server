package config

import "github.com/pion/webrtc/v3"

type Settings struct {
	RegistryAddress     string `env:"REGISTRY_ADDRESS"`
	StreamServerAddress string `env:"STREAMSERVER_ADDRESS"`
}

var PeerConnectionConfig = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	},
}
