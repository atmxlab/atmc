package ast

import (
	"strconv"

	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Int struct {
	literalNode
	int64
}

func NewInt(int string) (Int, error) {
	i, err := strconv.ParseInt(int, 10, 64)
	if err != nil {
		return Int{}, errors.Wrap(err, "error parsing integer")
	}

	return Int{int64: i}, nil
}

type Float struct {
	literalNode
	float64
}

func NewFloat(float64 string) (Float, error) {
	fl, err := strconv.ParseFloat(float64, 10)
	if err != nil {
		return Float{}, errors.Wrap(err, "error parsing float")
	}

	return Float{float64: fl}, nil
}

type String struct {
	literalNode
	string
}

func NewString(string string) String {
	return String{string: string}
}

type Bool struct {
	literalNode
	bool
}

func NewBool(bool string) (Bool, error) {
	switch bool {
	case "true":
		return Bool{bool: true}, nil
	case "false":
		return Bool{bool: false}, nil
	default:
		return Bool{}, errors.New("invalid bool string")
	}
}
