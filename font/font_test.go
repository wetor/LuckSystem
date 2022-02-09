package font

import (
	"fmt"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"image/png"
	"log"
	"lucksystem/charset"
	"lucksystem/pak"
	"os"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestFont(t *testing.T) {

	restruct.EnableExprBeta()
	pak := pak.NewPak(&pak.PakFileOptions{
		FileName: "../data/LB_EN/FONT.PAK",
		Coding:   charset.UTF_8,
	})
	err := pak.Open()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("pak header", pak.PakHeader)

	font := LoadLucaFontPak(pak, "モダン", 32)
	img := font.GetStringImage("ЁАБ #12ABjgkloa!.理樹@「…は？　こんな夜に？　どこで？」")
	f, _ := os.Create("../data/LB_EN/IMAGE/str/str.png")
	png.Encode(f, img)
	f.Close()
	// imgs := font.GetStringImageList("真人@「…戦いさ」")

	// for i, img := range imgs {
	// 	f, _ := os.Create("../data/LB_EN/IMAGE/str/" + strconv.Itoa(i) + ".png")
	// 	png.Encode(f, img)
	// 	f.Close()
	// }

}
func TestFreeTypeFont(t *testing.T) {
	data, err := os.ReadFile("../data/Other/Font/ARHei-400.ttf")
	if err != nil {
		log.Println(err)
		return
	}
	font, err := opentype.Parse(data)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(font.NumGlyphs())
	pic := image.NewRGBA(image.Rect(0, 0, 800, 600))
	c := NewContext(pic)

	c.SetFontFace(font, &opentype.FaceOptions{
		Size: 24,
		DPI:  72,
	})

	//c.Text(100, 100, "测试汉字，#12ABjgkloa!.理樹@「…は？　こんな夜に？　どこで？」")

	xx1, xx2, _ := c.fontDrawer.Face.GlyphBounds('б')
	fmt.Println(xx1.Max.X.Ceil(), xx1.Max.Y.Ceil(), xx2.Ceil())
	_, img, _, width, ok := c.fontDrawer.Face.Glyph(fixed.Point26_6{X: 0, Y: 0}, 'б')
	draw.Draw(pic, pic.Bounds().Add(image.Pt(10, 10)), img, img.Bounds().Min, draw.Src)
	draw.Draw(pic, pic.Bounds().Add(image.Pt(34, 10)), img, img.Bounds().Min, draw.Src)
	draw.Draw(pic, pic.Bounds().Add(image.Pt(58, 10)), img, img.Bounds().Min, draw.Src)
	_, img, _, width, ok = c.fontDrawer.Face.Glyph(fixed.Point26_6{X: 0, Y: 0}, '测')
	draw.Draw(pic, pic.Bounds().Add(image.Pt(10, 10)), img, img.Bounds().Min, draw.Src)
	fmt.Println(width.Ceil(), ok)
	f, _ := os.Create("../data/Other/Font/ARHei-400.ttf.png")
	defer f.Close()
	png.Encode(f, pic)

}

func TestCreateLucaFont(t *testing.T) {
	restruct.EnableExprBeta()
	f := CreateLucaFont(24, "../data/Other/Font/ARHei-400.ttf", " !@#$%.,abcdefgABCDEFG12345测试中文汉字")
	f.Export("../data/Other/Font/ARHei-400.ttf.png", true)
	f.Write("../data/Other/Font/ARHei-400.ttf.png", true)
}

func TestEidtLucaFont(t *testing.T) {
	restruct.EnableExprBeta()
	pak := pak.NewPak(&pak.PakFileOptions{
		FileName: "../data/LB_EN/FONT.PAK",
		Coding:   charset.UTF_8,
	})
	err := pak.Open()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("pak header", pak.PakHeader)

	f := LoadLucaFontPak(pak, "モダン", 32)
	f.ReplaceChars("../data/Other/Font/ARHei-400.ttf", " !@#$%.,abcdefgABCDEFG12345测试中文汉字", 7113, false)
	//f := CreateLucaFont("测试字体", 24, "../data/Other/Font/ARHei-400.ttf", " !@#$%.,abcdefgABCDEFG12345测试中文汉字")
	f.Export("../data/Other/Font/モダン32e.png", true)
	f.Write("../data/Other/Font/モダン32e", true)
}
func TestSPFont(t *testing.T) {

	restruct.EnableExprBeta()

	pak := pak.NewPak(&pak.PakFileOptions{
		FileName: "/Volumes/NTFS/WorkSpace/Github/SummerPockets/font/FONT.PAK",
		Coding:   charset.UTF_8,
	})
	err := pak.Open()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("pak header", pak.PakHeader)

	for i, f := range pak.Files {
		if i < 160 {
			continue
		}
		fmt.Println(f.ID, f.Name, f.Offset, f.Length, f.Replace)
	}

	f := LoadLucaFontPak(pak, "モダン", 32)
	f.Export("../data/SP/IMAGE/モダン32.png")

}
