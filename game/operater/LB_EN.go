package operater

import (
	"lucascript/charset"
	"lucascript/game/context"
	"lucascript/game/engine"
	"lucascript/script"
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
	var end uint8

	next := GetParam(code.ParamBytes, &voiceId)
	next = GetParam(code.ParamBytes, &msgStr_jp, next, 0, g.TextCharset)
	next = GetParam(code.ParamBytes, &msgStr_en, next, 0, g.TextCharset)
	GetParam(code.ParamBytes, &end, next)
	ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
		voiceId,
		&script.StringParam{
			Data:   msgStr_jp,
			Coding: g.TextCharset,
		}, &script.StringParam{
			Data:   msgStr_en,
			Coding: g.TextCharset,
		},
		end,
		[]bool{true, true, true, false},
	)
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
	var var3 uint16
	var var4 uint16
	var var5 uint16

	next := GetParam(code.ParamBytes, &varID)
	next = GetParam(code.ParamBytes, &var0, next)
	next = GetParam(code.ParamBytes, &var1, next)
	next = GetParam(code.ParamBytes, &var2, next)
	next = GetParam(code.ParamBytes, &msgStr_jp, next, 0, g.TextCharset)
	next = GetParam(code.ParamBytes, &msgStr_en, next, 0, g.TextCharset)

	next = GetParam(code.ParamBytes, &var3, next)
	next = GetParam(code.ParamBytes, &var4, next)
	GetParam(code.ParamBytes, &var5, next)
	ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
		varID,
		var0,
		var1,
		var2,
		&script.StringParam{
			Data:   msgStr_jp,
			Coding: g.TextCharset,
		}, &script.StringParam{
			Data:   msgStr_en,
			Coding: g.TextCharset,
		},
		var3,
		var4,
		var5,
		[]bool{true, false, false, false, true, true, false, false, false},
	)
	return func() {

		selectID := ctx.Engine.SELECT(msgStr_jp)

		//fun.Call(&varID, msgStr_en)
		utils.Logf("SELECT #%d = %d", varID, selectID)
		ctx.Variable.Set(ToString("#%d", varID), selectID)
		ctx.ChanEIP <- 0
	}
}
