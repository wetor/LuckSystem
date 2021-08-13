package operater

import (
	"lucksystem/charset"
	"lucksystem/game/context"
	"lucksystem/game/engine"
	"lucksystem/script"
	"lucksystem/utils"
)

// Operater 需定制指令
type Operater interface {
	MESSAGE(ctx *context.Context) engine.HandlerFunc
	SELECT(ctx *context.Context) engine.HandlerFunc

	IMAGELOAD(ctx *context.Context) engine.HandlerFunc
}

// LucaOperater 通用指令
type LucaOperater interface {
	UNDEFINE(ctx *context.Context, opname string) engine.HandlerFunc
	UNKNOW0(ctx *context.Context) engine.HandlerFunc
	EQU(ctx *context.Context) engine.HandlerFunc
	EQUN(ctx *context.Context) engine.HandlerFunc
	ADD(ctx *context.Context) engine.HandlerFunc
	RANDOM(ctx *context.Context) engine.HandlerFunc
	IFN(ctx *context.Context) engine.HandlerFunc
	IFY(ctx *context.Context) engine.HandlerFunc
	GOTO(ctx *context.Context) engine.HandlerFunc
	JUMP(ctx *context.Context) engine.HandlerFunc
	FARCALL(ctx *context.Context) engine.HandlerFunc

	MOVE(ctx *context.Context) engine.HandlerFunc
}

// LucaOperate 通用指令
type LucaOperate struct {
	ExprCharset charset.Charset
	TextCharset charset.Charset
	LabelMap    map[int]int
}

func (g *LucaOperate) UNDEFINE(ctx *context.Context, opcode string) engine.HandlerFunc {
	code := ctx.Code()
	if len(opcode) == 0 {
		opcode = ToString("%X", code.Opcode)
	}
	list, end := AllToUint16(code.ParamBytes)
	if end > 0 {
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			list,
			code.ParamBytes[end],
			opcode,
		)
	} else {
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			list,
			opcode,
		)
	}
	return func() {
		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}

func (g *LucaOperate) UNKNOW0(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var value uint16
	var exprStr string

	next := GetParam(code.ParamBytes, &value)
	GetParam(code.ParamBytes, &exprStr, next, 0, g.ExprCharset)
	ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
		value,
		exprStr,
		g.ExprCharset,
	)
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码

		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}

func (g *LucaOperate) EQU(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var key uint16
	var value uint16

	next := GetParam(code.ParamBytes, &key)
	if next < len(code.ParamBytes) {
		GetParam(code.ParamBytes, &value, next)
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			key,
			value,
		)
	} else {
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			key,
		)
	}

	//utils.Logf("EQU #%d = %d", key, value)
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码
		var keyStr string
		if key <= 1 {
			keyStr = ToString("%d", key)
		} else {
			keyStr = ToString("#%d", key)
		}
		ctx.Variable.Set(keyStr, int(value))

		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}

// EQUN 等价于EQU
func (g *LucaOperate) EQUN(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var key uint16
	var value uint16

	next := GetParam(code.ParamBytes, &key)
	if next < len(code.ParamBytes) {
		GetParam(code.ParamBytes, &value, next)
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			key,
			value,
		)
	} else {
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			key,
		)
	}
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码
		var keyStr string
		if key <= 1 {
			keyStr = ToString("%d", key)
		} else {
			keyStr = ToString("#%d", key)
		}
		ctx.Variable.Set(keyStr, int(value))
		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}

func (g *LucaOperate) ADD(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var value uint16
	var exprStr string

	next := GetParam(code.ParamBytes, &value)
	GetParam(code.ParamBytes, &exprStr, next, 0, g.ExprCharset)
	ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
		value,
		exprStr,
		g.ExprCharset,
	)
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码

		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}
func (g *LucaOperate) RANDOM(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var value uint16
	var lowerStr string
	var upperStr string

	next := GetParam(code.ParamBytes, &value)
	next = GetParam(code.ParamBytes, &lowerStr, next, 0, g.ExprCharset)
	GetParam(code.ParamBytes, &upperStr, next, 0, g.ExprCharset)
	ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
		value,
		lowerStr,
		upperStr,
		g.ExprCharset,
	)
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码

		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}

