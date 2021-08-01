package operater

import (
	"lucascript/charset"
	"lucascript/game/context"
	"lucascript/game/engine"
	"lucascript/script"
	"lucascript/utils"
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

func (g *SP) MESSAGE(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var voiceId uint16
	var msgStr lstring
	var end uint8

	next := GetParam(code.ParamBytes, &voiceId)
	next = GetParam(code.ParamBytes, &msgStr, next, 0, g.TextCharset)
	GetParam(code.ParamBytes, &end, next)
	ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode, GetOperateName(), voiceId, &script.StringParam{
		Data:   string(msgStr),
		Coding: g.TextCharset,
		HasLen: true,
	}, end, []bool{true, true, false})
	return func() {
		// 这里是执行内容
		ctx.Engine.MESSAGE(voiceId, msgStr)
		ctx.ChanEIP <- 0
	}

}
func (g *SP) SELECT(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var varID uint16
	var var0 uint16
	var var1 uint16
	var var2 uint16
	var msgStr string

	next := GetParam(code.ParamBytes, &varID)
	next = GetParam(code.ParamBytes, &var0, next)
	next = GetParam(code.ParamBytes, &var1, next)
	next = GetParam(code.ParamBytes, &var2, next)
	GetParam(code.ParamBytes, &msgStr, next, 0, g.TextCharset)
	ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode, GetOperateName(), varID, var0, var1, var2, &script.StringParam{
		Data:   msgStr,
		Coding: g.TextCharset,
	}, []bool{true, false, false, false, true})
	return func() {

		selectID := ctx.Engine.SELECT(msgStr)

		//fun.Call(&varID, msgStr_en)
		utils.Logf("SELECT #%d = %d", varID, selectID)
		ctx.Variable.Set(ToString("#%d", varID), selectID)
		ctx.ChanEIP <- 0
	}
}
