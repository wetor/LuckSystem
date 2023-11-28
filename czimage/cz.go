package czimage

import (
	"encoding/binary"
	"image"
	"io"
	"os"

	"github.com/go-restruct/restruct"
	"github.com/golang/glog"
)

// CzHeader
//
//	Description 长度为15 byte
type CzHeader struct {
	Magic        []byte `struct:"size=4"`
	HeaderLength uint32
	Width        uint16
	Heigth       uint16
	Colorbits    uint16
	Colorblock   uint8
}

// CzData
//
//	Description cz解析后的结构
type CzData struct {
	Raw        []byte        // Load()
	OutputInfo *CzOutputInfo // Load()
	Image      image.Image   // Export()
	PngImage   image.Image   // Import()

}

// CzBlockInfo
//
//	Description 块大小
type CzBlockInfo struct {
	CompressedSize uint32
	RawSize        uint32
}

// CzOutputInfo
//
//	Description 文件分块信息
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
	Export(w io.Writer) error
	Import(r io.Reader, fillSize bool) error
	Write(w io.Writer) error
}

func LoadCzImageFile(file string) CzImage {
	data, err := os.ReadFile(file)
	if err != nil {
		glog.Fatalln(err)
	}
	return LoadCzImage(data)
}
func LoadCzImage(data []byte) CzImage {
	header := CzHeader{}
	err := restruct.Unpack(data[:15], binary.LittleEndian, &header)
	if err != nil {
		glog.Fatalln(err)
	}
	glog.V(6).Infoln("cz header", header)
	var cz CzImage
	switch string(header.Magic[:3]) {
	case "CZ0":
		cz = new(Cz0Image)
		cz.Load(header, data)
	case "CZ1":
		cz = new(Cz1Image)
		cz.Load(header, data)
	case "CZ2":
		cz = new(Cz2Image)
		cz.Load(header, data)
	case "CZ3":
		cz = new(Cz3Image)
		cz.Load(header, data)
	default:
		glog.Fatalln("Unknown Cz image type")
	}

	return cz
}
