package photon

import (
	"os"

	"github.com/arvasyn/warden/internal/pkg/sandbox"
	"github.com/rs/zerolog/log"
	"go.yaml.in/yaml/v4"
)

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
