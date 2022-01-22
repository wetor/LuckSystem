package operater

import (
	"github.com/golang/glog"
	"lucksystem/charset"
	"lucksystem/game/context"
	"lucksystem/game/engine"
	"lucksystem/script"
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
		glog.V(3).Infof("SELECT #%d = %d\n", varID, selectID)
		ctx.Variable.Set(ToString("#%d", varID), selectID)
		ctx.ChanEIP <- 0
	}
}

func (g *LB_EN) IMAGELOAD(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var mode uint16
	var imgID uint16
	var var1 uint16
	var xPos uint16
	var yPos uint16

	next := GetParam(code.ParamBytes, &mode)
	next = GetParam(code.ParamBytes, &imgID, next)
	if mode == 0 {
		// 背景
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			mode,
			imgID,
		)
	} else if mode == 1 {
		// 立绘
		next = GetParam(code.ParamBytes, &var1, next)
		next = GetParam(code.ParamBytes, &xPos, next)
		GetParam(code.ParamBytes, &yPos, next)
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
func (g *LB_EN) BATTLE(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var battleId uint16
	var msgStr_jp string
	var msgStr_en string
	var var1 uint16

	next := GetParam(code.ParamBytes, &battleId)
	if len(code.ParamBytes) <= next {
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode, battleId, []bool{true})

	} else if battleId == 301 || battleId == 302 {
		GetParam(code.ParamBytes, &var1, next)
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			battleId,
			var1,
			[]bool{true, true},
		)
	} else if battleId == 300 {
		GetParam(code.ParamBytes, &msgStr_jp, next, 0, g.ExprCharset)
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			battleId,
			&script.StringParam{
				Data:   msgStr_jp,
				Coding: g.ExprCharset,
			},
			[]bool{true, true},
		)
	} else {
		next = GetParam(code.ParamBytes, &msgStr_jp, next, 0, g.TextCharset)
		GetParam(code.ParamBytes, &msgStr_en, next, 0, g.TextCharset)
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			battleId,
			&script.StringParam{
				Data:   msgStr_jp,
				Coding: g.TextCharset,
			}, &script.StringParam{
				Data:   msgStr_en,
				Coding: g.TextCharset,
			},
			[]bool{true, true, true},
		)
	}

	return func() {
		// 这里是执行内容
		ctx.Engine.MESSAGE(battleId, msgStr_jp)
		ctx.ChanEIP <- 0
	}
}
