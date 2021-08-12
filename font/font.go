package font

import (
	"image"
	"image/draw"
	"lucascript/czimage"
	"lucascript/pak"
	"strconv"
)

type LucaFont struct {
	Name    string
	Size    int
	CzImage czimage.CzImage
	Info    *FontInfo
	Image   *image.RGBA
}

var FONT_PAK = "../data/LB_EN/FONT.PAK"

// モダン/明朝/丸ゴシック/ゴシック
// 12 14 16 18 20 24 28 30 32 36 72
func LoadLucaFont(pak *pak.PakFile, name string, size int) *LucaFont {
	font := new(LucaFont)

	infoFile, _ := pak.Get("info" + strconv.Itoa(size))
	font.Info = LoadFontInfo(infoFile.Data)

	imageFile, _ := pak.Get(name + strconv.Itoa(size))
	font.CzImage, _ = czimage.LoadCzImage(imageFile.Data)

	font.Image = font.CzImage.Get().(*image.RGBA)
	return font
}
func (f *LucaFont) GetCharImage(unicode rune) image.Image {

	index, _ := f.Info.Get(unicode)
	size := int(f.Info.CharSize)
	y := index / 100
	x := index % 100
	return f.Image.SubImage(image.Rect(x*size, y*size, (x+1)*size, (y+1)*size))
}
func (f *LucaFont) GetStringImageList(str string) []image.Image {
	imgs := make([]image.Image, 0, len(str))
	for _, r := range str {
		imgs = append(imgs, f.GetCharImage(r))
	}
	return imgs
}

func (f *LucaFont) GetStringImage(str string) image.Image {
	imgW := int(f.Info.CharSize)
	imgs := f.GetStringImageList(str)
	pic := image.NewRGBA(image.Rect(0, 0, len(imgs)*imgW, imgW))
	for i, img := range imgs {
		draw.Draw(pic, pic.Bounds().Add(image.Pt(i*imgW, 0)), img, img.Bounds().Min, draw.Src)
	}
	return pic
}
