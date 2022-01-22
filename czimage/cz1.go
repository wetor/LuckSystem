package czimage

import (
	"github.com/golang/glog"
	"image"
	"image/color"
	"image/png"
	"os"
)

// Cz1Image
//  Description Cz1.Load() 载入并解压数据，转化成Image
type Cz1Image struct {
	CzHeader
	ColorPanel [][]byte // []BGRA
	CzData
}

func (cz *Cz1Image) Load(header CzHeader, data []byte) {
	cz.CzHeader = header
	pic := image.NewRGBA(image.Rect(0, 0, int(header.Width), int(header.Heigth)))
	offset := int(cz.HeaderLength)
	switch cz.Colorbits {
	case 4:
		// TODO 未测试
		cz.ColorPanel = make([][]byte, 16)
		for i := 0; i < 16; i++ {
			cz.ColorPanel[i] = data[offset : offset+4]
			offset += 4
		}
		glog.V(6).Infoln("cz1 colorPanel", len(cz.ColorPanel))
		cz.OutputInfo = GetOutputInfo(data[offset:])
		buf := Decompress(data[offset+cz.OutputInfo.Offset:], cz.OutputInfo)
		glog.V(6).Infoln("uncompress size", len(buf))
		i := 0
		var index uint8
		for y := 0; y < int(header.Heigth); y++ {
			for x := 0; x < int(header.Width); x++ {
				if i%2 == 0 {
					index = buf[i/2] & 0x0F // low4bit
				} else {
					index = (buf[i/2] & 0xF0) >> 4 // high4bit
				}
				pic.SetRGBA(x, y, color.RGBA{
					R: cz.ColorPanel[index][2],
					G: cz.ColorPanel[index][1],
					B: cz.ColorPanel[index][0],
					A: cz.ColorPanel[index][3],
				})
				i++
			}
		}
	case 8:
		cz.ColorPanel = make([][]byte, 256)
		for i := 0; i < 256; i++ {
			cz.ColorPanel[i] = data[offset : offset+4]
			offset += 4
		}
		glog.V(6).Infoln("cz1 colorPanel", len(cz.ColorPanel))
		cz.OutputInfo = GetOutputInfo(data[offset:])
		buf := Decompress(data[offset+cz.OutputInfo.Offset:], cz.OutputInfo)
		glog.V(6).Infoln("uncompress size", len(buf))
		// B,G,R,A
		// 0,1,2,3
		i := 0
		for y := 0; y < int(header.Heigth); y++ {
			for x := 0; x < int(header.Width); x++ {
				pic.SetRGBA(x, y, color.RGBA{
					R: cz.ColorPanel[buf[i]][2],
					G: cz.ColorPanel[buf[i]][1],
					B: cz.ColorPanel[buf[i]][0],
					A: cz.ColorPanel[buf[i]][3],
				})
				i++
			}
		}
	case 24:
		// TODO 未测试
		// RGB
		cz.OutputInfo = GetOutputInfo(data[offset:])
		buf := Decompress(data[offset+cz.OutputInfo.Offset:], cz.OutputInfo)
		glog.V(6).Infoln("uncompress size", len(buf))
		i := 0
		for y := 0; y < int(header.Heigth); y++ {
			for x := 0; x < int(header.Width); x++ {
				pic.SetRGBA(x, y, color.RGBA{
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
		cz.OutputInfo = GetOutputInfo(data[offset:])
		buf := Decompress(data[offset+cz.OutputInfo.Offset:], cz.OutputInfo)
		glog.V(6).Infoln("uncompress size", len(buf))
		pic.Pix = buf
	}

	cz.Image = pic
}
func (cz *Cz1Image) Export(path string) {
	f, _ := os.Create(path)
	defer f.Close()
	png.Encode(f, cz.Image)
}

func (cz *Cz1Image) GetImage() image.Image {
	return cz.Image
}
func (cz *Cz1Image) Import(file string) {
	pngFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer pngFile.Close()
	cz.PngImage, err = png.Decode(pngFile)
	if err != nil {
		panic(err)
	}
	pic := cz.PngImage.(*image.NRGBA)
	width := int(cz.Width)
	height := int(cz.Heigth)
	if width != pic.Rect.Size().X || height != pic.Rect.Size().Y {
		glog.V(2).Infof("图片大小不匹配，应该为 w%d h%d\n", width, height)
		return
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
	compressed, info := Compress(data, blockSize)

	cz1File, _ := os.Create(file + ".cz1")
	defer cz1File.Close()
	glog.V(6).Infoln(cz.CzHeader)
	err = WriteStruct(cz1File, &cz.CzHeader, cz.ColorPanel, info)

	if err != nil {
		panic(err)
	}
	cz1File.Write(compressed)
	glog.V(6).Infoln(cz.OutputInfo)
	glog.V(6).Infoln(info)
	cz.OutputInfo.TotalRawSize = 0
	cz.OutputInfo.TotalCompressedSize = 0
	for _, block := range info.BlockInfo {
		cz.OutputInfo.TotalRawSize += int(block.RawSize)
		cz.OutputInfo.TotalCompressedSize += int(block.CompressedSize)
	}
	cz.OutputInfo.Offset = 4 + int(cz.OutputInfo.FileCount)*8
}
