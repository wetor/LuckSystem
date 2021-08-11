package voice

import (
	"encoding/binary"
	"lucascript/utils"

	"github.com/go-restruct/restruct"
)

type OggPak struct {
	Magic []byte     `struct:"size=7"`
	Files []*OggFile `struct:"while=!_eof"`
	Index int        `struct:"-"`
}
type OggFile struct {
	SampleRate uint32
	Length     uint32
	Data       []byte `struct:"size=Length"`
}

func LoadOggPak(index int, data []byte) (*OggPak, error) {
	oggPak := &OggPak{}
	err := restruct.Unpack(data, binary.LittleEndian, oggPak)
	if err != nil {
		utils.Log("restruct.Unpack", err.Error())
		return nil, err
	}
	oggPak.Index = index
	return oggPak, nil
}
