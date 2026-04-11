package gallium

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/arvasyn/warden/internal/pkg/apperr"
	"github.com/rs/zerolog/log"
)

func Run(app Manifest, rootDir string, userArgs []string) error {
	if app.Type == UnitTypeApplication && !app.UseSandbox {
		log.Warn().
			Type("application", app.Application.Bundle).
			Msg("Tried disabling sandbox")

		return apperr.ErrApplicationNoDisableSandbox
	}

	if app.Application.Exec == "" {
		return apperr.ErrEmptyExec
	}

	if strings.Contains(app.Application.Exec, "..") {
		log.Error().
			Str("application", app.Application.Bundle).
			Str("exec", app.Application.Exec).
			Msg("Application tried using path traversal")

		return apperr.ErrAttemptedPathTraversal
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

	socket := "/tmp/dbus-proxy-" + app.Application.Bundle

	proxy, err := config.AddDBusProxy(app, socket)
	if err != nil {
		return err
	}

	defer func() {
		proxy.Process.Kill()
		proxy.Wait()
		os.Remove(socket)
	}()

	if err := config.NewTempDirectory(app.Application.Bundle, "/var/cache/fontconfig"); err != nil {
		return err
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
		log.Debug().
			Type("application", app.Application.Bundle).
			Msgf("Appending user args %s", strings.Join(userArgs, " "))

		config.Arguments = append(config.Arguments, userArgs...)
	}

	cmd := exec.Command("/usr/bin/bwrap", config.Arguments...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: config.UID,
			Gid: config.GID,
		},
	}

	err = cmd.Run()
	if err != nil {
		log.Error().
			Err(err).
			Msgf("Error executing command: %v", err)

		return err
	}

	log.Debug().
		Msg("Command succeeded")

	return nil
}
