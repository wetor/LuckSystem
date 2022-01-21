package czimage

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"lucksystem/utils"

	"github.com/go-restruct/restruct"
)

type CzHeader struct {
	Magic        []byte `struct:"size=4"`
	HeaderLength uint32
	Width        uint16
	Heigth       uint16
	Colorbits    uint16
	Colorblock   uint16
}

type CzData struct {
	Raw        []byte        // Load()
	OutputInfo *CzOutputInfo // Load()
	Image      image.Image   // Export()
	PngImage   image.Image   // Import()

}

type CzBlockInfo struct {
	CompressedSize uint32
	RawSize        uint32
}
type CzOutputInfo struct {
	Offset              int `struct:"-"`
	TotalRawSize        int `struct:"-"`
	TotalCompressedSize int `struct:"-"`
	FileCount           uint32
	BlockInfo           []CzBlockInfo `struct:"size=FileCount"`
}

type CzImage interface {
	Load(header CzHeader, data []byte)
	GetImage() image.Image
	Export(path string)
	Import(file string)
}

func LoadCzImage(data []byte) (CzImage, error) {
	header := CzHeader{}
	err := restruct.Unpack(data[:16], binary.LittleEndian, &header)
	if err != nil {
		utils.Log("restruct.Unpack", err.Error())
		return nil, err
	}
	fmt.Println("cz header", header)
	var cz CzImage
	switch string(header.Magic[:3]) {
	case "CZ1":
		cz = new(Cz1Image)
		cz.Load(header, data)
	//case "CZ2":
	//	cz = new(Cz2Image)
	//	cz.Load(header, data)
	case "CZ3":
		cz = new(Cz3Image)
		cz.Load(header, data)
	default:
		return nil, errors.New("未知类型")
	}

	return cz, nil
}
