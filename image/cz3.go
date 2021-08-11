package czimage

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

type Cz3Image struct {
	Header *CzHeader
	Image  image.Image
}

func (cz *Cz3Image) Load(header *CzHeader, data []byte) {
	buf := Decompress(data[header.HeaderLength:])
	fmt.Println(len(buf))
	cz.Header = header
	cz.Image = LineDiff(header, buf)
}
func (cz *Cz3Image) Save(path string) {
	f, _ := os.Create(path)
	defer f.Close()
	png.Encode(f, cz.Image)
}

func (cz *Cz3Image) Get() image.Image {
	return cz.Image
}
