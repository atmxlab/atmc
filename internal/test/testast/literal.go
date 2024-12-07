package testast

import (
	"testing"

	"github.com/atmxlab/atmcfg/internal/parser/ast"
	"github.com/atmxlab/atmcfg/internal/types"
)

func MustNewInt(t *testing.T, str string) ast.Int {
	return MustNewIntWithLocation(t, str, types.Location{})
}

func MustNewFloat(t *testing.T, str string) ast.Float {
	return MustNewFloatWithLocation(t, str, types.Location{})
}

func MustNewBool(t *testing.T, str string) ast.Bool {
	return MustNewBoolWithLocation(t, str, types.Location{})
}

func MustNewIntWithLocation(t *testing.T, str string, loc types.Location) ast.Int {
	i, err := ast.NewInt(str, loc)
	if err != nil {
		t.Fatal(err)
	}

	return i
}

func MustNewBoolWithLocation(t *testing.T, str string, loc types.Location) ast.Bool {
	i, err := ast.NewBool(str, loc)
	if err != nil {
		t.Fatal(err)
	}

	return i
}

func MustNewFloatWithLocation(t *testing.T, str string, loc types.Location) ast.Float {
	i, err := ast.NewFloat(str, loc)
	if err != nil {
		t.Fatal(err)
	}

	return i
}
