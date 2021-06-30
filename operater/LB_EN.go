package operater

import (
	"lucascript/charset"
	"lucascript/script"
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

func (g *LB_EN) MESSAGE(code *script.CodeLine) string {
	opcode := "message"
	voiceId := ToUint16(code.CodeBytes[0:2])
	jpStr, next := DecodeString(code.CodeBytes, 2, 0, g.TextCharset)
	enStr, _ := DecodeString(code.CodeBytes, next, 0, g.TextCharset)
	return ToString(`%d:%s (%d, "%s", "%s")`, code.Pos, opcode, voiceId, jpStr, enStr)
}
