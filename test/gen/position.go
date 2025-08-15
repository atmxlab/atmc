package gen

import (
	"github.com/atmxlab/atmc/test/testutils"
	"github.com/atmxlab/atmc/types"
)

func RandPosition(hooks ...func(pos types.Position)) types.Position {
	pos := types.NewPosition(
		RandUInt(),
		RandUInt(),
		RandUInt(),
	)

	testutils.ApplyHooks(pos, hooks)

	return pos
}
