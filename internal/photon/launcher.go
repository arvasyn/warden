package photon

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	"github.com/arvasyn/warden/internal/pkg/apperr"
	"github.com/arvasyn/warden/internal/pkg/sandbox"
	"github.com/rs/zerolog/log"
)

var uid, _ = user.Current()

var defaultArgs = []string{
	"--ro-bind", "/usr", "/usr",
	"--ro-bind", "/lib", "/lib",
	"--ro-bind", "/lib64", "/lib64",
	"--ro-bind", "/bin", "/bin",
	"--ro-bind", "/sbin", "/sbin",
	"--tmpfs", "/tmp",
	"--ro-bind", fmt.Sprintf("/run/user/%s/wayland-0", uid.Uid), fmt.Sprintf("/run/user/%s/wayland-0", uid.Uid),
	"--setenv", "XDG_RUNTIME_DIR", fmt.Sprintf("/run/user/%s", uid.Uid),
}

var allowSharedPid = true

func Run(app sandbox.Manifest, userArgs []string, rootDir string) error {
	if app.Type == sandbox.UnitTypeApplication && app.UseSandbox == false {
		log.Warn().Type("application", app.Application.Bundle).Msgf("Tried disabling sandbox")
		return apperr.ErrApplicationNoDisableSandbox
	}

	if app.Application.Exec == "" {
		return apperr.ErrEmptyExec
	}

mountLoop:
	for _, mount := range app.Sandbox.Filesystem.Mounts {
		var args []string

		switch mount.Type {
		case sandbox.MountTypeBind:
			args = []string{"--bind", mount.Source, mount.Target}
		case sandbox.MountTypeROBind:
			args = []string{"--ro-bind", mount.Source, mount.Target}
		case sandbox.MountTypeTmpfs:
			args = []string{"--tmpfs", "/tmp"}
			defaultArgs = append(defaultArgs, args...)
			continue
		case sandbox.MountTypeProc:
			defaultArgs = append(defaultArgs, "--proc", "/proc")
			allowSharedPid = false
			break mountLoop
		case sandbox.MountTypeDev:
			defaultArgs = append(defaultArgs, "--dev", "/dev")
			break mountLoop
		}

		ok := IsPathBlacklisted(mount.Source)
		if !ok {
			log.Warn().Type("application", app.Application.Bundle).Msgf("Tried mounting a blacklist path: %s", mount.Source)
			return apperr.ErrBlacklistedPath
		}

		defaultArgs = append(defaultArgs, args...)
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
		log.Debug().Type("application", app.Application.Bundle).Msgf("Passing environment variable %s", name)

		value, ok := os.LookupEnv(name)
		if !ok {
			log.Warn().Type("application", app.Application.Bundle).Msgf("Unable to get variable %s. Environment variable not set.", name)
			continue
		}

		defaultArgs = append(defaultArgs, "--setenv", name, value)
	}

	for name, value := range app.Sandbox.Env.Set {
		log.Debug().Type("application", app.Application.Bundle).Msgf("Setting environment variable %s to %s", name, value)
		defaultArgs = append(defaultArgs, "--setenv", name, value)
	}

	for _, name := range app.Sandbox.Env.Unset {
		log.Debug().Type("application", app.Application.Bundle).Msgf("Unsetting environment variable %s", name)
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

	defaultArgs = append(defaultArgs, "--ro-bind", rootDir, "/app")
	defaultArgs = append(defaultArgs, "--", path.Join("/app", app.Application.Exec))

	if app.Application.AppendUserArgs == true {
		log.Debug().Type("application", app.Application.Bundle).Msgf("Appending user args %s", strings.Join(defaultArgs, " "))
		defaultArgs = append(defaultArgs, userArgs...)
	}

	cmd := exec.Command("/usr/bin/bwrap", defaultArgs...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Str("output", string(out)).Msgf("Error executing command: %v", err)
		return err
	}

	log.Debug().Str("output", string(out)).Msgf("Command output: %s", string(out))
	return nil
}
