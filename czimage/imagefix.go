package czimage

import (
	"github.com/golang/glog"
	"image"
	"image/color"
)

func PanelImage(header *CzHeader, colorPanel [][]byte, data []byte) image.Image {
	width := int(header.Width)
	height := int(header.Heigth)
	pic := image.NewNRGBA(image.Rect(0, 0, width, height))
	// B,G,R,A
	// 0,1,2,3
	i := 0
	for y := 0; y < int(header.Heigth); y++ {
		for x := 0; x < int(header.Width); x++ {
			pic.SetNRGBA(x, y, color.NRGBA{
				R: colorPanel[data[i]][2],
				G: colorPanel[data[i]][1],
				B: colorPanel[data[i]][0],
				A: colorPanel[data[i]][3],
			})
			i++
		}
	}
	return pic
}

// DiffLine 图像拆分
//  Description 图像拆分，cz3用 png->data
//  Param header CzHeader
//  Param img image.Image
//  Return data
//
func DiffLine(header CzHeader, pic *image.NRGBA) (data []byte) {
	width := int(header.Width)
	height := int(header.Heigth)

	if width != pic.Rect.Size().X || height != pic.Rect.Size().Y {
		glog.V(2).Infof("图片大小不匹配，应该为 w%d h%d\n", width, height)
		return nil
	}
	data = make([]byte, len(pic.Pix))
	
	// PATCH YOREMI: blockHeight calculation (GARbro algorithm)
	blockHeight := (height + 2) / 3
	
	glog.V(0).Infof("DiffLine: height=%d, colorblock=%d, blockHeight=%d\n", height, header.Colorblock, blockHeight)
	
	pixelByteCount := int(header.Colorbits >> 3)
	lineByteCount := width * pixelByteCount
	
	preLine := make([]byte, lineByteCount)
	currLine := make([]byte, lineByteCount)  // Buffer pour la ligne courante
	
	i := 0
	for y := 0; y < height; y++ {
		// PATCH YOREMI: Copier pic.Pix dans currLine au lieu de créer un slice alias
		// Bug original: currLine = pic.Pix[i:...] créait un alias qui modifiait pic.Pix
		// Solution: Copier dans un buffer séparé
		copy(currLine, pic.Pix[i:i+lineByteCount])
		
		// Algorithme EXACT du code original
		if y%blockHeight != 0 {
			for x := 0; x < lineByteCount; x++ {
				currLine[x] -= preLine[x]
				// 因为是每一行较上一行的变化，故逆向执行时需要累加差异
				preLine[x] += currLine[x]
			}
		} else {
			copy(preLine, currLine)
		}

		copy(data[i:i+lineByteCount], currLine)
		i += lineByteCount
	}
	
	return data
}

// LineDiff 拆分图像还原
//  Description 拆分图像还原，cz3用 data->png
//  Param header *CzHeader
//  Param data []byte
//  Return image.Image
//
func LineDiff(header *CzHeader, data []byte) image.Image {
	//os.WriteFile("../data/LB_EN/IMAGE/ld.data", data, 0666)
	width := int(header.Width)
	height := int(header.Heigth)
	pic := image.NewNRGBA(image.Rect(0, 0, width, height))
	
	// PATCH YOREMI: blockHeight calculation (must match DiffLine)
	blockHeight := (height + 2) / 3
	
	glog.V(0).Infof("LineDiff: height=%d, colorblock=%d, blockHeight=%d\n", height, header.Colorblock, blockHeight)
	
	pixelByteCount := int(header.Colorbits >> 3)
	lineByteCount := width * pixelByteCount
	var currLine []byte
	preLine := make([]byte, lineByteCount)
	i := 0
	for y := 0; y < height; y++ {
		currLine = data[i : i+lineByteCount]
		if y%blockHeight != 0 {
			for x := 0; x < lineByteCount; x++ {
				currLine[x] += preLine[x]
			}
		}
		// PATCH YOREMI: Copier les données au lieu d'aliaser
		// Bug original: preLine = currLine créait un alias, pas une copie
		copy(preLine, currLine)
		
		if pixelByteCount == 4 {
			// y*pic.Stride : (y+1)*pic.Stride
			copy(pic.Pix[i:i+lineByteCount], currLine)
		} else if pixelByteCount == 3 {
			for x := 0; x < lineByteCount; x += 3 {
				pic.SetNRGBA(x/3, y, color.NRGBA{R: currLine[x], G: currLine[x+1], B: currLine[x+2], A: 0xFF})
			}
		}
		i += lineByteCount
	}
	
	//os.WriteFile("../data/LB_EN/IMAGE/ld.data.pix", pic.Pix, 0666)
	return pic
}

