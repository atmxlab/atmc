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

func (osb *OSBuilder) File(hook func(fb *FileBuilder)) *OSBuilder {
	fb := NewFileBuilder()
	hook(fb)
	path, content := fb.Build()

	osb.contentByFile[path] = content
	return osb
}

func (osb *OSBuilder) Build() OS {
	return OS{
		contentByFile: osb.contentByFile,
	}
}

type FileBuilder struct {
	path    string
	content []byte
}

func NewFileBuilder() *FileBuilder {
	return &FileBuilder{}
}

func (fb *FileBuilder) Path(path string) *FileBuilder {
	fb.path = path
	return fb
}

func (fb *FileBuilder) Content(content string) *FileBuilder {
	fb.content = []byte(content)
	return fb
}

func (fb *FileBuilder) Build() (string, []byte) {
	return fb.path, fb.content
}
