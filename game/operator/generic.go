package operator

import (
	"lucksystem/charset"
	"lucksystem/game/runtime"
)

// Generic is a fallback operator for games without a dedicated handler.
// It embeds the default, undefined, and expr operators to handle
// common opcodes (IFN, IFY, GOTO, JUMP, FARCALL, GOSUB, EQU, etc.)
// Unknown opcodes are handled by UNDEFINED which dumps params as uint16.
type Generic struct {
	LucaOperateUndefined
	LucaOperateDefault
	LucaOperateExpr
}

func NewGeneric() *Generic {
	return &Generic{}
}

func (g *Generic) Init(ctx *runtime.Runtime) {
	ctx.Init(charset.ShiftJIS, charset.Unicode, true)
}
