package script

import (
	"encoding/binary"
	"lucascript/utils"
	"os"

	"github.com/go-restruct/restruct"
)

type CodeLine struct {
	CodeInfo   `struct:"-"`
	Len        uint16
	Opcode     uint8
	FixedFlag  uint8
	FixedParam []uint16 `struct:"-"`
	CodeBytes  []byte   `struct:"size=Len - 4"` //`struct:"size=((Len+ 1)& ^1)- 4"`
	Align      []byte   `struct:"size=Len & 1"`
}

type CodeInfo struct {
	Index      int // 序号
	Pos        int // 文件偏移
	LabelIndex int // 跳转标记，和Pos关联
}

func (s *ScriptFile) Read() error {

	data, err := os.ReadFile(s.FileName)
	if err != nil {
		utils.Log("os.ReadFile", err.Error())
		return err
	}
	err = restruct.Unpack(data, binary.LittleEndian, s)
	if err != nil {
		utils.Log("restruct.Unpack", err.Error())
		// return err
	}
	s.CodeNum = len(s.Codes)

	pos := 0
	// 预处理 FixedParam
	for i, code := range s.Codes {
		code.Index = i
		code.Pos = pos
		pos += ((int(code.Len) + 1) & ^1) // 向上对齐2
		//(4 + len(code.CodeBytes))
		if code.FixedFlag > 0 {
			if code.FixedFlag >= 2 {
				code.FixedParam = make([]uint16, 2)
				code.FixedParam[0] = binary.LittleEndian.Uint16(code.CodeBytes[0:2])
				code.FixedParam[1] = binary.LittleEndian.Uint16(code.CodeBytes[2:4])
				code.CodeBytes = code.CodeBytes[4:]

			} else {
				code.FixedParam = make([]uint16, 1)
				code.FixedParam[0] = binary.LittleEndian.Uint16(code.CodeBytes[0:2])
				code.CodeBytes = code.CodeBytes[2:]
			}
		}
	}
	return nil
}
