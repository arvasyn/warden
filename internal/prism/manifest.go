package prism

type Configuration struct {
	// Current Warden configuration
	Version    int         `yaml:"version"`
	Enrollment *Enrollment `yaml:"enrollment"`
	Modules    []Module    `yaml:"modules"`
}

type Enrollment struct {
	// This should be the organization's slug (example.com)
	Organization string `yaml:"organization"`
	Token        string `yaml:"token"`
}

type Module struct {
	Name    string            `yaml:"name"`
	Enabled bool              `yaml:"enabled"`
	Data    map[string]string `yaml:"data"`
}
