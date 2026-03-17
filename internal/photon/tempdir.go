package photon

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/rs/zerolog/log"
)

func (c *Config) NewTempDirectory(bundle string, target string) error {
	tmpPath := fmt.Sprintf("/tmp/photon/%s/%d", bundle, rand.Int())

	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		log.Error().Err(err).Msg("Failed to create temporary directory")
		return err
	}

	c.Arguments = append(c.Arguments, "--bind", tmpPath, target)
	return nil
}
