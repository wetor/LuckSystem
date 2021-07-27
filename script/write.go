package script

import (
	"encoding/binary"
	"os"

	"github.com/go-restruct/restruct"
)

func (s *ScriptFile) Write() error {
	data, err := restruct.Pack(binary.LittleEndian, s)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.FileName+".out", data, 0666)
	if err != nil {
		return err
	}
	return nil
}