// DiffLine4 CZ4 encode: NRGBA pixels → separated [RGB w*h*3][Alpha w*h] with delta encoding
//  Description CZ4 stores RGB and Alpha in separate sections, each with independent
//  delta line encoding. This differs from CZ3 which stores interleaved RGBA.
//  Param header CzHeader
//  Param pic *image.NRGBA
//  Return data []byte (length = w*h*3 + w*h = w*h*4)
//
func DiffLine4(header CzHeader, pic *image.NRGBA) []byte {
	width := int(header.Width)
	height := int(header.Heigth)

	if width != pic.Rect.Size().X || height != pic.Rect.Size().Y {
		glog.V(2).Infof("DiffLine4: image size mismatch, expected w%d h%d\n", width, height)
		return nil
	}

	blockHeight := (height + 2) / 3

	rgbSize := width * height * 3
	alphaSize := width * height
	data := make([]byte, rgbSize+alphaSize)

	prevRGB := make([]byte, width*3)
	currRGB := make([]byte, width*3)
	prevAlpha := make([]byte, width)
	currAlpha := make([]byte, width)

	rgbOffset := 0
	alphaOffset := rgbSize

	for y := 0; y < height; y++ {
		// Extract RGB and Alpha from NRGBA pixels (R,G,B,A interleaved)
		pixOffset := y * pic.Stride
		for x := 0; x < width; x++ {
			currRGB[x*3] = pic.Pix[pixOffset+x*4]
			currRGB[x*3+1] = pic.Pix[pixOffset+x*4+1]
			currRGB[x*3+2] = pic.Pix[pixOffset+x*4+2]
			currAlpha[x] = pic.Pix[pixOffset+x*4+3]
		}

		if y%blockHeight != 0 {
			// Delta encode RGB
			for x := 0; x < width*3; x++ {
				currRGB[x] -= prevRGB[x]
				prevRGB[x] += currRGB[x]
			}
			// Delta encode Alpha
			for x := 0; x < width; x++ {
				currAlpha[x] -= prevAlpha[x]
				prevAlpha[x] += currAlpha[x]
			}
		} else {
			copy(prevRGB, currRGB)
			copy(prevAlpha, currAlpha)
		}

		copy(data[rgbOffset:rgbOffset+width*3], currRGB)
		copy(data[alphaOffset:alphaOffset+width], currAlpha)
		rgbOffset += width * 3
		alphaOffset += width
	}

	glog.V(0).Infof("DiffLine4: %dx%d, blockHeight=%d, output=%d bytes (RGB=%d + Alpha=%d)\n",
		width, height, blockHeight, len(data), rgbSize, alphaSize)

	return data
}

// LineDiff4 CZ4 decode: separated [RGB w*h*3][Alpha w*h] with delta → NRGBA image
//  Description Reconstructs NRGBA image from CZ4's separated channel layout.
//  Param header *CzHeader
//  Param data []byte
//  Return image.Image
//
func LineDiff4(header *CzHeader, data []byte) image.Image {
	width := int(header.Width)
	height := int(header.Heigth)
	pic := image.NewNRGBA(image.Rect(0, 0, width, height))

	blockHeight := (height + 2) / 3

	rgbSize := width * height * 3

	prevRGB := make([]byte, width*3)
	prevAlpha := make([]byte, width)

	rgbOffset := 0
	alphaOffset := rgbSize

	for y := 0; y < height; y++ {
		// Read current RGB and Alpha lines
		currRGB := make([]byte, width*3)
		copy(currRGB, data[rgbOffset:rgbOffset+width*3])

		currAlpha := make([]byte, width)
		copy(currAlpha, data[alphaOffset:alphaOffset+width])

		if y%blockHeight != 0 {
			// Delta decode RGB
			for x := 0; x < width*3; x++ {
				currRGB[x] += prevRGB[x]
			}
			// Delta decode Alpha
			for x := 0; x < width; x++ {
				currAlpha[x] += prevAlpha[x]
			}
		}

		copy(prevRGB, currRGB)
		copy(prevAlpha, currAlpha)

		// Interleave to NRGBA pixels
		pixOffset := y * pic.Stride
		for x := 0; x < width; x++ {
			pic.Pix[pixOffset+x*4] = currRGB[x*3]
			pic.Pix[pixOffset+x*4+1] = currRGB[x*3+1]
			pic.Pix[pixOffset+x*4+2] = currRGB[x*3+2]
			pic.Pix[pixOffset+x*4+3] = currAlpha[x]
		}

		rgbOffset += width * 3
		alphaOffset += width
	}

	glog.V(0).Infof("LineDiff4: %dx%d, blockHeight=%d, decoded from %d bytes\n",
		width, height, blockHeight, len(data))

	return pic
}
