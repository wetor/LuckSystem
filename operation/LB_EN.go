package operation

import (
	"encoding/binary"
	"fmt"
	"lucascript/charset"
	"lucascript/script"
)

type LB_EN struct {
}

func GetLB_EN() Operation {
	return &LB_EN{}
}
func (g *LB_EN) UNDEFINE(code *script.CodeLine, opcode string) string {
	if len(opcode) == 0 {
		opcode = fmt.Sprintf("%X", code.Opcode)
	}
	return fmt.Sprintf(`%d:%s (UnDefine)`, code.Pos, opcode)
}
func (g *LB_EN) EQU(code *script.CodeLine) string {
	opcode := "EQU"
	key := binary.LittleEndian.Uint16(code.CodeBytes[0:2])
	value := binary.LittleEndian.Uint16(code.CodeBytes[2:4])
	return fmt.Sprintf(`%d:%s #%d = %d`, code.Pos, opcode, key, value)
}
func (g *LB_EN) ADD(code *script.CodeLine) string {
	opcode := "ADD"
	key := binary.LittleEndian.Uint16(code.CodeBytes[0:2])
	exprStr, _ := ReadString(code.CodeBytes, 2, charset.UTF_8)
	return fmt.Sprintf(`%d:%s (#%d, %s)`, code.Pos, opcode, key, exprStr)
}
func (g *LB_EN) MESSAGE(code *script.CodeLine) string {
	opcode := "MESSAGE"
	voiceId := binary.LittleEndian.Uint16(code.CodeBytes[0:2])
	jpStr, next := ReadString(code.CodeBytes, 2, charset.Unicode)
	enStr, _ := ReadString(code.CodeBytes, next, charset.Unicode)
	return fmt.Sprintf(`%d:%s (%d, "%s", "%s")`, code.Pos, opcode, voiceId, jpStr, enStr)
}

func (g *LB_EN) IFN(code *script.CodeLine) string {
	opcode := "IFN"
	exprStr, next := ReadString(code.CodeBytes, 0, charset.UTF_8)
	jumpPos := binary.LittleEndian.Uint32(code.CodeBytes[next : next+4])
	return fmt.Sprintf(`%d:%s%s goto %d`, code.Pos, opcode, exprStr, jumpPos)
}
func (g *LB_EN) GOTO(code *script.CodeLine) string {
	opcode := "GOTO"
	jumpPos := binary.LittleEndian.Uint32(code.CodeBytes[0:4])
	return fmt.Sprintf(`%d:%s %d`, code.Pos, opcode, jumpPos)
}
