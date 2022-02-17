package font

import (
	"bytes"
	"errors"
	"github.com/golang/glog"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"image/png"
	"io"
	"lucksystem/czimage"
	"lucksystem/pak"
	"math"
	"os"
	"strconv"
)

type LucaFont struct {
	Size    int
	CzImage czimage.CzImage
	Info    *Info
	Image   *image.RGBA
}

// LoadLucaFontPak 通过pak加载LucaFont
//  Description
//  Param pak *pak.PakFile
//  Param fontName string モダン/明朝/丸ゴシック/ゴシック
//  Param size int 12 14 16 18 20 24 28 30 32 36 72
//  Return *LucaFont
//
func LoadLucaFontPak(pak *pak.PakFile, fontName string, size int) *LucaFont {
	infoFile, err := pak.Get("info" + strconv.Itoa(size))
	if err != nil {
		glog.Fatalln(err)
	}
	imageFile, err := pak.Get(fontName + strconv.Itoa(size))
	if err != nil {
		glog.Fatalln(err)
	}
	return LoadLucaFont(infoFile.Data, imageFile.Data)
}

// LoadLucaFontFile 通过文件名加载LucaFont
//  Description
//  Param infoFilename string
//  Param imageFilename string
//  Return *LucaFont
//
func LoadLucaFontFile(infoFilename, imageFilename string) *LucaFont {
	infoFile, err := os.ReadFile(infoFilename)
	if err != nil {
		glog.Fatalln(err)
	}
	imageFile, err := os.ReadFile(imageFilename)
	if err != nil {
		glog.Fatalln(err)
	}
	return LoadLucaFont(infoFile, imageFile)
}

// LoadLucaFont 通过字节数据加载LucaFont
//  Description
//  Param infoFile []byte
//  Param imageFile []byte
//  Return *LucaFont
//
func LoadLucaFont(infoFile, imageFile []byte) *LucaFont {
	font := &LucaFont{}
	font.Info = LoadFontInfo(infoFile)
	font.Size = int(font.Info.FontSize)
	font.CzImage, _ = czimage.LoadCzImage(imageFile)

	font.Image = font.CzImage.GetImage().(*image.RGBA)
	return font
}

// GetCharImage 获取单个字符图像和偏移信息
//  Description
//  Receiver f *LucaFont
//  Param unicode rune
//  Return image.Image
//  Return DrawSize
//
func (f *LucaFont) GetCharImage(unicode rune) (image.Image, DrawSize) {

	index, draw, _ := f.Info.Get(unicode)
	size := int(f.Info.BlockSize)
	y := index / 100
	x := index % 100
	return f.Image.SubImage(image.Rect(x*size, y*size, (x+1)*size, (y+1)*size)), draw
}

// GetStringImageList 获取字符串每个字符的图像和偏移信息
//  Description
//  Receiver f *LucaFont
//  Param str string
//  Return []image.Image
//  Return []DrawSize
//
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

// GetStringImage 将字符串转化为图像
//  Description
//  Receiver f *LucaFont
//  Param str string
//  Return image.Image
//
func (f *LucaFont) GetStringImage(str string) image.Image {
	imgW := int(f.Info.BlockSize)
	imgs, draws := f.GetStringImageList(str)
	pic := image.NewRGBA(image.Rect(0, 0, len(imgs)*imgW, imgW*2))
	X := 0
	for i, img := range imgs {

		draw.Draw(pic, pic.Bounds().Add(image.Pt(X+int(draws[i].X), int(draws[i].Y))), img, img.Bounds().Min, draw.Src)
		X += int(draws[i].W)
	}
	_ = draws
	return pic
}

// CreateLucaFont 创建全新的字体
//  Description
//  Param fontSize int 字体大小
//  Param fontFile io.Reader 字体文件
//  Param allChar string 所有字符
//  Return *LucaFont
//
func CreateLucaFont(fontSize int, fontFile io.Reader, allChar string) *LucaFont {
	font := &LucaFont{
		Size: fontSize,
	}
	font.Info = CreateFontInfo(fontSize, fontSize+1)
	//font.Info.SetChars(, 20)
	font.ReplaceChars(fontFile, allChar, 0, true)

	return font
}

