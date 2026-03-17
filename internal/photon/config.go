package photon

import (
	"fmt"
	"os"
	"os/user"

	"github.com/arvasyn/warden/internal/pkg/sandbox"
	"github.com/rs/zerolog/log"
	"go.yaml.in/yaml/v4"
)

type Config struct {
	Arguments      []string
	AllowSharedPID bool
}

func NewConfig() (*Config, error) {
	uid, err := user.Current()
	if err != nil {
		return nil, err
	}

	return &Config{
		Arguments: []string{
			"--bind", fmt.Sprintf("/run/user/%s/wayland-0", uid.Uid),
			fmt.Sprintf("/run/user/%s/wayland-0", uid.Uid),
			"--ro-bind", "/usr", "/usr",
			"--ro-bind", "/lib", "/lib",
			"--ro-bind", "/lib64", "/lib64",
			"--ro-bind", "/bin", "/bin",
			"--ro-bind", "/sbin", "/sbin",
			"--ro-bind", "/etc/fonts", "/etc/fonts",
			"--setenv", "XDG_RUNTIME_DIR", fmt.Sprintf("/run/user/%s", uid.Uid),
		},
		AllowSharedPID: true,
	}, nil
}

func Parse(path string) sandbox.Manifest {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to read application manifest")
		return sandbox.Manifest{}
	}

	manifest := sandbox.Manifest{}
	err = yaml.Unmarshal(file, &manifest)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse application manifest")
		return sandbox.Manifest{}
	}

	return manifest
}
