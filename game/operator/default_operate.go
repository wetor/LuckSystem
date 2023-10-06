package operator

import (
	"lucksystem/game/engine"
	"lucksystem/game/runtime"
)

type LucaOperateDefault struct {
}

//func (g *LucaOperateDefault) UNKNOW0(ctx *runtime.Runtime) engine.HandlerFunc {
//	code := ctx.Code()
//	var value uint16
//	var exprStr string
//
//	next := GetParam(code.ParamBytes, &value)
//	GetParam(code.ParamBytes, &exprStr, next, 0, ctx.ExprCharset)
//	ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
//		value,
//		exprStr,
//		ctx.ExprCharset,
//	)
//	return func() {
//		// 这里是执行 与虚拟机逻辑有关的代码
//
//		// 下一步执行地址，为0则表示紧接着向下
//		ctx.ChanEIP <- 0
//	}
//}

func (g *LucaOperateDefault) IFN(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()

	op := NewOP(ctx, code.ParamBytes, 0)
	op.ReadString(true, ctx.ExprCharset)
	op.ReadJump(true)
	op.SetOperateParams()

	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码
		//eip := 0
		//
		//res, err := ctx.Variable.TestExpr(exprStr)
		//if err != nil {
		//	panic(err)
		//}
		//if !res {
		//	glog.V(3).Infof("IFN %s => %d\n", exprStr, !res)
		//	eip = int(jumpPos)
		//}
		// 这里执行与游戏相关代码，内部与虚拟机无关联

		ctx.ChanEIP <- 0
	}
}
func (g *LucaOperateDefault) IFY(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	op := NewOP(ctx, code.ParamBytes, 0)
	op.ReadString(true, ctx.ExprCharset)
	op.ReadJump(true)
	op.SetOperateParams()

	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码
		// 这里执行与游戏相关代码，内部与虚拟机无关联

		ctx.ChanEIP <- 0
	}
}
func (g *LucaOperateDefault) FARCALL(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()

	op := NewOP(ctx, code.ParamBytes, 0)
	index := op.ReadUInt16(true)
	fileStr := op.ReadString(true, ctx.ExprCharset)
	jumpPos := op.ReadFileJump(true, fileStr)
	op.SetOperateParams()

	return func() {
		// 这里是执行内容

		ctx.Engine.FARCALL(index, fileStr, jumpPos)
		ctx.ChanEIP <- 0
	}
}

func (g *LucaOperateDefault) GOTO(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	op := NewOP(ctx, code.ParamBytes, 0)
	op.ReadJump(true)
	op.SetOperateParams()

	return func() {
		// 这里是执行内容
		ctx.ChanEIP <- 0
	}
}

func (g *LucaOperateDefault) GOSUB(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	op := NewOP(ctx, code.ParamBytes, 0)
	op.ReadUInt16(true)
	op.ReadJump(true)
	op.SetOperateParams()

	return func() {
		// 这里是执行内容
		ctx.ChanEIP <- 0
	}
}

func (g *LucaOperateDefault) JUMP(ctx *runtime.Runtime) engine.HandlerFunc {
	code := ctx.Code()
	var fileStr string
	var jumpPos uint32
	op := NewOP(ctx, code.ParamBytes, 0)
	fileStr = op.ReadString(true, ctx.ExprCharset)
	if op.CanRead() {
		jumpPos = op.ReadFileJump(true, fileStr)
	}
	op.SetOperateParams()

	// ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode, &script.JumpParam{
	// 	ScriptName: fileStr,
	// 	Position:   int(jumpPos),
	// })
	return func() {
		// 这里是执行内容

		ctx.Engine.JUMP(fileStr, jumpPos)
		ctx.ChanEIP <- 0
	}

}

//func (g *LucaOperateDefault) MOVE(ctx *runtime.Runtime) engine.HandlerFunc {
//	code := ctx.Code()
//	var val1 uint8
//	var val2 uint16
//	var val3 uint16
//	var height uint16
//	var width uint16
//
//	next := GetParam(code.ParamBytes, &val1)
//	next = GetParam(code.ParamBytes, &val2, next)
//	next = GetParam(code.ParamBytes, &val3, next)
//	next = GetParam(code.ParamBytes, &height, next)
//	GetParam(code.ParamBytes, &width, next)
//	ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
//		val1,
//		val2,
//		val3,
//		height,
//		width,
//	)
//	return func() {
//		// 这里是执行 与虚拟机逻辑有关的代码
//
//		// 下一步执行地址，为0则表示紧接着向下
//		ctx.ChanEIP <- 0
//	}
//}
