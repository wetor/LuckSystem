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
	ExtendedHeader []byte        // Raw bytes of extended header (HeaderLength - 15)
	ColorPanel     []color.NRGBA // []BGRA
	CzData
}

func (cz *Cz1Image) Load(header CzHeader, data []byte) {
	cz.CzHeader = header
	cz.Raw = data

	// Normalize Colorbits: values > 32 (e.g. 248=0xF8) are proprietary
	// markers for 8-bit indexed palette mode (same behavior as lbee-utils)
	if cz.Colorbits > 32 {
		glog.V(4).Infof("CZ1: Colorbits=%d > 32, normalizing to 8 (indexed palette)\n", cz.Colorbits)
		cz.Colorbits = 8
	}

	// Save extended header bytes (between base header and palette/data)
	if cz.HeaderLength > 15 {
		cz.ExtendedHeader = make([]byte, cz.HeaderLength-15)
		copy(cz.ExtendedHeader, cz.Raw[15:cz.HeaderLength])
	}

	offset := int(cz.HeaderLength)
	if cz.Colorbits == 4 || cz.Colorbits == 8 {
		colorCount := 1 << cz.Colorbits
		cz.ColorPanel = make([]color.NRGBA, colorCount)
		for i := 0; i < colorCount; i++ {
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
		offset += (1 << cz.Colorbits) * 4
	}
	buf := Decompress(cz.Raw[offset+cz.OutputInfo.Offset:], cz.OutputInfo)
	glog.V(6).Infoln("uncompress size", len(buf))

	switch cz.Colorbits {
	case 4:
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
		// RGB → NRGBA with full alpha
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
		// RGBA direct copy (CZ1 32-bit stores pixels as RGBA)
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
//	Description Import PNG and compress back to CZ1 format
//	Receiver cz *Cz1Image
//	Param r io.Reader
//	Param fillSize bool
//	Return error
func (cz *Cz1Image) Import(r io.Reader, fillSize bool) error {
	var err error
	cz.PngImage, err = png.Decode(r)
	if err != nil {
		panic(err)
	}

	// Convert any PNG format to NRGBA
	pic := ImageToNRGBA(cz.PngImage)

	width := int(cz.Width)
	height := int(cz.Heigth)
	if fillSize {
		pic = FillImage(pic, width, height)
	}

	if width != pic.Rect.Size().X || height != pic.Rect.Size().Y {
		glog.V(2).Infof("Image size mismatch, expected w%d h%d\n", width, height)
		return err
	}

	var data []byte

	switch cz.Colorbits {
	case 4:
		// 4-bit indexed: find closest palette entry for each pixel
		data = make([]byte, (width*height+1)/2)
		i := 0
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				c := pic.NRGBAAt(x, y)
				idx := cz.findClosestPaletteEntry(c)
				if i%2 == 0 {
					data[i/2] = idx & 0x0F
				} else {
					data[i/2] |= (idx & 0x0F) << 4
				}
				i++
			}
		}
	case 8:
		// 8-bit indexed: find closest palette entry for each pixel
		data = make([]byte, width*height)
		i := 0
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				c := pic.NRGBAAt(x, y)
				data[i] = cz.findClosestPaletteEntry(c)
				i++
			}
		}
	case 24:
		// RGB: 3 bytes per pixel
		data = make([]byte, width*height*3)
		i := 0
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				c := pic.NRGBAAt(x, y)
				data[i+0] = c.R
				data[i+1] = c.G
				data[i+2] = c.B
				i += 3
			}
		}
	case 32:
		// RGBA direct copy: 4 bytes per pixel
		data = pic.Pix
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

// findClosestPaletteEntry finds the palette index that best matches the given color
func (cz *Cz1Image) findClosestPaletteEntry(c color.NRGBA) uint8 {
	bestIdx := uint8(0)
	bestDist := int(^uint(0) >> 1) // max int
	for i, pc := range cz.ColorPanel {
		dr := int(c.R) - int(pc.R)
		dg := int(c.G) - int(pc.G)
		db := int(c.B) - int(pc.B)
		da := int(c.A) - int(pc.A)
		dist := dr*dr + dg*dg + db*db + da*da
		if dist == 0 {
			return uint8(i)
		}
		if dist < bestDist {
			bestDist = dist
			bestIdx = uint8(i)
		}
	}
	return bestIdx
}

func (cz *Cz1Image) Write(w io.Writer) error {
	var err error
	// Fix magic byte (same pattern as CZ3/CZ4)
	cz.CzHeader.Magic = []byte{'C', 'Z', '1', 0}
	glog.V(6).Infoln(cz.CzHeader)

	// Write base header (15 bytes)
	err = WriteStruct(w, &cz.CzHeader)
	if err != nil {
		return err
	}

	// Write extended header if present (HeaderLength - 15 bytes)
	if len(cz.ExtendedHeader) > 0 {
		_, err = w.Write(cz.ExtendedHeader)
		if err != nil {
			return err
		}
	}

	// Write palette in BGRA order (file format) if present
	if len(cz.ColorPanel) > 0 {
		for _, c := range cz.ColorPanel {
			bgra := [4]byte{c.B, c.G, c.R, c.A}
			_, err = w.Write(bgra[:])
			if err != nil {
				return err
			}
		}
	}

	// Write block table and compressed data
	err = WriteStruct(w, cz.OutputInfo)
	if err != nil {
		return err
	}
	_, err = w.Write(cz.Raw)

	return err
}
