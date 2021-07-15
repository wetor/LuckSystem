package operater

import (
	"fmt"
	"lucascript/charset"
	"lucascript/game/context"
	"lucascript/game/engine"
	"lucascript/game/expr"
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
	UNDEFINE(code *script.CodeLine, opname string) string
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

func (g *LucaOperate) UNDEFINE(code *script.CodeLine, opcode string) string {
	if len(opcode) == 0 {
		opcode = ToString("%X", code.Opcode)
	}
	list, end := AllToUint16(code.CodeBytes)
	str := ""
	for _, num := range list {
		str += fmt.Sprintf(", %d", num)
	}
	if end > 0 {
		str += fmt.Sprintf(", 0x%X", code.CodeBytes[end])
	}
	if len(str) >= 2 {
		str = str[2:]
	}
	return ToString(`%s (%s)`, opcode, str)
}
func (g *LucaOperate) EQU(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var key uint16
	var value uint16

	next := GetParam(code.CodeBytes, &key)
	GetParam(code.CodeBytes, &value, next)

	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码
		var keyStr string
		if key <= 1 {
			keyStr = ToString("%d", key)
		} else {
			keyStr = ToString("#%d", key)
		}
		ctx.Variable.Set(keyStr, int(value))
		// 这里执行与游戏相关代码，内部与虚拟机无关联
		utils.Logf("EQU #%d = %d", key, value)
		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}

// EQUN 等价于EQU
func (g *LucaOperate) EQUN(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var key uint16
	var value uint16

	next := GetParam(code.CodeBytes, &key)
	GetParam(code.CodeBytes, &value, next)

	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码
		var keyStr string
		if key <= 1 {
			keyStr = ToString("%d", key)
		} else {
			keyStr = ToString("#%d", key)
		}
		ctx.Variable.Set(keyStr, int(value))
		// 这里执行与游戏相关代码，内部与虚拟机无关联
		utils.Logf("EQUN #%d = %d", key, value)
		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}

// func (g *LucaOperate) ADD(code *script.CodeLine) string {
// 	opcode := "add"
// 	key := ToUint16(code.CodeBytes[0:2])
// 	exprStr, _ := DecodeString(code.CodeBytes, 2, 0, g.ExprCharset)
// 	return ToString(`%d:%s (#%d, %s)`, code.Pos, opcode, key, exprStr)
// }

func (g *LucaOperate) IFN(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var jumpPos uint32
	var exprStr string
	next := GetParam(code.CodeBytes, &exprStr, 0, 0, g.ExprCharset)
	GetParam(code.CodeBytes, &jumpPos, next, 4)

	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码
		eip := 0

		res, err := expr.RunExpr(exprStr, ctx.Variable.ValueMap)
		if err != nil {
			panic(err)
		}
		if !res {
			utils.Logf("IFN %s => %d", exprStr, !res)
			eip = int(jumpPos)
		}
		// 这里执行与游戏相关代码，内部与虚拟机无关联
		utils.Logf("IFN %s{goto %d}", exprStr, jumpPos)
		ctx.ChanEIP <- eip
	}
}
func (g *LucaOperate) IFY(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var jumpPos uint32
	var exprStr string
	next := GetParam(code.CodeBytes, &exprStr, 0, 0, g.ExprCharset)
	GetParam(code.CodeBytes, &jumpPos, next, 4)

	return func() {
		// 这里是执行 与虚拟机逻辑有关的代码
		eip := 0
		res := true // res:=expr(ifExprStr)
		if res {
			eip = 0
		}
		// 这里执行与游戏相关代码，内部与虚拟机无关联
		utils.Logf("IFN %s{goto %d}", exprStr, jumpPos)
		ctx.ChanEIP <- eip
	}
}
func (g *LucaOperate) FARCALL(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var index uint16
	var fileStr string
	var jumpPos uint32

	next := GetParam(code.CodeBytes, &index)
	next = GetParam(code.CodeBytes, &fileStr, next, 0, g.ExprCharset)
	GetParam(code.CodeBytes, &jumpPos, next)

	return func() {
		// 这里是执行内容
		utils.Logf("FARCALL (%d) {goto \"%s\", %d}", index, fileStr, jumpPos)
		ctx.Engine.FARCALL(index, fileStr, jumpPos)
		ctx.ChanEIP <- 0
	}
}

func (g *LucaOperate) GOTO(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()

	var jumpPos uint32
	GetParam(code.CodeBytes, &jumpPos)

	return func() {
		// 这里是执行内容
		ctx.ChanEIP <- 0
	}
}

func (g *LucaOperate) JUMP(ctx *context.Context) engine.HandlerFunc {
	code := ctx.Code()
	var jumpPos uint32
	var fileStr string
	next := GetParam(code.CodeBytes, &fileStr, 0, 0, g.ExprCharset)
	GetParam(code.CodeBytes, &jumpPos, next, 4)

	return func() {
		// 这里是执行内容
		utils.Logf("JUMP {goto \"%s\", %d}", fileStr, jumpPos)
		ctx.Engine.JUMP(fileStr, jumpPos)
		ctx.ChanEIP <- 0
	}

}