// ReplaceChars 替换字体中的字符
//  Description 替换字体中的字符信息以及图像, 如果startIndex=0且allChar为空，则为修改原字体
//  Receiver f *LucaFont
//  Param fontFile io.Reader 字体文件
//  Param allChar string 所替换的字符
//  Param startIndex int 开始序号（图像从上到下，从左到右计算）
//  Param reDraw bool 是否用新字体重绘startIndex之前的字符
//
func (f *LucaFont) ReplaceChars(fontFile io.Reader, allChar string, startIndex int, reDraw bool) {

	if f.Info == nil {
		glog.Fatalln("需要先载入或创建LucaFont")
		return
	}
	if len(allChar) == 0 && startIndex == 0 && !reDraw {
		// 什么都不做
		return
	}
	f.Info.SetChars(fontFile, allChar, startIndex, reDraw)
	size := int(f.Info.BlockSize)
	imageW := size*100 + 4                                         // 100个字符宽度+4
	imageH := size * int(math.Ceil(float64(f.Info.CharNum)/100.0)) // 对应行数高度
	oldImageH := size * int(math.Ceil(float64(startIndex)/100.0))

	pic := image.NewRGBA(image.Rect(0, 0, imageW, imageH))
	if !reDraw && f.Image != nil {
		img := f.Image.SubImage(image.Rect(0, 0, imageW, oldImageH))
		draw.Draw(pic, pic.Bounds().Add(image.Pt(0, 0)), img, img.Bounds().Min, draw.Src)
	}

	alphaMask := image.NewAlpha(image.Rect(0, 0, size, size))
	if reDraw {
		startIndex = 0
	}
	for i := startIndex; i < int(f.Info.CharNum); i++ {
		y := i / 100
		x := i % 100
		point := fixed.Point26_6{
			X: fixed.Int26_6(x * 64),
			Y: fixed.Int26_6(y * 64),
		}
		_, img, _, _, _ := f.Info.FontFace.Glyph(point, f.Info.IndexUnicode[i])
		// yOffset := dr.Min.Y + fontSize
		// fmt.Println(string(font.Info.IndexFont[i]), " ", dr.Min.Y+fontSize)
		if y == startIndex/100 {
			draw.Draw(pic, pic.Bounds().Add(image.Pt(x*size, y*size)), alphaMask, alphaMask.Bounds().Min, draw.Src)
		}
		draw.Draw(pic, pic.Bounds().Add(image.Pt(x*size, y*size)), img, img.Bounds().Min, draw.Src)
	}
	f.Image = pic
}

// Export
//  Description
//  Receiver f *LucaFont
//  Param w io.Writer
//  Param opt ...interface{}
//    opt[0] txtFilename io.Writer 导出的全字符文件
//  Return error
//
func (f *LucaFont) Export(w io.Writer, opt ...interface{}) error {
	err := png.Encode(w, f.Image)
	if err != nil {
		return err
	}
	if len(opt) > 0 {
		err = f.Info.Export(opt[0].(io.Writer))
	}
	return err
}

// Import
//  Description
//  Receiver f *LucaFont
//  Param r io.Reader 字体文件
//  Param opt ...interface{}
//    opt[0]	onlyRedraw	bool 仅使用新字体重绘，不增加新字符
//      or
//    opt[0]	allChar	string	增加的全字符，若startIndex==0，且第一个字符不是空格，会自动补充为空格
//    opt[1]	startIndex	int	开始位置，即字库上面跳过多少字符，-1为添加到最后
//    opt[2]	redraw	bool	是否用新字体重绘startIndex之前的字符
//  Return error
//
func (f *LucaFont) Import(r io.Reader, opt ...interface{}) error {
	if onlyRedraw, ok := opt[0].(bool); ok {
		if onlyRedraw {
			// 仅重绘
			f.ReplaceChars(r, "", 0, true)
			return nil
		} else {
			// 不做任何修改
			return nil
		}
	}
	if opt[1].(int) == -1 {
		opt[1] = int(f.Info.CharNum)
	}
	f.ReplaceChars(r, opt[0].(string), opt[1].(int), opt[2].(bool))
	return nil
}

// Write
//  Description
//  Receiver f *LucaFont
//  Param w io.Writer
//  Param opt ...interface{}
//    opt[0] infoFilename io.Writer 导出新的info文件
//  Return error
//
func (f *LucaFont) Write(w io.Writer, opt ...interface{}) error {
	var err error
	if f.CzImage != nil {
		// load

		img := bytes.NewBuffer(nil)
		err = png.Encode(img, f.Image)
		if err != nil {
			return err
		}
		czImg := bytes.NewBuffer(nil)
		err = f.CzImage.Import(img, true)
		if err != nil {
			return err
		}
		err = f.CzImage.Write(czImg)
		if err != nil {
			return err
		}
		_, err = w.Write(czImg.Bytes())
		if err != nil {
			return err
		}
	} else {
		// create
		return errors.New("LucaFont.Write 目前不支持创建的字体")
	}
	if len(opt) > 0 {
		err = f.Info.Write(opt[0].(io.Writer))
	}
	return err
}
