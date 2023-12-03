package czimage

import (
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/golang/glog"
)

// Cz1Image
//
//	Description Cz1.Load() 载入并解压数据，转化成Image
type Cz1Image struct {
	CzHeader
	ColorPanel []color.NRGBA // []BGRA
	CzData
}

func (cz *Cz1Image) Load(header CzHeader, data []byte) {
	cz.CzHeader = header
	cz.Raw = data

	offset := int(cz.HeaderLength)
	if cz.Colorbits == 4 || cz.Colorbits == 8 {
		cz.ColorPanel = make([]color.NRGBA, 1<<cz.Colorbits)
		for i := 0; i < (1 << cz.Colorbits); i++ {
			cz.ColorPanel[i] = color.NRGBA{
				B: cz.Raw[offset+0],
				G: cz.Raw[offset+1],
				R: cz.Raw[offset+2],
				A: cz.Raw[offset+3],
			}
			offset += 4
		}
		glog.V(6).Infoln("cz1 colorPanel", len(cz.ColorPanel))
	}
	cz.OutputInfo = GetOutputInfo(cz.Raw[offset:])
}

// decompress
//
//	Description 解压数据
//	Receiver cz *Cz1Image
func (cz *Cz1Image) decompress() {
	pic := image.NewNRGBA(image.Rect(0, 0, int(cz.CzHeader.Width), int(cz.CzHeader.Heigth)))
	offset := int(cz.HeaderLength)
	if cz.Colorbits == 4 || cz.Colorbits == 8 {
		offset += 1 << (cz.Colorbits + 2)
	}
	buf := Decompress(cz.Raw[offset+cz.OutputInfo.Offset:], cz.OutputInfo)
	glog.V(6).Infoln("uncompress size", len(buf))

	switch cz.Colorbits {
	case 4:
		// TODO 未测试
		i := 0
		var index uint8
		for y := 0; y < int(cz.CzHeader.Heigth); y++ {
			for x := 0; x < int(cz.CzHeader.Width); x++ {
				if i%2 == 0 {
					index = buf[i/2] & 0x0F // low4bit
				} else {
					index = (buf[i/2] & 0xF0) >> 4 // high4bit
				}
				pic.SetNRGBA(x, y, cz.ColorPanel[index])
				i++
			}
		}
	case 8:
		i := 0
		for y := 0; y < int(cz.CzHeader.Heigth); y++ {
			for x := 0; x < int(cz.CzHeader.Width); x++ {
				pic.SetNRGBA(x, y, cz.ColorPanel[buf[i]])
				i++
			}
		}
	case 24:
		// TODO 未测试
		// RGB
		i := 0
		for y := 0; y < int(cz.CzHeader.Heigth); y++ {
			for x := 0; x < int(cz.CzHeader.Width); x++ {
				pic.SetNRGBA(x, y, color.NRGBA{
					R: buf[i+0],
					G: buf[i+1],
					B: buf[i+2],
					A: 0xFF,
				})
				i += 3
			}
		}
	case 32:
		// TODO 未测试
		// RGBA
		pic.Pix = buf
	}
	cz.Image = pic
}

// GetImage
//
//	Description 取得解压后的图像数据
//	Receiver cz *Cz1Image
//	Return image.Image
func (cz *Cz1Image) GetImage() image.Image {
	if cz.Image == nil {
		cz.decompress()
	}
	return cz.Image
}

// Export
//
//	Description 导出图像到文件
//	Receiver cz *Cz1Image
//	Param w io.Writer
//	Return error
func (cz *Cz1Image) Export(w io.Writer) error {
	if cz.Image == nil {
		cz.decompress()
	}
	return png.Encode(w, cz.Image)
}

// Import
//
//	Description
//	Receiver cz *Cz1Image
//	Param r io.Reader
//	Param fillSize bool 是否填充大小
//	Return error
func (cz *Cz1Image) Import(r io.Reader, fillSize bool) error {
	var err error
	cz.PngImage, err = png.Decode(r)
	if err != nil {
		panic(err)
	}
	pic := cz.PngImage.(*image.NRGBA)
	width := int(cz.Width)
	height := int(cz.Heigth)
	if fillSize == true {
		// 填充大小
		pic = FillImage(pic, width, height)
	}

	if width != pic.Rect.Size().X || height != pic.Rect.Size().Y {
		glog.V(2).Infof("图片大小不匹配，应该为 w%d h%d\n", width, height)
		return err
	}
	data := make([]byte, width*height)
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			data[i] = pic.At(x, y).(color.NRGBA).A
			i++
		}
	}
	blockSize := 0
	if len(cz.OutputInfo.BlockInfo) != 0 {
		blockSize = int(cz.OutputInfo.BlockInfo[0].CompressedSize)
	}
	cz.Raw, cz.OutputInfo = Compress(data, blockSize)

	cz.OutputInfo.TotalRawSize = 0
	cz.OutputInfo.TotalCompressedSize = 0
	for _, block := range cz.OutputInfo.BlockInfo {
		cz.OutputInfo.TotalRawSize += int(block.RawSize)
		cz.OutputInfo.TotalCompressedSize += int(block.CompressedSize)
	}
	cz.OutputInfo.Offset = 4 + int(cz.OutputInfo.FileCount)*8

	return nil
}

func (cz *Cz1Image) Write(w io.Writer) error {
	var err error
	glog.V(6).Infoln(cz.CzHeader)
	err = WriteStruct(w, &cz.CzHeader, cz.ColorPanel, cz.OutputInfo)

	if err != nil {
		return err
	}
	_, err = w.Write(cz.Raw)

	return err

}
