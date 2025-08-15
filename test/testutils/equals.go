package testutils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func AssertEmptyDiff[T any](t *testing.T, expected, actual T, opts ...cmp.Option) {
	t.Helper()

	require.Empty(t, cmp.Diff(expected, actual, opts...))
}
