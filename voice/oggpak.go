package voice

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"github.com/golang/glog"
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
		glog.V(8).Infoln("restruct.Unpack", err.Error())
		return nil, err
	}
	oggPak.Index = index
	return oggPak, nil
}
