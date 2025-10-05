package testos

import (
	"path/filepath"

	"github.com/atmxlab/atmcfg/pkg/errors"
)

type OS struct {
	contentByFile map[string][]byte
}

func (o OS) ReadFile(s string) ([]byte, error) {
	if content, ok := o.contentByFile[s]; ok {
		return content, nil
	}

	return nil, errors.NotFound("file not found")
}

func (o OS) AbsPath(baseDir, relPath string) (string, error) {
	if filepath.IsAbs(relPath) {
		return relPath, nil
	}

	return filepath.Join(baseDir, relPath), nil
}

type OSBuilder struct {
	contentByFile map[string][]byte
}

func NewOSBuilder() *OSBuilder {
	return &OSBuilder{
		contentByFile: make(map[string][]byte),
	}
}

func (osb *OSBuilder) Content(path string, content []byte) *OSBuilder {
	osb.contentByFile[path] = content
	return osb
}

func (osb *OSBuilder) Build() OS {
	return OS{
		contentByFile: osb.contentByFile,
	}
}
