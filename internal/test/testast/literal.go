package testast

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/types"
)

func MustNewInt(t *testing.T, str string) ast.Int {
	i, err := ast.NewInt(str)
	if err != nil {
		t.Fatal(err)
	}

	return i
}

func MustNewFloat(t *testing.T, str string) ast.Float {
	i, err := ast.NewFloat(str)
	if err != nil {
		t.Fatal(err)
	}

	return i
}

func MustNewBool(t *testing.T, str string) ast.Bool {
	i, err := ast.NewBool(str, types.NewInitialLocation())
	if err != nil {
		t.Fatal(err)
	}

	return i
}