func (g *LucaOperate) IFN(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var exprStr string
	var jumpPos uint32

	next := GetParam(code.ParamBytes, &exprStr, 0, 0, g.ExprCharset)
	GetParam(code.ParamBytes, &jumpPos, next, 4)
	ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
		exprStr,
		&script.JumpParam{
			Position: int(jumpPos),
		},
		g.ExprCharset,
	)
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码
		eip := 0

		res, err := ctx.Variable.TestExpr(exprStr)
		if err != nil {
			panic(err)
		}
		if !res {
			utils.Logf("IFN %s => %d", exprStr, !res)
			eip = int(jumpPos)
		}
		// 这里执行与游戏相关代码，内部与虚拟机无关联

		ctx.ChanEIP <- eip
	}
}
func (g *LucaOperate) IFY(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var exprStr string
	var jumpPos uint32

	next := GetParam(code.ParamBytes, &exprStr, 0, 0, g.ExprCharset)
	GetParam(code.ParamBytes, &jumpPos, next, 4)
	ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
		exprStr,
		&script.JumpParam{
			Position: int(jumpPos),
		},
		g.ExprCharset,
	)
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码
		eip := 0
		res := true // res:=expr(ifExprStr)
		if res {
			eip = 0
		}
		// 这里执行与游戏相关代码，内部与虚拟机无关联

		ctx.ChanEIP <- eip
	}
}
func (g *LucaOperate) FARCALL(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var index uint16
	var fileStr string
	var jumpPos uint32

	next := GetParam(code.ParamBytes, &index)
	next = GetParam(code.ParamBytes, &fileStr, next, 0, g.ExprCharset)
	GetParam(code.ParamBytes, &jumpPos, next)
	// ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode, index, &script.JumpParam{
	// 	ScriptName: fileStr,
	// 	Position:   int(jumpPos),
	// })
	// 文件外跳转
	ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
		index,
		fileStr,
		jumpPos,
		g.ExprCharset,
	)
	return func() {
		// 这里是执行内容

		ctx.Engine.FARCALL(index, fileStr, jumpPos)
		ctx.ChanEIP <- 0
	}
}

func (g *LucaOperate) GOTO(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()

	var jumpPos uint32
	GetParam(code.ParamBytes, &jumpPos)
	ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
		&script.JumpParam{
			Position: int(jumpPos),
		},
	)
	return func() {
		// 这里是执行内容
		ctx.ChanEIP <- 0
	}
}

func (g *LucaOperate) JUMP(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var fileStr string
	var jumpPos uint32

	next := GetParam(code.ParamBytes, &fileStr, 0, 0, g.ExprCharset)
	if next < len(code.ParamBytes) {
		GetParam(code.ParamBytes, &jumpPos, next, 4)
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			fileStr,
			jumpPos,
			g.ExprCharset,
		)
	} else {
		ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
			fileStr,
			g.ExprCharset,
		)
	}
	// ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode, &script.JumpParam{
	// 	ScriptName: fileStr,
	// 	Position:   int(jumpPos),
	// })
	return func() {
		// 这里是执行内容

		ctx.Engine.JUMP(fileStr, jumpPos)
		ctx.ChanEIP <- 0
	}

}

func (g *LucaOperate) MOVE(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var val1 uint8
	var val2 uint16
	var val3 uint16
	var height uint16
	var width uint16

	next := GetParam(code.ParamBytes, &val1)
	next = GetParam(code.ParamBytes, &val2, next)
	next = GetParam(code.ParamBytes, &val3, next)
	next = GetParam(code.ParamBytes, &height, next)
	GetParam(code.ParamBytes, &width, next)
	ctx.Script().SetOperateParams(ctx.CIndex, ctx.RunMode,
		val1,
		val2,
		val3,
		height,
		width,
	)
	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码

		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}
