package config

type Settings struct {
	RegistryAddress     string `env:"REGISTRY_ADDRESS"`
	StreamServerAddress string `env:"STREAMSERVER_ADDRESS"`
}
