package prism

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/arvasyn/warden/internal/pkg/apperr"
	"go.yaml.in/yaml/v4"
)

var Config *Configuration = nil

func Configure(anchor string) error {
	splitAnchor := strings.Split(anchor, ":")

	if len(splitAnchor) != 2 {
		return apperr.ErrInvalidAnchorFormat
	}

	switch splitAnchor[0] {
	case "token":
		return apperr.ErrNotImplemented
	case "path":
		file, err := os.Open(splitAnchor[1])
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return apperr.ErrInvalidPath
			}

			return err
		}

		defer file.Close()
		content, err := io.ReadAll(file)

		config, err := parseConfig(content)
		if err != nil {
			return err
		}

		Config = config
		return nil
	default:
		return apperr.ErrInvalidAnchorFormat
	}
}

func parseConfig(content []byte) (*Configuration, error) {
	var cfg Configuration

	err := yaml.Unmarshal(content, &cfg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
