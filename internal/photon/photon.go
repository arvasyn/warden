package photon

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/arvasyn/warden/internal/pkg/apperr"
	"github.com/arvasyn/warden/internal/pkg/sandbox"
	"github.com/rs/zerolog/log"
)

func Run(app sandbox.Manifest, userArgs []string, rootDir string) error {
	if app.Type == sandbox.UnitTypeApplication && !app.UseSandbox {
		log.Warn().Type("application", app.Application.Bundle).
			Msg("Tried disabling sandbox")

		return apperr.ErrApplicationNoDisableSandbox
	}

	if app.Application.Exec == "" {
		return apperr.ErrEmptyExec
	}

	config, err := NewConfig()
	if err != nil {
		return err
	}

	if err := config.AddMounts(app); err != nil {
		return err
	}

	config.AddNamespaces(app)
	config.AddEnvironment(app)
	config.AddCapabilities(app)

	if err := config.NewTempDirectory(app.Application.Bundle, "/var/cache/fontconfig"); err != nil {
		return err
	}

	if dbus := os.Getenv("DBUS_SESSION_BUS_ADDRESS"); dbus != "" {
		if after, ok := strings.CutPrefix(dbus, "unix:path="); ok {
			dbus = strings.SplitN(after, ",", 2)[0]
		}

		config.Arguments = append(config.Arguments,
			"--ro-bind", dbus, dbus,
			"--setenv", "DBUS_SESSION_BUS_ADDRESS", dbus,
		)
	}

	if xauth := os.Getenv("XAUTHORITY"); xauth != "" {
		config.Arguments = append(config.Arguments,
			"--ro-bind", xauth, xauth,
			"--setenv", "XAUTHORITY", xauth,
		)
	}

	config.Arguments = append(config.Arguments, "--ro-bind", rootDir, "/app")
	config.Arguments = append(config.Arguments, "--", path.Join("/app", app.Application.Exec))

	if app.Application.AppendUserArgs {
		log.Debug().Type("application", app.Application.Bundle).
			Msgf("Appending user args %s", strings.Join(userArgs, " "))

		config.Arguments = append(config.Arguments, userArgs...)
	}

	cmd := exec.Command("/usr/bin/bwrap", config.Arguments...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Str("output", string(out)).
			Msgf("Error executing command: %v", err)
		return err
	}

	log.Debug().Str("output", string(out)).Msg("Command succeeded")
	return nil
}
