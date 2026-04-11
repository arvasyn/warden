package gallium

import (
	"os"
	"path/filepath"

	"github.com/arvasyn/warden/internal/pkg/apperr"
	"github.com/rs/zerolog/log"
)

func (c *Config) AddMounts(manifest Manifest) error {
	for _, mount := range manifest.Sandbox.Filesystem.Mounts {
		if mount.Source == "" || mount.Target == "" {
			continue
		}

		if mount.Type == MountTypeTmpfs || mount.Type == MountTypeProc || mount.Type == MountTypeDev {
			args, _, err := c.ProcessMountType(mount, manifest, "", "")
			if err != nil {
				return err
			}

			if args != nil {
				c.Arguments = append(c.Arguments, args...)
			}

			continue
		}

		canonicalSource, err := filepath.EvalSymlinks(filepath.Clean(mount.Source))
		if err != nil {
			log.Error().
				Str("source", mount.Source).
				Err(err).
				Msg("Failed to resolve source path")

			continue
		}

		canonicalTarget := filepath.Clean(mount.Target)

		if _, err := os.Lstat(canonicalTarget); err == nil {
			if resolved, err := filepath.EvalSymlinks(canonicalTarget); err == nil {
				canonicalTarget = resolved
			}
		}

		if IsPathBlacklisted(canonicalSource) || IsPathBlacklisted(canonicalTarget) {
			log.Warn().
				Str("application", manifest.Application.Bundle).
				Msgf("Tried mounting a blacklist path: %s -> %s", canonicalSource, canonicalTarget)

			return apperr.ErrBlacklistedPath
		}

		if !Ask(manifest, canonicalSource, manifest.Sandbox.Permissions[mount.Source]) {
			log.Info().
				Str("application", manifest.Application.Bundle).
				Msgf("User denied access to %s", canonicalSource)

			continue
		}

		log.Info().
			Str("application", manifest.Application.Bundle).
			Msgf("User allowed access to %s", canonicalSource)

		args, skip, err := c.ProcessMountType(mount, manifest, canonicalSource, canonicalTarget)
		if err != nil {
			return err
		}

		if skip {
			continue
		}

		c.Arguments = append(c.Arguments, args...)
	}

	return nil
}

func (c *Config) ProcessMountType(mount Mount, app Manifest, canonicalSource, canonicalTarget string) ([]string, bool, error) {
	switch mount.Type {
	case MountTypeBind:
		return c.ProcessBind(mount, app, canonicalSource, canonicalTarget, false)
	case MountTypeROBind:
		return c.ProcessBind(mount, app, canonicalSource, canonicalTarget, true)
	case MountTypeTmpfs:
		c.Arguments = append(c.Arguments, "--tmpfs", "/tmp")
		return nil, true, nil
	case MountTypeProc:
		c.Arguments = append(c.Arguments, "--proc", "/proc")
		c.AllowSharedPID = false
		return nil, true, nil
	case MountTypeDev:
		c.Arguments = append(c.Arguments, "--dev", "/dev")
		return nil, true, nil
	default:
		return nil, true, nil
	}
}

func (c *Config) ProcessBind(mount Mount, app Manifest, canonicalSource, canonicalTarget string, readOnly bool) ([]string, bool, error) {
	permission, ok := app.Sandbox.Permissions[mount.Source]

	if !ok {
		log.Warn().
			Str("application", app.Application.Bundle).
			Msg("Custom file mounts must be declared as a permission")

		return nil, true, nil
	}

	expectedType := PermissionTypeReadWrite
	bindType := "--bind"

	if readOnly {
		expectedType = PermissionTypeRead
		bindType = "--ro-bind"
	}

	if permission.Type != expectedType {
		log.Warn().
			Str("application", app.Application.Bundle).
			Msg("Invalid permission type")

		return nil, true, nil
	}

	return []string{bindType, canonicalSource, canonicalTarget}, false, nil
}
