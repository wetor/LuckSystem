package operater

import (
	"lucascript/charset"
	"lucascript/game/context"
	"lucascript/game/engine"
	"lucascript/script"
	"lucascript/utils"
)

// Operater 需定制指令
type Operater interface {
	MESSAGE(ctx *context.Context) engine.HandlerFunc
	SELECT(ctx *context.Context) engine.HandlerFunc
}

// LucaOperater 通用指令
type LucaOperater interface {
	UNDEFINE(ctx *context.Context, opname string) engine.HandlerFunc
	EQU(ctx *context.Context) engine.HandlerFunc
	EQUN(ctx *context.Context) engine.HandlerFunc
	// ADD(code *script.CodeLine) string
	IFN(ctx *context.Context) engine.HandlerFunc
	IFY(ctx *context.Context) engine.HandlerFunc
	GOTO(ctx *context.Context) engine.HandlerFunc
	JUMP(ctx *context.Context) engine.HandlerFunc
	FARCALL(ctx *context.Context) engine.HandlerFunc
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
		ctx.Script.AddCodeParams(ctx.CIndex, opcode, list, code.ParamBytes[end])
	} else {
		ctx.Script.AddCodeParams(ctx.CIndex, opcode, list)
	}
	return func() {
		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}
func (g *LucaOperate) EQU(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var key uint16
	var value uint16

	next := GetParam(code.ParamBytes, &key)
	GetParam(code.ParamBytes, &value, next)
	// 这里执行与游戏相关代码，内部与虚拟机无关联
	ctx.Script.AddCodeParams(ctx.CIndex, "EQU", key, value)
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
	GetParam(code.ParamBytes, &value, next)
	// 这里执行与游戏相关代码，内部与虚拟机无关联
	ctx.Script.AddCodeParams(ctx.CIndex, "EQUN", key, value)
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

// func (g *LucaOperate) ADD(code *script.CodeLine) string {
// 	opcode := "add"
// 	key := ToUint16(code.ParamBytes[0:2])
// 	exprStr, _ := DecodeString(code.ParamBytes, 2, 0, g.ExprCharset)
// 	return ToString(`%d:%s (#%d, %s)`, code.Pos, opcode, key, exprStr)
// }

func (g *LucaOperate) IFN(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var exprStr string
	var jumpPos uint32

	next := GetParam(code.ParamBytes, &exprStr, 0, 0, g.ExprCharset)
	GetParam(code.ParamBytes, &jumpPos, next, 4)
	ctx.Script.AddCodeParams(ctx.CIndex, "IFN", exprStr, &script.JumpParam{
		Position: int(jumpPos),
	})
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
	ctx.Script.AddCodeParams(ctx.CIndex, "IFY", exprStr, &script.JumpParam{
		Position: int(jumpPos),
	})
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
	// ctx.Script.AddCodeParams(ctx.CIndex, "FARCALL", index, &script.JumpParam{
	// 	ScriptName: fileStr,
	// 	Position:   int(jumpPos),
	// })
	// 文件外跳转
	ctx.Script.AddCodeParams(ctx.CIndex, "FARCALL", index, fileStr, jumpPos)
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
	ctx.Script.AddCodeParams(ctx.CIndex, "GOTO", &script.JumpParam{
		Position: int(jumpPos),
	})
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
	GetParam(code.ParamBytes, &jumpPos, next, 4)
	// ctx.Script.AddCodeParams(ctx.CIndex, "JUMP", &script.JumpParam{
	// 	ScriptName: fileStr,
	// 	Position:   int(jumpPos),
	// })
	ctx.Script.AddCodeParams(ctx.CIndex, "JUMP", fileStr, jumpPos)
	return func() {
		// 这里是执行内容

		ctx.Engine.JUMP(fileStr, jumpPos)
		ctx.ChanEIP <- 0
	}

}
