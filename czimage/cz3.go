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
	
	// PATCH YOREMI: S'assurer que l'image est en NRGBA pour avoir 4 bytes/pixel
	var nrgbaImg *image.NRGBA
	if img, ok := cz.Image.(*image.NRGBA); ok {
		nrgbaImg = img
	} else {
		glog.V(0).Infof("Export: Converting image to NRGBA (was %T)\n", cz.Image)
		nrgbaImg = ImageToNRGBA(cz.Image)
	}
	
	glog.V(0).Infof("Export: %dx%d, Colorbits=%d, pix_len=%d bytes\n",
		nrgbaImg.Rect.Dx(), nrgbaImg.Rect.Dy(), cz.CzHeader.Colorbits, len(nrgbaImg.Pix))
	
	// Note: png.Encode() va optimiser en RGB si tous les alpha sont 255
	// C'est un comportement standard de la lib PNG de Go
	// Import() doit gérer ça en reconvertissant RGB → RGBA
	return png.Encode(w, nrgbaImg)
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
	
	// PATCH YOREMI: LOG le format PNG source pour debug
	glog.V(0).Infof("Import: PNG source type=%T, bounds=%v\n", 
		cz.PngImage, cz.PngImage.Bounds())
	
	// FORCER conversion en NRGBA (4 bytes/pixel)
	pic, ok := cz.PngImage.(*image.NRGBA)
	if !ok {
		glog.V(0).Infof("Import: Converting to NRGBA (source was %T)\n", cz.PngImage)
		pic = ImageToNRGBA(cz.PngImage)
	}
	
	// VÉRIFICATION CRITIQUE: pic.Pix doit avoir exactement width * height * 4 bytes
	expectedPixLen := pic.Rect.Dx() * pic.Rect.Dy() * 4
	if len(pic.Pix) != expectedPixLen {
		glog.Errorf("Import: CRITICAL ERROR - Pix length mismatch! expected=%d, got=%d\n",
			expectedPixLen, len(pic.Pix))
	}
	
	// VÉRIFICATION: CzHeader.Colorbits doit être 32 (AIR attend toujours RGBA)
	if cz.CzHeader.Colorbits != 32 {
		glog.Warningf("Import: CzHeader.Colorbits=%d (expected 32), may cause corruption\n",
			cz.CzHeader.Colorbits)
	}
	
	glog.V(0).Infof("Import: Calling DiffLine with Colorbits=%d, pic=%dx%d stride=%d\n",
		cz.CzHeader.Colorbits, pic.Rect.Dx(), pic.Rect.Dy(), pic.Stride)
	
	data := DiffLine(cz.CzHeader, pic)
	
	// VÉRIFICATION: data doit être width * height * (Colorbits/8)
	expectedDataLen := int(cz.Width) * int(cz.Heigth) * int(cz.CzHeader.Colorbits/8)
	if len(data) != expectedDataLen {
		glog.Errorf("Import: CRITICAL ERROR - DiffLine output size mismatch! expected=%d, got=%d\n",
			expectedDataLen, len(data))
	} else {
		glog.V(0).Infof("Import: DiffLine OK, generated %d bytes\n", len(data))
	}
	
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
	
	glog.V(0).Infof("Import: Compressed to %d bytes (RawSize=%d)\n",
		cz.OutputInfo.TotalCompressedSize, cz.OutputInfo.TotalRawSize)

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
	// PATCH YOREMI: Force le Magic à rester "CZ3" au lieu de "CZ0"
	// Le bug causait la conversion CZ3 → CZ0, rendant les fichiers illisibles par le jeu
	cz.CzHeader.Magic = []byte{'C', 'Z', '3', 0}
	
	// PATCH YOREMI: Forcer Colorbits à 32 (AIR attend toujours RGBA)
	if cz.CzHeader.Colorbits != 32 {
		glog.Warningf("Write: Forcing Colorbits from %d to 32\n", cz.CzHeader.Colorbits)
		cz.CzHeader.Colorbits = 32
	}
	
	glog.V(6).Infoln(cz.CzHeader)
	glog.V(0).Infof("Write: CZ3 %dx%d, Colorbits=%d, %d blocks, RawSize=%d\n",
		cz.Width, cz.Heigth, cz.CzHeader.Colorbits, cz.OutputInfo.FileCount,
		cz.OutputInfo.TotalRawSize)
	err = WriteStruct(w, &cz.CzHeader, &cz.Cz3Header, cz.OutputInfo)

	if err != nil {
		return err
	}
	_, err = w.Write(cz.Raw)

	return err

}
