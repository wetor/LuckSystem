package operater

import (
	"fmt"
	"lucascript/charset"
	"lucascript/function"
	"lucascript/paramter"
	"lucascript/script"
)

// Operater 需定制指令
type Operater interface {
	MESSAGE(code *script.CodeLine) string
}

// LucaOperater 通用指令
type LucaOperater interface {
	UNDEFINE(code *script.CodeLine, opname string) string
	EQU(code *script.CodeLine) string
	ADD(code *script.CodeLine) string
	IFN(code *script.CodeLine) []paramter.Paramter
	IFY(code *script.CodeLine) string
	GOTO(code *script.CodeLine) string
	JUMP(code *script.CodeLine) string
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
	return ToString(`%d:%s (%s)`, code.Pos, opcode, str)
}
func (g *LucaOperate) EQU(code *script.CodeLine) string {
	opcode := "equ"
	key := ToUint16(code.CodeBytes[0:2])
	value := ToUint16(code.CodeBytes[2:4])
	return ToString(`%d:%s #%d = %d`, code.Pos, opcode, key, value)
}
func (g *LucaOperate) ADD(code *script.CodeLine) string {
	opcode := "add"
	key := ToUint16(code.CodeBytes[0:2])
	exprStr, _ := DecodeString(code.CodeBytes, 2, 0, g.ExprCharset)
	return ToString(`%d:%s (#%d, %s)`, code.Pos, opcode, key, exprStr)
}

// func (g *LucaOperate) IFN(code *script.CodeLine) string {
// 	opcode := "ifN"
// 	exprStr, next := DecodeString(code.CodeBytes, 0, 0, g.ExprCharset)
// 	jumpPos := ToUint32(code.CodeBytes[next : next+4])
// 	return ToString(`%d:%s ("%s", %d)`, code.Pos, opcode, exprStr, jumpPos)
// }
func (g *LucaOperate) IFN(code *script.CodeLine) func() int {
	opcode := "ifNot"
	params := make([]paramter.Paramter, 0, 2)
	exprStr, next := DecodeString(code.CodeBytes, 0, 0, g.ExprCharset)
	params = append(params, &paramter.LString{
		Data:    exprStr,
		Charset: g.ExprCharset,
	})
	jumpPos := ToUint32(code.CodeBytes[next : next+4])
	params = append(params, &paramter.LUint32{
		Data: jumpPos,
	})
	fun := &function.IfNot{
		Name: opcode,
	}
	return func() int {
		// 这里是执行内容
		return fun.Call(params)
	}
}
func (g *LucaOperate) IFY(code *script.CodeLine) string {
	opcode := "ifY"
	exprStr, next := DecodeString(code.CodeBytes, 0, 0, g.ExprCharset)
	jumpPos := ToUint32(code.CodeBytes[next : next+4])
	return ToString(`%d:%s ("%s", %d)`, code.Pos, opcode, exprStr, jumpPos)
}
func (g *LucaOperate) GOTO(code *script.CodeLine) string {
	opcode := "goto"
	jumpPos := ToUint32(code.CodeBytes[0:4])
	return ToString(`%d:%s %d`, code.Pos, opcode, jumpPos)
}

func (g *LucaOperate) JUMP(code *script.CodeLine) string {
	opcode := "jump"
	exprStr, next := DecodeString(code.CodeBytes, 0, 0, g.ExprCharset)
	jumpPos := ToUint32(code.CodeBytes[next : next+4])
	return ToString(`%d:%s ("%s", %d)`, code.Pos, opcode, exprStr, jumpPos)
}
