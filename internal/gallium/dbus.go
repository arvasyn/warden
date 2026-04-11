package gallium

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/arvasyn/warden/internal/pkg/apperr"
	"github.com/rs/zerolog/log"
)

func (c *Config) AddDBusProxy(app Manifest, socket string) (*exec.Cmd, error) {
	hostBusAddress, err := GetDBusSocket(c.UID)
	if err != nil {
		return nil, err
	}

	args := []string{
		hostBusAddress,
		socket,
		"--filter",
	}

	args = append(args, "--talk=org.freedesktop.DBus")

	for _, portal := range app.Sandbox.Portals {
		args = append(args, fmt.Sprintf("--talk=%s", portal))

		if portal == "org.freedesktop.portal.Notification" {
			args = append(args, "--talk=org.freedesktop.Notifications")
		}
	}

	cmd := exec.Command("xdg-dbus-proxy", args...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: c.UID,
			Gid: c.GID,
		},
	}

	if err := cmd.Start(); err != nil {
		log.Error().
			Str("application", app.Application.Bundle).
			Str("exec", app.Application.Exec).
			Err(err).
			Msg("Failed to execute proxy command")

		return nil, apperr.ErrDBusProxyFailed
	}

	if err := WaitForSocket(socket, 3*time.Second); err != nil {
		return nil, apperr.ErrDBusProxyFailed
	}

	c.Arguments = append(c.Arguments,
		"--bind", socket, fmt.Sprintf("/run/user/%d/bus", c.UID),
		"--setenv", "DBUS_SESSION_BUS_ADDRESS", fmt.Sprintf("unix:path=/run/user/%d/bus", c.UID),
	)

	return cmd, nil
}

func WaitForSocket(path string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return nil
		}

		time.Sleep(10 * time.Millisecond)
	}

	return apperr.ErrDBusProxyTimeoutReached
}

func GetDBusSocket(uid uint32) (string, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if _, err := strconv.Atoi(entry.Name()); err != nil {
			continue
		}

		environPath := "/proc/" + entry.Name() + "/environ"

		info, err := os.Stat(environPath)
		if err != nil {
			continue
		}

		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok || stat.Uid != uid {
			continue
		}

		comm, err := os.ReadFile("/proc/" + entry.Name() + "/comm")
		if err != nil {
			continue
		}

		if strings.TrimSpace(string(comm)) != "dbus-daemon" {
			continue
		}

		data, err := os.ReadFile(environPath)
		if err != nil {
			continue
		}

		for _, env := range bytes.Split(data, []byte{0}) {
			if bytes.HasPrefix(env, []byte("DBUS_SESSION_BUS_ADDRESS=")) {
				return string(bytes.TrimPrefix(env, []byte("DBUS_SESSION_BUS_ADDRESS="))), nil
			}
		}
	}

	return "", apperr.ErrFailedToFindDBus
}
