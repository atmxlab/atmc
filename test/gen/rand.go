package gen

import (
	"math/rand"

	"github.com/google/uuid"
)

// RandString возвращает случайную строку.
func RandString() string {
	return uuid.NewString()
}

func RandInt() int {
	return rand.Int()
}

func RandUInt() uint {
	return uint(rand.Uint64())
}

// RandInt64 возвращает случайное целое число.
func RandInt64() int64 {
	return rand.Int63()
}

// RandInt32 возвращает случайное целое число.
func RandInt32() int32 {
	return rand.Int31()
}

func RandPositiveInt64() int64 {
	return rand.Int63n(100_000_000_00) + 1_000_00
}

// RandIntInRange возвращает случайное значение из интервала [min, max].
func RandIntInRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// RandElement возвращает рандомный элемент.
func RandElement[T any](elem T, elements ...T) T {
	elements = append(elements, elem)

	return elements[RandIntInRange(0, len(elements)-1)]
}
