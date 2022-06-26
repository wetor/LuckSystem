package czimage

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"github.com/golang/glog"
	"image"
	"image/png"
	"io"
)

type Cz3Header struct {
	Flag    uint8
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

// Load
//  Description 载入cz图像
//  Receiver cz *Cz3Image
//  Param header CzHeader
//  Param data []byte
//
func (cz *Cz3Image) Load(header CzHeader, data []byte) {
	cz.CzHeader = header
	cz.Raw = data
	err := restruct.Unpack(cz.Raw[15:], binary.LittleEndian, &cz.Cz3Header)
	if err != nil {
		panic(err)
	}
	glog.V(6).Infoln("cz3 header ", cz.Cz3Header)
	cz.OutputInfo = GetOutputInfo(cz.Raw[int(cz.HeaderLength):])
}

// decompress
//  Description 解压数据
//  Receiver cz *Cz3Image
//
func (cz *Cz3Image) decompress() {
	//os.WriteFile("../data/LB_EN/IMAGE/2.lzw", cz.Raw[int(cz.HeaderLength)+cz.OutputInfo.Offset:], 0666)
	buf := Decompress(cz.Raw[int(cz.HeaderLength)+cz.OutputInfo.Offset:], cz.OutputInfo)
	glog.V(6).Infoln("uncompress size ", len(buf))
	cz.Image = LineDiff(&cz.CzHeader, buf)
}

// GetImage
//  Description 取得解压后的图像数据
//  Receiver cz *Cz3Image
//  Return image.Image
//
func (cz *Cz3Image) GetImage() image.Image {
	if cz.Image == nil {
		cz.decompress()
	}
	return cz.Image
}

// Export
//  Description 导出图像到文件
//  Receiver cz *Cz3Image
//  Param w io.Writer
//  Return error
//
func (cz *Cz3Image) Export(w io.Writer) error {
	if cz.Image == nil {
		cz.decompress()
	}
	return png.Encode(w, cz.Image)
}

// Import
//  Description 导入图像
//  Receiver cz *Cz3Image
//  Param r io.Reader
//  Param fillSize bool
//  Return error
//
func (cz *Cz3Image) Import(r io.Reader, fillSize bool) error {
	var err error
	cz.PngImage, err = png.Decode(r)
	if err != nil {
		return err
	}
	pic, ok := cz.PngImage.(*image.NRGBA)
	if !ok {
		pic = ImageToNRGBA(cz.PngImage)
	}
	data := DiffLine(cz.CzHeader, pic)
	blockSize := 0
	if len(cz.OutputInfo.BlockInfo) != 0 {
		blockSize = int(cz.OutputInfo.BlockInfo[0].CompressedSize)
	}
	glog.V(6).Infoln(cz.OutputInfo)
	cz.Raw, cz.OutputInfo = Compress(data, blockSize)
	glog.V(6).Infoln(cz.OutputInfo)
	cz.OutputInfo.TotalRawSize = 0
	cz.OutputInfo.TotalCompressedSize = 0
	for _, block := range cz.OutputInfo.BlockInfo {
		cz.OutputInfo.TotalRawSize += int(block.RawSize)
		cz.OutputInfo.TotalCompressedSize += int(block.CompressedSize)
	}
	cz.OutputInfo.Offset = 4 + int(cz.OutputInfo.FileCount)*8

	return nil
}

// Write
//  Description 将图像保存为cz
//  Receiver cz *Cz3Image
//  Param w io.Writer
//  Return error
//
func (cz *Cz3Image) Write(w io.Writer) error {
	var err error
	glog.V(6).Infoln(cz.CzHeader)
	err = WriteStruct(w, &cz.CzHeader, &cz.Cz3Header, cz.OutputInfo)

	if err != nil {
		return err
	}
	_, err = w.Write(cz.Raw)

	return err

}
