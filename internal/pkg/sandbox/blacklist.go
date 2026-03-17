package sandbox

import (
	"slices"
	"strings"
)

func IsPathBlacklisted(path string) bool {
	badPaths := []string{
		"/etc/passwd",
		"/etc/shadow",
		"/etc/group",
		"/etc/sudoers",
		"/etc/hosts",
		"/etc/hostname",
		"/etc/ssh/sshd_config",
		"/etc/ssh/ssh_config",
		"/root",
		"/var/log/auth.log",
		"/var/log/syslog",
		"/var/log/secure",
		"/proc/self/environ",
		"/proc/self/cmdline",
		"/proc/mounts",
		"/proc/net/tcp",
		"/boot/grub/grub.cfg",
		"/etc/crontab",
		"/etc/cron.d",
		"/etc/cron.daily",
	}

	if path == "" {
		return false
	}

	trimmed, ok := strings.CutPrefix(path, "/")

	if !ok {
		return true
	}

	splitPath := strings.Split(trimmed, "/")

	if len(splitPath) < 2 {
		return true
	}

	if slices.Contains(badPaths, path) {
		return true
	}

	switch splitPath[0] {
	case "dev":
		return true
	case "core":
		return true
	}

	return false
}
