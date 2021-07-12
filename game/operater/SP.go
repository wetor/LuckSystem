package operater

import (
	"lucascript/charset"
	"lucascript/function"
	"lucascript/game/context"
)

type SP struct {
	LucaOperate
}

func GetSP() Operater {
	return &SP{
		LucaOperate: LucaOperate{
			ExprCharset: charset.ShiftJIS,
			TextCharset: charset.Unicode,
		},
	}
}

func (g *SP) MESSAGE(ctx *context.Context) function.HandlerFunc {
	code := ctx.Code()
	var voiceId uint16
	var msgStr string
	var end uint16

	next := GetParam(code.CodeBytes, &voiceId)
	next = GetParam(code.CodeBytes, &msgStr, next, 0, g.TextCharset)
	GetParam(code.CodeBytes, &end, next)

	fun := function.MESSAGE{}
	return func() {
		// 这里是执行内容
		fun.Call(voiceId, msgStr)
		ctx.ChanEIP <- 0
	}

}
func (g *SP) SELECT(ctx *context.Context) function.HandlerFunc {

	return func() {}
}
