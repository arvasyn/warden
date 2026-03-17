package photon

import (
	"slices"

	"github.com/arvasyn/warden/internal/pkg/sandbox"
	"github.com/rs/zerolog/log"
)

func (c *Config) AddCapabilities(app sandbox.Manifest) {
	var allowedCapabilities = []string{
		"CAP_NET_BIND_SERVICE",
		"CAP_SETUID",
		"CAP_SETGID",
		"CAP_CHOWN",
		"CAP_FOWNER",
	}

	for _, capability := range app.Sandbox.Capabilities.Add {
		if !slices.Contains(allowedCapabilities, capability) {
			log.Warn().Str("application", app.Application.Bundle).
				Str("capability", capability).
				Msg("Application tried adding unallowed capability")

			continue
		}

		log.Debug().Str("application", app.Application.Bundle).
			Str("capability", capability).
			Msg("Adding capability")

		c.Arguments = append(c.Arguments, "--cap-add", capability)
	}

	if app.Sandbox.Capabilities.DropAll {
		c.Arguments = append(c.Arguments, "--cap-drop", "ALL")
	}
}
