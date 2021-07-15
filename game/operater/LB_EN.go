package operater

import (
	"lucascript/charset"
	"lucascript/game/context"
	"lucascript/game/engine"
	"lucascript/utils"
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

func (g *LB_EN) MESSAGE(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var voiceId uint16
	var msgStr_jp string
	var msgStr_en string

	next := GetParam(code.CodeBytes, &voiceId)
	next = GetParam(code.CodeBytes, &msgStr_jp, next, 0, g.TextCharset)
	GetParam(code.CodeBytes, &msgStr_en, next, 0, g.TextCharset)

	return func() {
		// 这里是执行内容
		ctx.Engine.MESSAGE(voiceId, msgStr_jp)
		ctx.ChanEIP <- 0
	}
}
func (g *LB_EN) SELECT(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var varID uint16
	var var0 uint16
	var var1 uint16
	var var2 uint16
	var msgStr_jp string
	var msgStr_en string

	next := GetParam(code.CodeBytes, &varID)
	next = GetParam(code.CodeBytes, &var0, next)
	next = GetParam(code.CodeBytes, &var1, next)
	next = GetParam(code.CodeBytes, &var2, next)
	next = GetParam(code.CodeBytes, &msgStr_jp, next, 0, g.TextCharset)
	GetParam(code.CodeBytes, &msgStr_en, next, 0, g.TextCharset)

	return func() {

		selectID := ctx.Engine.SELECT(msgStr_jp)

		//fun.Call(&varID, msgStr_en)
		utils.Logf("SELECT #%d = %d", varID, selectID)
		ctx.Variable.Set(ToString("#%d", varID), selectID)
		ctx.ChanEIP <- 0
	}
}
