package adapter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/atmxlab/atmc/pkg/errors"
)

type OS struct {
}

func NewOS() OS {
	return OS{}
}

func (O OS) ReadFile(name string) ([]byte, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		return nil, errors.Wrap(err, "os.ReadFile")
	}

	return content, nil
}

func (O OS) AbsPath(baseDir, relPath string) (string, error) {
	if filepath.IsAbs(relPath) {
		return relPath, nil
	}

	return filepath.Join(baseDir, relPath), nil
}

func (O OS) EnvVariables() map[string]string {
	envMap := make(map[string]string)

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}

	return envMap
}
