package gallium

import (
	"slices"

	"github.com/rs/zerolog/log"
)

func (c *Config) AddNamespaces(manifest Manifest) {
	var allowedNamespaces = []string{"mnt", "pid", "user", "net", "ipc", "uts"}

	for namespace, unshare := range manifest.Sandbox.Namespaces {
		if !slices.Contains(allowedNamespaces, namespace) {
			log.Warn().
				Str("application", manifest.Application.Bundle).
				Str("namespace", namespace).
				Msg("Attempted argument poisoning?")

			continue
		}

		if unshare {
			c.Arguments = append(c.Arguments, "--unshare-"+namespace)
		} else {
			if namespace == "pid" && c.AllowSharedPID == true {
				c.Arguments = append(c.Arguments, "--share-"+namespace)
			}
		}
	}
}
