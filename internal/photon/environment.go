package photon

import (
	"os"
	"regexp"
	"slices"

	"github.com/arvasyn/warden/internal/pkg/sandbox"
	"github.com/rs/zerolog/log"
)

func (c *Config) AddEnvironment(app sandbox.Manifest) {
	for _, name := range app.Sandbox.Env.Passthrough {
		if !IsEnvAllowed(name) {
			log.Warn().Str("application", app.Application.Bundle).
				Str("name", name).
				Msg("Environment variable not allowed")

			continue
		}

		value, ok := os.LookupEnv(name)

		if !ok {
			log.Warn().Str("application", app.Application.Bundle).
				Msgf("Unable to get variable %s. Environment variable not set.", name)

			continue
		}

		log.Debug().Str("application", app.Application.Bundle).
			Msgf("Passing environment variable %s", name)

		c.Arguments = append(c.Arguments, "--setenv", name, value)
	}

	for name, value := range app.Sandbox.Env.Set {
		if !IsEnvAllowed(name) {
			log.Warn().Str("application", app.Application.Bundle).
				Str("name", name).
				Msg("Environment variable not allowed")

			continue
		}

		log.Debug().Str("application", app.Application.Bundle).
			Msgf("Setting environment variable %s to %s", name, value)

		c.Arguments = append(c.Arguments, "--setenv", name, value)
	}

	for _, name := range app.Sandbox.Env.Unset {
		log.Debug().Str("application", app.Application.Bundle).
			Msgf("Unsetting environment variable %s", name)

		c.Arguments = append(c.Arguments, "--unsetenv", name)
	}
}

func IsEnvAllowed(env string) bool {
	var regex = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)
	var badEnvVariables = []string{"LD_PRELOAD", "LD_LIBRARY_PATH"}

	if !regex.MatchString(env) {
		return false
	}

	if slices.Contains(badEnvVariables, env) {
		return false
	}

	return true
}
