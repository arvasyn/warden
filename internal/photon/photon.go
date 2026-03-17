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

/*
var uid, _ = user.Current()

var defaultArgs = []string{
	"--bind", fmt.Sprintf("/run/user/%s/wayland-0", uid.Uid), fmt.Sprintf("/run/user/%s/wayland-0", uid.Uid),
	"--ro-bind", "/usr", "/usr",
	"--ro-bind", "/lib", "/lib",
	"--ro-bind", "/lib64", "/lib64",
	"--ro-bind", "/bin", "/bin",
	"--ro-bind", "/sbin", "/sbin",
	"--ro-bind", "/etc/fonts", "/etc/fonts",
	"--tmpfs", "/tmp",
	"--setenv", "XDG_RUNTIME_DIR", fmt.Sprintf("/run/user/%s", uid.Uid),
}

var allowSharedPid = true

func Run(app sandbox.Manifest, userArgs []string, rootDir string) error {
	if app.Type == sandbox.UnitTypeApplication && app.UseSandbox == false {
		log.Warn().Str("application", app.Application.Bundle).Msgf("Tried disabling sandbox")
		return apperr.ErrApplicationNoDisableSandbox
	}

	if app.Application.Exec == "" {
		return apperr.ErrEmptyExec
	}

	for namespace, unshare := range app.Sandbox.Namespaces {
		if unshare {
			defaultArgs = append(defaultArgs, "--unshare-"+namespace)
		} else {
			if namespace == "pid" && allowSharedPid == false {
				defaultArgs = append(defaultArgs, "--share-"+namespace)
			}
		}
	}

	for _, name := range app.Sandbox.Env.Passthrough {
		log.Debug().Str("application", app.Application.Bundle).Msgf("Passing environment variable %s", name)

		value, ok := os.LookupEnv(name)
		if !ok {
			log.Warn().Str("application", app.Application.Bundle).Msgf("Unable to get variable %s. Environment variable not set.", name)
			continue
		}

		defaultArgs = append(defaultArgs, "--setenv", name, value)
	}

	for name, value := range app.Sandbox.Env.Set {
		log.Debug().Str("application", app.Application.Bundle).Msgf("Setting environment variable %s to %s", name, value)
		defaultArgs = append(defaultArgs, "--setenv", name, value)
	}

	for _, name := range app.Sandbox.Env.Unset {
		log.Debug().Str("application", app.Application.Bundle).Msgf("Unsetting environment variable %s", name)
		defaultArgs = append(defaultArgs, "--unsetenv", name)
	}

	for _, c := range app.Sandbox.Capabilities.Add {
		log.Debug().Str("capability", c).Msg("Adding capability")
		defaultArgs = append(defaultArgs, "--cap-add", c)
	}

	if app.Sandbox.Capabilities.DropAll {
		defaultArgs = append(defaultArgs, "--cap-drop", "ALL")
	}

	if dbus := os.Getenv("DBUS_SESSION_BUS_ADDRESS"); dbus != "" {
		socketPath := dbus
		if after, ok := strings.CutPrefix(dbus, "unix:path="); ok {
			socketPath = strings.SplitN(after, ",", 2)[0]
		}

		defaultArgs = append(defaultArgs,
			"--ro-bind", socketPath, socketPath,
			"--setenv", "DBUS_SESSION_BUS_ADDRESS", dbus,
		)
	}

	if xauth := os.Getenv("XAUTHORITY"); xauth != "" {
		defaultArgs = append(defaultArgs,
			"--ro-bind", xauth, xauth,
			"--setenv", "XAUTHORITY", xauth,
		)
	}

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	tmpPath := fmt.Sprintf("/tmp/%s", seededRand.Int())

	err := os.MkdirAll(tmpPath, 0777)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create temporary directory")
	} else {
		defaultArgs = append(defaultArgs, "--bind", tmpPath, "/var/cache/fontconfig")
	}

	defaultArgs = append(defaultArgs, "--ro-bind", rootDir, "/app")
	defaultArgs = append(defaultArgs, "--", path.Join("/app", app.Application.Exec))

	if app.Application.AppendUserArgs == true {
		log.Debug().Str("application", app.Application.Bundle).Msgf("Appending user args %s", strings.Join(defaultArgs, " "))
		defaultArgs = append(defaultArgs, userArgs...)
	}

	cmd := exec.Command("/usr/bin/bwrap", defaultArgs...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Str("output", string(out)).Msgf("Error executing command: %v", err)
		return err
	}

	log.Debug().Str("output", string(out)).Msgf("Command succeeded")
	return nil
}*/

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
