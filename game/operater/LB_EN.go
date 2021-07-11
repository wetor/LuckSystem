package operater

import (
	"lucascript/charset"
	"lucascript/function"
	"lucascript/game/context"
	"lucascript/paramter"
)

type LB_EN struct {
	LucaOperate
}

func GetLB_EN() Operater {
	return &LB_EN{
		LucaOperate: LucaOperate{
			ExprCharset: charset.UTF_8,
			TextCharset: charset.Unicode,
		},
	}
}

func (g *LB_EN) MESSAGE(ctx *context.Context) function.HandlerFunc {
	code := ctx.Code()
	var voiceId paramter.LUint16
	var msgStr_jp paramter.LString
	var msgStr_en paramter.LString

	next := GetParam(code.CodeBytes, &voiceId)
	next = GetParam(code.CodeBytes, &msgStr_jp, next, 0, g.TextCharset)
	GetParam(code.CodeBytes, &msgStr_en, next, 0, g.TextCharset)

	fun := function.MESSAGE{}
	return func() {
		// 这里是执行内容
		fun.Call([]paramter.Paramter{&voiceId, &msgStr_jp})
		fun.Call([]paramter.Paramter{&voiceId, &msgStr_en})
		ctx.ChanEIP <- 0
	}
}
