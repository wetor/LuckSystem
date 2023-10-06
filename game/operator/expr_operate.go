package operator

import (
	"lucksystem/game/engine"
	"lucksystem/game/runtime"
)

type LucaOperateExpr struct {
}

func (g *LucaOperateExpr) EQU(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var key uint16
	var value uint16

	next := GetParam(code.ParamBytes, &key)
	if next < len(code.ParamBytes) {
		GetParam(code.ParamBytes, &value, next)
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			key,
			value,
		)
	} else {
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			key,
		)
	}

	//utils.Logf("EQU #%d = %d", key, value)
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码
		//var keyStr string
		//if key <= 1 {
		//	keyStr = ToString("%d", key)
		//} else {
		//	keyStr = ToString("#%d", key)
		//}
		//ctx.Variable.Set(keyStr, int(value))

		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}

// EQUN 等价于EQU
func (g *LucaOperateExpr) EQUN(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var key uint16
	var value uint16

	next := GetParam(code.ParamBytes, &key)
	if next < len(code.ParamBytes) {
		GetParam(code.ParamBytes, &value, next)
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			key,
			value,
		)
	} else {
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			key,
		)
	}
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码
		//var keyStr string
		//if key <= 1 {
		//	keyStr = ToString("%d", key)
		//} else {
		//	keyStr = ToString("#%d", key)
		//}
		//ctx.Variable.Set(keyStr, int(value))
		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}

func (g *LucaOperateExpr) ADD(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var value uint16
	var exprStr string

	next := GetParam(code.ParamBytes, &value)
	GetParam(code.ParamBytes, &exprStr, next, 0, ctx.ExprCharset)
	ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
		value,
		exprStr,
		ctx.ExprCharset,
	)
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码

		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}
func (g *LucaOperateExpr) RANDOM(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var value uint16
	var lowerStr string
	var upperStr string

	next := GetParam(code.ParamBytes, &value)
	next = GetParam(code.ParamBytes, &lowerStr, next, 0, ctx.ExprCharset)
	GetParam(code.ParamBytes, &upperStr, next, 0, ctx.ExprCharset)
	ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
		value,
		lowerStr,
		upperStr,
		ctx.ExprCharset,
	)
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码

		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}
