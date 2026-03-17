package photon

import (
	"slices"

	"github.com/arvasyn/warden/internal/pkg/sandbox"
	"github.com/rs/zerolog/log"
)

func (c *Config) AddNamespaces(manifest sandbox.Manifest) {
	var allowedNamespaces = []string{"mnt", "user", "pid", "net", "ipc", "uts"}

	for namespace, unshare := range manifest.Sandbox.Namespaces {
		if !slices.Contains(allowedNamespaces, namespace) {
			log.Warn().Str("application", manifest.Application.Bundle).
				Str("namespace", namespace).
				Msg("Attempted argument poisoning?")

			continue
		}

		if unshare {
			c.Arguments = append(c.Arguments, "--unshare-"+namespace)
		} else {
			if namespace == "pid" && c.AllowSharedPID == false {
				c.Arguments = append(c.Arguments, "--share-"+namespace)
			}
		}
	}
}
