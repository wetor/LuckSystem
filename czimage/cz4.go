package czimage

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"github.com/golang/glog"
	"image"
	"image/png"
	"io"
)

// Cz4Image
//  Description CZ4 format: LZW-compressed with separated RGB/Alpha channels
//  and per-channel delta line encoding.
//  Same header structure as CZ3, same LZW compression, but pixels are stored
//  as [RGB: w*h*3 bytes][Alpha: w*h bytes] instead of interleaved RGBA.
type Cz4Image struct {
	CzHeader
	Cz3Header // CZ4 uses the same extended header as CZ3
	CzData
}

// Load
//  Description Load CZ4 image data
//  Receiver cz *Cz4Image
//  Param header CzHeader
//  Param data []byte
//
func (cz *Cz4Image) Load(header CzHeader, data []byte) {
	cz.CzHeader = header
	cz.Raw = data
	err := restruct.Unpack(cz.Raw[15:], binary.LittleEndian, &cz.Cz3Header)
	if err != nil {
		panic(err)
	}
	glog.V(6).Infoln("cz4 header ", cz.Cz3Header)
	cz.OutputInfo = GetOutputInfo(cz.Raw[int(cz.HeaderLength):])
}

// decompress
//  Description Decompress and decode CZ4 data
//  Receiver cz *Cz4Image
//
func (cz *Cz4Image) decompress() {
	buf := Decompress(cz.Raw[int(cz.HeaderLength)+cz.OutputInfo.Offset:], cz.OutputInfo)
	glog.V(6).Infoln("uncompress size ", len(buf))
	cz.Image = LineDiff4(&cz.CzHeader, buf)
}

// GetImage
//  Description Get decoded image
//  Receiver cz *Cz4Image
//  Return image.Image
//
func (cz *Cz4Image) GetImage() image.Image {
	if cz.Image == nil {
		cz.decompress()
	}
	return cz.Image
}

// Export
//  Description Export image as PNG
//  Receiver cz *Cz4Image
//  Param w io.Writer
//  Return error
//
func (cz *Cz4Image) Export(w io.Writer) error {
	if cz.Image == nil {
		cz.decompress()
	}

	var nrgbaImg *image.NRGBA
	if img, ok := cz.Image.(*image.NRGBA); ok {
		nrgbaImg = img
	} else {
		glog.V(0).Infof("Export CZ4: Converting image to NRGBA (was %T)\n", cz.Image)
		nrgbaImg = ImageToNRGBA(cz.Image)
	}

	glog.V(0).Infof("Export CZ4: %dx%d, Colorbits=%d\n",
		nrgbaImg.Rect.Dx(), nrgbaImg.Rect.Dy(), cz.CzHeader.Colorbits)

	return png.Encode(w, nrgbaImg)
}

// Import
//  Description Import PNG and encode as CZ4
//  Receiver cz *Cz4Image
//  Param r io.Reader
//  Param fillSize bool
//  Return error
//
func (cz *Cz4Image) Import(r io.Reader, fillSize bool) error {
	var err error
	cz.PngImage, err = png.Decode(r)
	if err != nil {
		return err
	}

	glog.V(0).Infof("Import CZ4: PNG source type=%T, bounds=%v\n",
		cz.PngImage, cz.PngImage.Bounds())

	// Force NRGBA conversion
	pic, ok := cz.PngImage.(*image.NRGBA)
	if !ok {
		glog.V(0).Infof("Import CZ4: Converting to NRGBA (source was %T)\n", cz.PngImage)
		pic = ImageToNRGBA(cz.PngImage)
	}

	// Verify pixel buffer size
	expectedPixLen := pic.Rect.Dx() * pic.Rect.Dy() * 4
	if len(pic.Pix) != expectedPixLen {
		glog.Errorf("Import CZ4: Pix length mismatch! expected=%d, got=%d\n",
			expectedPixLen, len(pic.Pix))
	}

	if cz.CzHeader.Colorbits != 32 {
		glog.Warningf("Import CZ4: Colorbits=%d (expected 32)\n", cz.CzHeader.Colorbits)
	}

	data := DiffLine4(cz.CzHeader, pic)

	// Verify output size: w*h*3 (RGB) + w*h (Alpha) = w*h*4
	expectedDataLen := int(cz.Width) * int(cz.Heigth) * 4
	if len(data) != expectedDataLen {
		glog.Errorf("Import CZ4: DiffLine4 output size mismatch! expected=%d, got=%d\n",
			expectedDataLen, len(data))
	} else {
		glog.V(0).Infof("Import CZ4: DiffLine4 OK, generated %d bytes\n", len(data))
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

	glog.V(0).Infof("Import CZ4: Compressed to %d bytes (RawSize=%d)\n",
		cz.OutputInfo.TotalCompressedSize, cz.OutputInfo.TotalRawSize)

	return nil
}

// Write
//  Description Write CZ4 image to output
//  Receiver cz *Cz4Image
//  Param w io.Writer
//  Return error
//
func (cz *Cz4Image) Write(w io.Writer) error {
	var err error
	// Force CZ4 magic
	cz.CzHeader.Magic = []byte{'C', 'Z', '4', 0}

	if cz.CzHeader.Colorbits != 32 {
		glog.Warningf("Write CZ4: Forcing Colorbits from %d to 32\n", cz.CzHeader.Colorbits)
		cz.CzHeader.Colorbits = 32
	}

	glog.V(0).Infof("Write CZ4: %dx%d, %d blocks, RawSize=%d\n",
		cz.Width, cz.Heigth, cz.OutputInfo.FileCount, cz.OutputInfo.TotalRawSize)

	err = WriteStruct(w, &cz.CzHeader, &cz.Cz3Header, cz.OutputInfo)
	if err != nil {
		return err
	}
	_, err = w.Write(cz.Raw)
	return err
}
