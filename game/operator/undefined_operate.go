package operator

import (
	"github.com/golang/glog"
	"lucksystem/game/engine"
	"lucksystem/game/runtime"
)

type LucaOperateUndefined struct {
}

func (g *LucaOperateUndefined) UNDEFINED(ctx *runtime.Runtime, opcode string) engine.HandlerFunc {
	glog.V(5).Infoln(ctx.CIndex, "Operation不存在", opcode)
	code := ctx.Code()
	if len(opcode) == 0 {
		opcode = ToString("%X", code.Opcode)
	}
	list, end := AllToUint16(code.ParamBytes)
	if end >= 0 {
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			list,
			code.ParamBytes[end],
		)
	} else {
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			list,
		)
	}
	return func() {
		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}
