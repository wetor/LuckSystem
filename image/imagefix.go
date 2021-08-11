package czimage

import (
	"image"
	"math"
)

func LineDiff(header *CzHeader, data []byte) image.Image {
	pic := image.NewRGBA(image.Rect(0, 0, int(header.Width), int(header.Heigth)))
	blockHeight := int(uint16(math.Ceil(float64(header.Heigth) / float64(header.Colorblock))))
	pixelByteCount := int(header.Colorbits >> 3)
	lineByteCount := int(header.Width) * pixelByteCount
	var currLine []byte
	preLine := make([]byte, lineByteCount)
	i := 0
	for y := 0; y < int(header.Heigth); y++ {
		currLine = data[i : i+lineByteCount]
		if y%blockHeight != 0 {
			for x := 0; x < lineByteCount; x++ {
				currLine[x] += preLine[x]
			}
		}
		if pixelByteCount == 4 {
			// y*pic.Stride : (y+1)*pic.Stride
			copy(pic.Pix[i:i+lineByteCount], currLine)
		}

		// for x := 0; x < int(header.Width); x++ {
		// 	if pixelByteCount == 4 {
		// 		pic.SetRGBA(x, y, color.RGBA{
		// 			R: currLine[x*pixelByteCount+0],
		// 			G: currLine[x*pixelByteCount+1],
		// 			B: currLine[x*pixelByteCount+2],
		// 			A: currLine[x*pixelByteCount+3],
		// 		})
		// 	} else if pixelByteCount == 3 {
		// 		pic.SetRGBA(x, y, color.RGBA{
		// 			R: currLine[x*pixelByteCount+0],
		// 			G: currLine[x*pixelByteCount+1],
		// 			B: currLine[x*pixelByteCount+2],
		// 			A: 0xFF,
		// 		})
		// 	}
		// }
		preLine = currLine
		i += lineByteCount
	}
	return pic
}
