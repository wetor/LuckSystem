package font

import (
	"image"
	"image/draw"
	"lucksystem/czimage"
	"lucksystem/pak"
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
func (f *LucaFont) GetCharImage(unicode rune) (image.Image, DrawSize) {

	index, draw := f.Info.Get(unicode)
	size := int(f.Info.CharSize)
	y := index / 100
	x := index % 100
	return f.Image.SubImage(image.Rect(x*size, y*size, (x+1)*size, (y+1)*size)), draw
}
func (f *LucaFont) GetStringImageList(str string) ([]image.Image, []DrawSize) {
	imgs := make([]image.Image, 0, len(str))
	draws := make([]DrawSize, 0, len(str))
	for _, r := range str {
		img, draw := f.GetCharImage(r)
		imgs = append(imgs, img)
		draws = append(draws, draw)
	}
	return imgs, draws
}

func (f *LucaFont) GetStringImage(str string) image.Image {
	imgW := int(f.Info.CharSize)
	imgs, draws := f.GetStringImageList(str)
	pic := image.NewRGBA(image.Rect(0, 0, len(imgs)*imgW, imgW))
	X := 0
	for i, img := range imgs {

		draw.Draw(pic, pic.Bounds().Add(image.Pt(X+int(draws[i].X), int(draws[i].Y))), img, img.Bounds().Min, draw.Src)
		X += int(draws[i].W)
	}
	_ = draws
	return pic
}
