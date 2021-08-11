package czimage

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"lucascript/utils"

	"github.com/go-restruct/restruct"
)

type CzHeader struct {
	Magic        string `struct:"size=4"`
	HeaderLength uint32
	Width        uint16
	Heigth       uint16
	Colorbits    uint16

	Colorblock uint16
	X          uint16
	Y          uint16
	Width1     uint16
	Heigth1    uint16

	Width2  uint16
	Heigth2 uint16
}

type CzBlockInfo struct {
	BlockIndex     uint32
	RawSize        uint32
	CompressedSize uint32
}
type CzOutputInfo struct {
	TotalRawSize        uint32
	TotalCompressedSize uint32
	FileCount           uint32
	BlockInfo           []CzBlockInfo
}

type CzImage interface {
	Load(header *CzHeader, data []byte)
	Get() image.Image
	Save(path string)
}

func LoadCzImage(index int, data []byte) (CzImage, error) {
	header := &CzHeader{}
	err := restruct.Unpack(data, binary.LittleEndian, header)
	if err != nil {
		utils.Log("restruct.Unpack", err.Error())
		return nil, err
	}
	fmt.Println(header)

	var cz CzImage
	switch header.Magic[:3] {
	case "CZ0":
	case "CZ3":
		cz = &Cz3Image{}
		cz.Load(header, data)
	default:
		return nil, errors.New("未知类型")
	}

	return cz, nil
}
