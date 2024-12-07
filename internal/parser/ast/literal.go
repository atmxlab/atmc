package ast

import (
	"strconv"

	"github.com/atmxlab/atmcfg/pkg/errors"
)

type Int = literalNode[int64]

func NewInt(int string) (Int, error) {
	i, err := strconv.ParseInt(int, 10, 64)
	if err != nil {
		return Int{}, errors.Wrap(err, "error parsing integer")
	}

	return Int{value: i}, nil
}

type Float = literalNode[float64]

func NewFloat(float64 string) (Float, error) {
	fl, err := strconv.ParseFloat(float64, 10)
	if err != nil {
		return Float{}, errors.Wrap(err, "error parsing float")
	}

	return Float{value: fl}, nil
}

type String = literalNode[string]

func NewString(string string) String {
	return String{value: string}
}

type Bool literalNode[bool]

func NewBool(bool string) (Bool, error) {
	switch bool {
	case "true":
		return Bool{value: true}, nil
	case "false":
		return Bool{value: false}, nil
	default:
		return Bool{}, errors.New("invalid bool string")
	}
}
