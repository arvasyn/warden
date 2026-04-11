package gallium

import (
	"fmt"
	"os"
	"os/user"
	"strconv"

	"github.com/arvasyn/warden/internal/pkg/apperr"
	"github.com/rs/zerolog/log"
	"go.yaml.in/yaml/v4"
)

type Config struct {
	UID            uint32
	GID            uint32
	Arguments      []string
	AllowSharedPID bool
}

func NewConfig() (*Config, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	if u.Uid == "0" {
		u.Uid = os.Getenv("SUDO_UID")

		if u.Uid == "" {
			return nil, apperr.ErrInvalidUID
		}
	}

	if u.Gid == "0" {
		u.Gid = os.Getenv("SUDO_GID")

		if u.Gid == "" {
			return nil, apperr.ErrInvalidGID
		}
	}

	uid, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		return nil, err
	}

	gid, err := strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		return nil, err
	}

	return &Config{
		UID: uint32(uid),
		GID: uint32(gid),
		Arguments: []string{
			"--bind", fmt.Sprintf("/run/user/%d/wayland-0", uid), fmt.Sprintf("/run/user/%d/wayland-0", uid),
			"--ro-bind", "/usr", "/usr",
			"--ro-bind", "/lib", "/lib",
			"--ro-bind", "/lib64", "/lib64",
			"--ro-bind", "/bin", "/bin",
			"--ro-bind", "/etc/fonts", "/etc/fonts",
		},
		AllowSharedPID: false,
	}, nil
}

func Parse(path string) (*Manifest, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Unable to read application manifest")

		return nil, apperr.ErrFailedToReadManifest
	}

	manifest := Manifest{}
	err = yaml.Unmarshal(file, &manifest)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Unable to parse application manifest")

		return nil, apperr.ErrFailedToParseManifest
	}

	err = manifest.Validate()
	if err != nil {
		return nil, apperr.ErrFailedToValidateManifest
	}

	return &manifest, nil
}
