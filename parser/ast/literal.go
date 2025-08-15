package ast

import (
	"strconv"

	"github.com/atmxlab/atmc/pkg/errors"
	"github.com/atmxlab/atmc/types"
)

type Int = literalNode[int64]

func NewInt(int string, loc types.Location) (Int, error) {
	value, err := strconv.ParseInt(int, 10, 64)
	if err != nil {
		return Int{}, errors.Wrap(err, "error parsing integer")
	}

	i := Int{value: value}
	i.loc = loc

	return i, nil
}

type Float = literalNode[float64]

func NewFloat(float64 string, loc types.Location) (Float, error) {
	value, err := strconv.ParseFloat(float64, 10)
	if err != nil {
		return Float{}, errors.Wrap(err, "error parsing float")
	}

	f := Float{value: value}
	f.loc = loc

	return f, nil
}

type String = literalNode[string]

func NewString(string string, loc types.Location) String {
	s := String{value: string}
	s.loc = loc

	return s
}

type Bool = literalNode[bool]

func NewBool(bool string, loc types.Location) (Bool, error) {
	switch bool {
	case "true":
		b := Bool{value: true}
		b.loc = loc

		return b, nil
	case "false":
		b := Bool{value: false}
		b.loc = loc

		return b, nil
	default:
		return Bool{}, errors.New("invalid bool string")
	}
}
