package gallium

import (
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
)

func (c *Config) AddEnvironment(app Manifest) {
	for _, name := range app.Sandbox.Env.Passthrough {
		if !IsEnvAllowed(name) {
			log.Warn().
				Str("application", app.Application.Bundle).
				Str("name", name).
				Msg("Environment variable not allowed")

			continue
		}

		value, ok := os.LookupEnv(name)

		if !ok {
			log.Warn().
				Str("application", app.Application.Bundle).
				Msgf("Unable to get variable %s. Environment variable not set.", name)

			continue
		}

		log.Debug().
			Str("application", app.Application.Bundle).
			Msgf("Passing environment variable %s", name)

		c.Arguments = append(c.Arguments, "--setenv", name, value)
	}

	for name, value := range app.Sandbox.Env.Set {
		if !IsEnvAllowed(name) {
			log.Warn().
				Str("application", app.Application.Bundle).
				Str("name", name).
				Msg("Environment variable not allowed")

			continue
		}

		log.Debug().
			Str("application", app.Application.Bundle).
			Msgf("Setting environment variable %s", name)

		c.Arguments = append(c.Arguments, "--setenv", name, value)
	}

	for _, name := range app.Sandbox.Env.Unset {
		if !IsEnvAllowed(name) {
			log.Warn().
				Str("application", app.Application.Bundle).
				Str("name", name).
				Msg("Environment variable not allowed")

			continue
		}

		log.Debug().
			Str("application", app.Application.Bundle).
			Msgf("Unsetting environment variable %s", name)

		c.Arguments = append(c.Arguments, "--unsetenv", name)
	}
}

func IsEnvAllowed(env string) bool {
	if strings.HasPrefix(env, "LD_") {
		return false
	}

	var regex = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)
	var badEnvVariables = []string{"GLIBC_TUNABLES"}

	if !regex.MatchString(env) {
		return false
	}

	if slices.Contains(badEnvVariables, env) {
		return false
	}

	return true
}
