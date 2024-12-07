package gen

import (
	"github.com/atmxlab/atmcfg/internal/test/testutils"
	"github.com/atmxlab/atmcfg/internal/types"
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
