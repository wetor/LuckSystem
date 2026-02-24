package czimage

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"github.com/golang/glog"
	"image"
	"image/color"
	"image/png"
	"io"
)

type Cz0Header struct {
	Flag    uint8
	X       uint16
	Y       uint16
	Width1  uint16
	Heigth1 uint16

	Width2  uint16
	Heigth2 uint16
}

// Cz0Image
//
//	Description Cz0.Load() 载入数据
//	Description Cz0.Export() 解压数据，转化成Image并导出
//	Description Cz0.GetImage() 解压数据，转化成Image
type Cz0Image struct {
	CzHeader
	Cz0Header

	// ExtraHeader preserves any extended header bytes between the fixed headers
	// (CzHeader + Cz0Header) and HeaderLength.
	// Some games store positioning/clip metadata there; dropping it can move UI elements.
	ExtraHeader []byte

	CzData
}

// Load
//
//	Description 载入cz图像
//	Receiver cz *Cz0Image
//	Param header CzHeader
//	Param data []byte
func (cz *Cz0Image) Load(header CzHeader, data []byte) {
	cz.CzHeader = header
	cz.Raw = data
	err := restruct.Unpack(cz.Raw[15:], binary.LittleEndian, &cz.Cz0Header)
	if err != nil {
		panic(err)
	}

	// Preserve extended header bytes between the fixed headers and HeaderLength.
	// For AIR Steam CZ0, HeaderLength can be 0x40 (64), while fixed headers are 28 bytes.
	start := 15 + 13
	end := int(cz.HeaderLength)
	if end > len(cz.Raw) {
		end = len(cz.Raw)
	}
	if end > start {
		cz.ExtraHeader = append([]byte(nil), cz.Raw[start:end]...)
	} else {
		cz.ExtraHeader = nil
	}
	glog.V(6).Infoln("cz0 header ", cz.Cz0Header)
	cz.OutputInfo = nil
}

// decompress
//
//	Description 解压数据
//	Receiver cz *Cz0Image
func (cz *Cz0Image) decompress() {
	//os.WriteFile("../data/LB_EN/IMAGE/2.lzw", cz.Raw[int(cz.HeaderLength)+cz.OutputInfo.Offset:], 0666)
	glog.V(6).Infoln("size ", len(cz.Raw))
	glog.V(0).Infof("Decompress CZ0: %dx%d, Colorbits=%d\n",
    cz.Width, cz.Heigth, cz.Colorbits)
	pic := image.NewNRGBA(image.Rect(0, 0, int(cz.Width), int(cz.Heigth)))
	offset := int(cz.HeaderLength)
	switch cz.Colorbits {
	case 32:
		for y := 0; y < int(cz.Heigth); y++ {
			for x := 0; x < int(cz.Width); x++ {
				pic.SetNRGBA(x, y, color.NRGBA{
					R: cz.Raw[offset+0],
					G: cz.Raw[offset+1],
					B: cz.Raw[offset+2],
					A: cz.Raw[offset+3]},
				)
				offset += 4
			}
		}
	}
	cz.Image = pic
}

// GetImage
//
//	Description 取得解压后的图像数据
//	Receiver cz *Cz0Image
//	Return image.Image
func (cz *Cz0Image) GetImage() image.Image {
	if cz.Image == nil {
		cz.decompress()
	}
	return cz.Image
}

// Export
//
//	Description 导出图像到文件
//	Receiver cz *Cz0Image
//	Param w io.Writer
//	Return error
func (cz *Cz0Image) Export(w io.Writer) error {
	if cz.Image == nil {
		cz.decompress()
	}
	return png.Encode(w, cz.Image)
}

// Import
//
//	Description 导入图像
//	Receiver cz *Cz0Image
//	Param r io.Reader
//	Param fillSize bool
//	Return error
func (cz *Cz0Image) Import(r io.Reader, fillSize bool) error {
	var err error
	cz.PngImage, err = png.Decode(r)
	return err

}

// Write
//
//	Description 将图像保存为cz
//	Receiver cz *Cz0Image
//	Param w io.Writer
//	Return error
func (cz *Cz0Image) Write(w io.Writer) error {
	var err error
	glog.V(6).Infoln(cz.CzHeader)
	err = WriteStruct(w, &cz.CzHeader, &cz.Cz0Header)
	if err != nil {
		return err
	}

	// Write preserved extended header bytes and/or pad up to HeaderLength.
	wantExtra := int(cz.HeaderLength) - (15 + 13)
	if wantExtra < 0 {
		wantExtra = 0
	}
	if wantExtra > 0 {
		extra := cz.ExtraHeader
		if len(extra) < wantExtra {
			buf := make([]byte, wantExtra)
			copy(buf, extra)
			extra = buf
		} else if len(extra) > wantExtra {
			extra = extra[:wantExtra]
		}
		if _, err := w.Write(extra); err != nil {
			return err
		}
	}
	pic, ok := cz.PngImage.(*image.NRGBA)
	if !ok {
		pic = ImageToNRGBA(cz.PngImage)
	}
	switch cz.Colorbits {
	case 32:
		_, err = w.Write(pic.Pix)
	}
	return err
}
