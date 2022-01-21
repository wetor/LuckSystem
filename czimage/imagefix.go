package czimage

import (
	"image"
	"image/color"
	"lucksystem/utils"
	"math"
)

func PanelImage(header *CzHeader, colorPanel [][]byte, data []byte) image.Image {
	width := int(header.Width)
	height := int(header.Heigth)
	pic := image.NewRGBA(image.Rect(0, 0, width, height))
	// B,G,R,A
	// 0,1,2,3
	i := 0
	for y := 0; y < int(header.Heigth); y++ {
		for x := 0; x < int(header.Width); x++ {
			pic.SetRGBA(x, y, color.RGBA{
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

// DiffLine png->data
func DiffLine(header *CzHeader, img image.Image) (data []byte) {
	width := int(header.Width)
	height := int(header.Heigth)

	pic := img.(*image.RGBA)
	if width != pic.Rect.Size().X || height != pic.Rect.Size().Y {
		utils.Logf("图片大小不匹配，应该为 w%d h%d", width, height)
		return nil
	}
	data = make([]byte, len(pic.Pix))
	blockHeight := int(uint16(math.Ceil(float64(height) / float64(header.Colorblock))))
	pixelByteCount := int(header.Colorbits >> 3)
	lineByteCount := width * pixelByteCount
	var currLine []byte
	preLine := make([]byte, lineByteCount)
	i := 0
	for y := 0; y < height; y++ {
		currLine = pic.Pix[i : i+lineByteCount]
		if y%blockHeight != 0 {
			for x := 0; x < lineByteCount; x++ {
				currLine[x] -= preLine[x]
				// 因为是每一行较上一行的变化，故逆向执行时需要累加差异
				preLine[x] += currLine[x]
			}
		} else {
			preLine = currLine
		}
		if pixelByteCount == 4 {
			// y*pic.Stride : (y+1)*pic.Stride
			copy(data[i:i+lineByteCount], currLine)
		}
		i += lineByteCount
	}
	return data
}

// LineDiff data->png
func LineDiff(header *CzHeader, data []byte) image.Image {
	//os.WriteFile("../data/LB_EN/IMAGE/ld.data", data, 0666)
	width := int(header.Width)
	height := int(header.Heigth)
	pic := image.NewRGBA(image.Rect(0, 0, width, height))
	blockHeight := int(uint16(math.Ceil(float64(height) / float64(header.Colorblock))))
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
		preLine = currLine
		if pixelByteCount == 4 {
			// y*pic.Stride : (y+1)*pic.Stride
			copy(pic.Pix[i:i+lineByteCount], currLine)
		}
		i += lineByteCount
	}
	//os.WriteFile("../data/LB_EN/IMAGE/ld.data.pix", pic.Pix, 0666)
	return pic
}
