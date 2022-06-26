package font

import (
	"fmt"
	"github.com/golang/glog"
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
	pak := pak.LoadPak(
		"../data/LB_EN/FONT.PAK",
		charset.UTF_8,
	)
	fmt.Println("pak header", pak.Header)

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
	pic := image.NewNRGBA(image.Rect(0, 0, 800, 600))

	fontFace, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 24,
		DPI:  72,
	})
	if err != nil {
		glog.Fatalln(err)
	}

	//c.Text(100, 100, "测试汉字，#12ABjgkloa!.理樹@「…は？　こんな夜に？　どこで？」")

	xx1, xx2, _ := fontFace.GlyphBounds('б')
	fmt.Println(xx1.Max.X.Ceil(), xx1.Max.Y.Ceil(), xx2.Ceil())
	_, img, _, width, ok := fontFace.Glyph(fixed.Point26_6{X: 0, Y: 0}, 'б')
	draw.Draw(pic, pic.Bounds().Add(image.Pt(10, 10)), img, img.Bounds().Min, draw.Src)
	draw.Draw(pic, pic.Bounds().Add(image.Pt(34, 10)), img, img.Bounds().Min, draw.Src)
	draw.Draw(pic, pic.Bounds().Add(image.Pt(58, 10)), img, img.Bounds().Min, draw.Src)
	_, img, _, width, ok = fontFace.Glyph(fixed.Point26_6{X: 0, Y: 0}, '测')
	draw.Draw(pic, pic.Bounds().Add(image.Pt(10, 10)), img, img.Bounds().Min, draw.Src)
	fmt.Println(width.Ceil(), ok)
	f, _ := os.Create("../data/Other/Font/ARHei-400.ttf.png")
	defer f.Close()
	png.Encode(f, pic)

}

func TestCreateLucaFont(t *testing.T) {
	restruct.EnableExprBeta()
	font, _ := os.Open("../data/Other/Font/ARHei-400.ttf")
	defer font.Close()
	f := CreateLucaFont(24, font, " !@#$%.,abcdefgABCDEFG12345测试中文汉字")
	pngFile, _ := os.Create("../data/Other/Font/ARHei-400.ttf.png")
	f.Export(pngFile, "../data/Other/Font/ARHei-400.allChar.txt")

	czFile, _ := os.Create("../data/Other/Font/ARHei-400.ttf.cz")
	infoFile, _ := os.Create("../data/Other/Font/ARHei-400.ttf.info")
	f.Write(czFile, infoFile)
}

func TestEidtLucaFont(t *testing.T) {
	restruct.EnableExprBeta()
	pak := pak.LoadPak(
		"../data/LB_EN/FONT.PAK",
		charset.UTF_8,
	)

	fmt.Println("pak header", pak.Header)

	f := LoadLucaFontPak(pak, "モダン", 32)
	font, _ := os.Open("../data/Other/Font/ARHei-400.ttf")
	defer font.Close()
	f.ReplaceChars(font, " !@#$%.,abcdefgABCDEFG12345测试中文汉字", 7113, false)
	//f := CreateLucaFont("测试字体", 24, "../data/Other/Font/ARHei-400.ttf", " !@#$%.,abcdefgABCDEFG12345测试中文汉字")

	pngFile, _ := os.Create("../data/Other/Font/モダン32e.png")
	f.Export(pngFile, "../data/Other/Font/モダン32e.allChar.txt")

	czFile, _ := os.Create("../data/Other/Font/モダン32e.cz")
	infoFile, _ := os.Create("../data/Other/Font/モダン32e.info")
	f.Write(czFile, infoFile)
}
func TestSPFont(t *testing.T) {

	restruct.EnableExprBeta()

	pak := pak.LoadPak(
		"/Volumes/NTFS/WorkSpace/Github/SummerPockets/font/FONT.PAK",
		charset.UTF_8,
	)

	fmt.Println("pak header", pak.Header)

	for i, f := range pak.Files {
		if i < 160 {
			continue
		}
		fmt.Println(f.ID, f.Name, f.Offset, f.Length, f.Replace)
	}

	f := LoadLucaFontPak(pak, "モダン", 32)
	file, _ := os.Create("../data/Other/Font/モダン32.png")
	defer file.Close()
	f.Export(file, "")

}
func TestLucaFont_Export(t *testing.T) {
	restruct.EnableExprBeta()
	var err error
	savePath := "../data/LB_EN/FONT/"
	infoFile := "info32"
	czFile := "明朝32"
	txtName := "info32.txt"
	pngName := "明朝32.png"

	infoData, err := os.ReadFile(savePath + infoFile)
	if err != nil {
		panic(err)
	}
	czData, err := os.ReadFile(savePath + czFile)
	if err != nil {
		panic(err)
	}
	font := LoadLucaFont(infoData, czData)

	txtFile := savePath + txtName

	pngFile, _ := os.Create(savePath + pngName)
	defer pngFile.Close()

	err = font.Export(pngFile, txtFile)
	if err != nil {
		panic(err)
	}
}

func TestLucaFont_Import(t *testing.T) {
	restruct.EnableExprBeta()
	var err error
	savePath := "../data/LB_EN/FONT/"
	infoFile := "info32"
	czFile := "明朝32"
	ttfFile := "../data/Other/Font/ARHei-400.ttf"
	addChars := "../data/Other/Font/allchar.txt"

	infoData, err := os.ReadFile(savePath + infoFile)
	if err != nil {
		panic(err)
	}
	czData, err := os.ReadFile(savePath + czFile)
	if err != nil {
		panic(err)
	}
	font := LoadLucaFont(infoData, czData)
	ttf, _ := os.Open(ttfFile)
	defer ttf.Close()

	//===============
	err = font.Import(ttf, 0, true, "")
	if err != nil {
		panic(err)
	}
	cz1, _ := os.Create(savePath + czFile + "_onlyRedraw")
	defer cz1.Close()
	info1, _ := os.Create(savePath + infoFile + "_onlyRedraw")
	defer info1.Close()
	err = font.Write(cz1, info1)
	if err != nil {
		panic(err)
	}
	//================
	err = font.Import(ttf, -1, false, addChars)
	if err != nil {
		panic(err)
	}
	cz2, _ := os.Create(savePath + czFile + "_addChar")
	defer cz2.Close()
	info2, _ := os.Create(savePath + infoFile + "_addChar")
	defer info2.Close()
	err = font.Write(cz2, info2)
	if err != nil {
		panic(err)
	}
}
