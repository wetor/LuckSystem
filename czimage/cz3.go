package czimage

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/go-restruct/restruct"
)

type Cz3Header struct {
	X       uint16
	Y       uint16
	Width1  uint16
	Heigth1 uint16

	Width2  uint16
	Heigth2 uint16
}

type Cz3Image struct {
	CzHeader
	Cz3Header
	Image    image.Image
	PngImage image.Image
}

func (cz *Cz3Image) Load(header CzHeader, data []byte) {
	cz.CzHeader = header
	err := restruct.Unpack(data[16:], binary.LittleEndian, &cz.Cz3Header)
	if err != nil {
		panic(err)
	}
	fmt.Println("cz3 header", cz.Cz3Header)
	buf := Decompress(data[cz.HeaderLength:])
	fmt.Println("uncompress size", len(buf))
	cz.Image = LineDiff(&cz.CzHeader, buf)
}
func (cz *Cz3Image) Save(path string) {
	f, _ := os.Create(path)
	defer f.Close()
	png.Encode(f, cz.Image)
}

func (cz *Cz3Image) Get() image.Image {
	return cz.Image
}

func (cz *Cz3Image) Import(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	cz.PngImage, err = png.Decode(f)
	if err != nil {
		panic(err)
	}
}
