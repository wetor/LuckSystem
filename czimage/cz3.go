package czimage

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"lucksystem/utils"
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

// Cz3Image
//  Description Cz3.Load() 载入数据
//  Description Cz3.Export() 解压数据，转化成Image并导出
//  Description Cz3.GetImage() 解压数据，转化成Image
type Cz3Image struct {
	CzHeader
	Cz3Header
	CzData
}

func (cz *Cz3Image) Load(header CzHeader, data []byte) {
	cz.CzHeader = header
	cz.Raw = data
	err := restruct.Unpack(cz.Raw[16:], binary.LittleEndian, &cz.Cz3Header)
	if err != nil {
		panic(err)
	}
	utils.LogA("cz3 header ", cz.Cz3Header)
	cz.OutputInfo = GetOutputInfo(cz.Raw[int(cz.HeaderLength):])
}

func (cz *Cz3Image) decompress() {
	//os.WriteFile("../data/LB_EN/IMAGE/2.lzw", cz.Raw[int(cz.HeaderLength)+cz.OutputInfo.Offset:], 0666)
	buf := Decompress(cz.Raw[int(cz.HeaderLength)+cz.OutputInfo.Offset:], cz.OutputInfo)
	utils.LogA("uncompress size ", len(buf))
	cz.Image = LineDiff(&cz.CzHeader, buf)
}

func (cz *Cz3Image) Export(path string) {
	if cz.Image == nil {
		cz.decompress()
	}
	f, _ := os.Create(path)
	defer f.Close()
	png.Encode(f, cz.Image)
}

func (cz *Cz3Image) GetImage() image.Image {
	if cz.Image == nil {
		cz.decompress()
	}
	return cz.Image
}

func (cz *Cz3Image) Import(file string) {
	pngFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer pngFile.Close()
	cz.PngImage, err = png.Decode(pngFile)
	if err != nil {
		panic(err)
	}
	data := DiffLine(cz.CzHeader, cz.PngImage)
	compressed, info := Compress(data, 0)

	cz3File, _ := os.Create(file + ".cz3")
	defer cz3File.Close()
	fmt.Println(cz.CzHeader)
	err = WriteStruct(cz3File, &cz.CzHeader, &cz.Cz3Header, info)

	if err != nil {
		panic(err)
	}
	cz3File.Write(compressed)
	fmt.Println(cz.OutputInfo)
	fmt.Println(info)
	cz.OutputInfo.TotalRawSize = 0
	cz.OutputInfo.TotalCompressedSize = 0
	for _, block := range info.BlockInfo {
		cz.OutputInfo.TotalRawSize += int(block.RawSize)
		cz.OutputInfo.TotalCompressedSize += int(block.CompressedSize)
	}
	cz.OutputInfo.Offset = 4 + int(cz.OutputInfo.FileCount)*8
}
