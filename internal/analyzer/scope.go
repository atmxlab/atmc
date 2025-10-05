package semantic

import (
	"github.com/atmxlab/atmcfg/pkg/errors"
)

// Scope Включает в себя контекст.
// Может хранить название переменных (можно дополнить)
type scope struct {
	variables map[string]*variable
}

func newScope() *scope {
	return &scope{variables: make(map[string]*variable)}
}

type variable struct {
	name string
	refs uint
}

func (v *variable) Refs() uint {
	return v.refs
}

func (v *variable) Name() string {
	return v.name
}

func (v *variable) incrRef() {
	v.refs++
}

func (v *variable) hasRefs() bool {
	return v.refs != 0
}

func (s *scope) addVariable(name string) {
	s.variables[name] = &variable{
		name: name,
		refs: 0,
	}
}

func (s *scope) hasVariable(name string) bool {
	_, ok := s.variables[name]
	return ok
}

func (s *scope) incrRef(name string) {
	s.variables[name].incrRef()
}

func (s *scope) checkVariableRefs() error {
	j := errors.NewJoiner()

	for _, v := range s.variables {
		if !v.hasRefs() {
			j.Join(errors.Wrapf(ErrUnusedVariable, "unuse variable: %s", v.name))
		}
	}

	return j.Err()
}
