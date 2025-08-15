package testos

import (
	"path/filepath"

	"github.com/atmxlab/atmc/pkg/errors"
)

type OS struct {
	contentByFile map[string][]byte
	env           map[string]string
}

func (o OS) EnvVariables() map[string]string {
	return o.env
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
	env           map[string]string
}

func NewOSBuilder() *OSBuilder {
	return &OSBuilder{
		contentByFile: make(map[string][]byte),
		env:           make(map[string]string),
	}
}

func (osb *OSBuilder) File(hook func(fb *FileBuilder)) *OSBuilder {
	fb := NewFileBuilder()
	hook(fb)
	path, content := fb.Build()

	osb.contentByFile[path] = content
	return osb
}

func (osb *OSBuilder) Env(hook func(eb *EnvBuilder)) *OSBuilder {
	fb := NewEnvBuilder()
	hook(fb)
	key, value := fb.Build()

	osb.env[key] = value
	return osb
}

func (osb *OSBuilder) Build() OS {
	return OS{
		contentByFile: osb.contentByFile,
		env:           osb.env,
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

type EnvBuilder struct {
	key   string
	value string
}

func NewEnvBuilder() *EnvBuilder {
	return &EnvBuilder{}
}

func (fb *EnvBuilder) Key(key string) *EnvBuilder {
	fb.key = key
	return fb
}

func (fb *EnvBuilder) Value(value string) *EnvBuilder {
	fb.value = value
	return fb
}

func (fb *EnvBuilder) Build() (string, string) {
	return fb.key, fb.value
}
