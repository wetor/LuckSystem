package operater

import (
	"lucascript/charset"
	"lucascript/script"
)

type SP struct {
	LucaOperate
}

func GetSP() Operater {
	return &SP{
		LucaOperate: LucaOperate{
			ExprCharset: charset.ShiftJIS,
			TextCharset: charset.Unicode,
		},
	}
}

func (g *SP) MESSAGE(code *script.CodeLine) string {
	opcode := "message"
	voiceId := ToUint16(code.CodeBytes[0:2])
	jpStrLen := ToUint16(code.CodeBytes[2:4]) * 2
	jpStr, next := DecodeString(code.CodeBytes, 4, int(jpStrLen), g.TextCharset)
	end := ToUint16(code.CodeBytes[next : next+2])
	return ToString(`%d:%s (%d, "%s", %d)`, code.Pos, opcode, voiceId, jpStr, end)
}
