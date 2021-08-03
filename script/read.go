package script

import (
	"encoding/binary"
	"fmt"
	"lucascript/pak"
	"lucascript/utils"
	"os"

	"github.com/go-restruct/restruct"
)

func (s *ScriptFile) ReadByEntry(entry *pak.FileEntry) error {
	s.FileName = entry.Name
	return s.ReadData(entry.Data)
}
func (s *ScriptFile) Read() error {

	data, err := os.ReadFile(s.FileName)
	if err != nil {
		utils.Log("os.ReadFile", err.Error())
		return err
	}
	return s.ReadData(data)
}

func (s *ScriptFile) ReadData(data []byte) error {
	fmt.Println(len(data))
	err := restruct.Unpack(data, binary.LittleEndian, s)
	if err != nil {
		utils.Log("restruct.Unpack", err.Error())
		// return err
	}
	s.CodeNum = len(s.Codes)
	// s.FormatCodes = make([]string, s.CodeNum)
	pos := 0
	// 预处理 FixedParam
	for i, code := range s.Codes {
		code.Index = i
		code.Pos = pos
		pos += ((int(code.Len) + 1) & ^1) // 向上对齐2
		//(4 + len(code.ParamBytes))
		if code.FixedFlag > 0 {
			if code.FixedFlag >= 2 {
				code.FixedParam = make([]uint16, 2)
				code.FixedParam[0] = binary.LittleEndian.Uint16(code.RawBytes[0:2])
				code.FixedParam[1] = binary.LittleEndian.Uint16(code.RawBytes[2:4])
				code.ParamBytes = make([]byte, len(code.RawBytes)-4)
				copy(code.ParamBytes, code.RawBytes[4:])

			} else {
				code.FixedParam = make([]uint16, 1)
				code.FixedParam[0] = binary.LittleEndian.Uint16(code.RawBytes[0:2])
				code.ParamBytes = make([]byte, len(code.RawBytes)-2)
				copy(code.ParamBytes, code.RawBytes[2:])
			}
		} else {
			code.ParamBytes = make([]byte, len(code.RawBytes))
			copy(code.ParamBytes, code.RawBytes)
		}
	}
	return nil
}
