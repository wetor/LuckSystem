package operater

import (
	"github.com/golang/glog"
	"lucksystem/charset"
	"lucksystem/game/context"
	"lucksystem/game/engine"
	"lucksystem/script"
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
	ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
		voiceId,
		&script.StringParam{
			Data:   string(msgStr),
			Coding: g.TextCharset,
			HasLen: true,
		},
		end,
		[]bool{true, true, false},
	)
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
	var msgStr lstring
	var var3 uint16
	var var4 uint16
	var var5 uint16

	next := GetParam(code.ParamBytes, &varID)
	next = GetParam(code.ParamBytes, &var0, next)
	next = GetParam(code.ParamBytes, &var1, next)
	next = GetParam(code.ParamBytes, &var2, next)
	next = GetParam(code.ParamBytes, &msgStr, next, 0, g.TextCharset)

	next = GetParam(code.ParamBytes, &var3, next)
	next = GetParam(code.ParamBytes, &var4, next)
	GetParam(code.ParamBytes, &var5, next)
	ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
		varID,
		var0,
		var1,
		var2,
		&script.StringParam{
			Data:   string(msgStr),
			Coding: g.TextCharset,
			HasLen: true,
		},
		var3,
		var4,
		var5,
		[]bool{true, false, false, false, true, false, false, false},
	)
	return func() {

		selectID := ctx.Engine.SELECT(msgStr)

		//fun.Call(&varID, msgStr_en)
		glog.V(3).Infof("SELECT #%d = %d\n", varID, selectID)
		ctx.Variable.Set(ToString("#%d", varID), selectID)
		ctx.ChanEIP <- 0
	}
}

func (g *SP) IMAGELOAD(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var mode uint16
	var imgID uint16
	var var1 uint16
	var xPos uint16
	var yPos uint16
	var var2 uint16

	var var3 uint8

	next := GetParam(code.ParamBytes, &mode)
	next = GetParam(code.ParamBytes, &imgID, next)
	if mode == 1795 {
		next = GetParam(code.ParamBytes, &var3, next)
	} else {
		next = GetParam(code.ParamBytes, &var1, next)
		next = GetParam(code.ParamBytes, &xPos, next)
		next = GetParam(code.ParamBytes, &yPos, next)
	}

	if mode == 1 {
		// 立绘
		GetParam(code.ParamBytes, &var2, next)
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			mode,
			imgID,
			var1,
			xPos,
			yPos,
			var2,
		)
	} else if mode == 1795 {
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			mode,
			imgID,
			var3,
		)
	} else {
		// 其他
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			mode,
			imgID,
			var1,
			xPos,
			yPos,
		)
	}

	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码

		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}
