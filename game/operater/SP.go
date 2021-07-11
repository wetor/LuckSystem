package operater

import (
	"lucascript/charset"
	"lucascript/function"
	"lucascript/game/context"
	"lucascript/paramter"
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
	var voiceId paramter.LUint16
	var msgStr paramter.LString
	var end paramter.LUint16

	next := GetParam(code.CodeBytes, &voiceId)
	next = GetParam(code.CodeBytes, &msgStr, next, 0, g.TextCharset)
	GetParam(code.CodeBytes, &end, next)

	fun := function.MESSAGE{}
	return func() {
		// 这里是执行内容
		fun.Call([]paramter.Paramter{&voiceId, &msgStr})
		ctx.ChanEIP <- 0
	}

}
