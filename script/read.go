package script

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/go-restruct/restruct"
)

type ScriptFile struct {
	FileName string      `struct:"-"`
	Version  uint8       `struct:"-"`
	CodeNum  int         `struct:"-"`
	Code     []*CodeLine `struct:"while=true"`
}

type CodeLine struct {
	Pos       int `struct:"-"`
	Len       uint16
	Opcode    uint8
	InfoFlag  uint8
	Info      []uint16 `struct:"-"`
	CodeBytes []byte   `struct:"size=((Len+ 1)& ^1)- 4"`
}

func (s *ScriptFile) Read() error {

	data, err := os.ReadFile(s.FileName)
	if err != nil {
		fmt.Println("os.ReadFile", err.Error())
		return err
	}
	err = restruct.Unpack(data, binary.LittleEndian, s)
	if err != nil {
		fmt.Println("restruct.Unpack", err.Error())
		// return err
	}
	pos := 0
	s.CodeNum = len(s.Code)
	for _, code := range s.Code {
		code.Pos = pos
		pos += (4 + len(code.CodeBytes))
		if code.InfoFlag > 0 {
			if code.InfoFlag >= 2 {
				code.Info = make([]uint16, 2)
				code.Info[0] = binary.LittleEndian.Uint16(code.CodeBytes[0:2])
				code.Info[1] = binary.LittleEndian.Uint16(code.CodeBytes[2:4])
				code.CodeBytes = code.CodeBytes[4:]

			} else {
				code.Info = make([]uint16, 1)
				code.Info[0] = binary.LittleEndian.Uint16(code.CodeBytes[0:2])
				code.CodeBytes = code.CodeBytes[2:]
			}
		}
	}
	return nil
}
