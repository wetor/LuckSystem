package operator

import (
	"github.com/golang/glog"
	"lucksystem/charset"
	"lucksystem/game/engine"
	"lucksystem/game/runtime"
)

type SP struct {
	LucaOperateUndefined
	LucaOperateDefault
	LucaOperateExpr
}

func NewSP() *SP {
	return &SP{}
}

func (g *SP) Init(ctx *runtime.Runtime) {
	ctx.Init(charset.ShiftJIS, charset.Unicode, true)
}

func (g *SP) MESSAGE(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	op := NewOP(ctx, code.ParamBytes, 0)

	var voiceId = op.ReadUInt16(true)
	var msgStr = op.ReadLString(true, ctx.TextCharset)
	op.ReadUInt8(false)
	op.SetOperateParams()

	return func() {
		// 这里是执行内容
		ctx.Engine.MESSAGE(voiceId, msgStr)
		ctx.ChanEIP <- 0
	}

}
func (g *SP) SELECT(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	op := NewOP(ctx, code.ParamBytes, 0)

	var varID = op.ReadUInt16(true)
	op.ReadUInt16(false)
	op.ReadUInt16(false)
	op.ReadUInt16(false)
	var msgStr = op.ReadLString(true, ctx.TextCharset)
	op.ReadUInt16(false)
	op.ReadUInt16(false)
	op.ReadUInt16(false)
	op.SetOperateParams()

	return func() {

		selectID := ctx.Engine.SELECT(msgStr)

		//fun.Call(&varID, msgStr_en)
		glog.V(3).Infof("SELECT #%d = %d\n", varID, selectID)
		ctx.ChanEIP <- 0
	}
}

func (g *SP) IMAGELOAD(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	op := NewOP(ctx, code.ParamBytes, 0)

	var mode = op.ReadUInt16(true)
	op.ReadUInt16(true)

	if mode == 1795 {
		op.ReadUInt8(true)
	} else {
		op.ReadUInt16(true)
		op.ReadUInt16(true)
		op.ReadUInt16(true)
	}

	if mode == 1 {
		// 立绘
		op.ReadUInt16(true)
	} else if mode == 1795 {
	} else {
		// 其他
	}
	op.SetOperateParams()
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码

		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}
