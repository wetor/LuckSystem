package operater

import (
	"lucascript/charset"
	"lucascript/game/context"
	"lucascript/game/engine"
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
	var msgStr string
	var end uint16

	next := GetParam(code.ParamBytes, &voiceId)
	next = GetParam(code.ParamBytes, &msgStr, next, 0, g.TextCharset)
	GetParam(code.ParamBytes, &end, next)
	ctx.Script.AddCodeParams(ctx.CIndex, "MESSAGE", voiceId, msgStr)
	return func() {
		// 这里是执行内容
		ctx.Engine.MESSAGE(voiceId, msgStr)
		ctx.ChanEIP <- 0
	}

}
func (g *SP) SELECT(ctx *context.Context) engine.HandlerFunc {

	return func() {}
}
