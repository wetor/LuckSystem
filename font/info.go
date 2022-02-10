package font

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"github.com/golang/glog"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"io"
	"os"
)

type DrawSize struct {
	X uint8 // x offset
	W uint8 // width
	Y uint8 // y offset
}
type CharSize struct {
	X uint8
	W uint8
}

type Info struct {
	FontSize     uint16     // 字体大小
	BlockSize    uint16     // 字体所在区域大小（超过会被切割）
	CharNum      uint16     // 字符数量
	CharNum2     uint16     `struct:"if=CharNum==100"`
	DrawSize     []DrawSize `struct:"size=CharNum==100?CharNum2:CharNum"`
	UnicodeIndex []uint16   `struct:"size=65536"` // unicode -> imgindex
	UnicodeSize  []CharSize `struct:"size=65536"`
	// FontMap     map[rune]uint16 `struct:"-"` // unicode -> imgindex
	FontFace     font.Face `struct:"-"`
	IndexUnicode []rune    `struct:"-"` // imgindex -> unicode
}

func LoadFontInfo(data []byte) *Info {
	info := new(Info)
	err := restruct.Unpack(data, binary.LittleEndian, info)
	if err != nil {
		glog.Fatalln("restruct.Unpack", err)
	}
	if info.CharNum == 100 {
		info.CharNum = info.CharNum2
		info.CharNum2 = 100
	}
	// info.FontMap = make(map[rune]uint16)
	info.IndexUnicode = make([]rune, info.CharNum)
	// 6 + 3*7112
	// fmt.Println(info.FontSize, info.CharSize, info.CharNum)
	for unicode, index := range info.UnicodeIndex {
		if index != 0 || unicode == 32 {
			// info.FontMap[rune(unicode)] = index
			info.IndexUnicode[int(index)] = rune(unicode)
		}
	}
	//glog.V(6).Infoln("font info", info.FontSize, info.CharNum, info.BlockSize, len(info.FontMap))
	//
	//for i, v := range info.UnicodeSize {
	//	if v.W != 0 {
	//		glog.V(6).Infof("%d %v UnicodeSize:%v", i, string(rune(i)), v)
	//	}
	//}

	return info
	// for unicode, index := range info.FontMap {

	// 	uni := make([]byte, 2)
	// 	binary.LittleEndian.PutUint16(uni, uint16(unicode))
	// 	str, _ := charset.ToUTF8(charset.Unicode, uni)

	// 	fmt.Println(index, unicode, str, info.DrawSize[index], info.UnicodeSize[unicode])

	// }
}

func (i *Info) Get(unicode rune) (int, DrawSize, CharSize) {
	index := i.UnicodeIndex[unicode]
	if unicode != 32 && index == 0 {
		panic("不存在此字符 " + string(unicode))
	}
	return int(index), i.DrawSize[index], i.UnicodeSize[unicode]
}

// CreateFontInfo 创建字体Info信息
//  Description
//  Param fontSize int 字体实际大小
//  Param blockSize int 字体所在区域大小（超过会被切割）
//  Return *FontInfo
//
func CreateFontInfo(fontSize, blockSize int) *Info {

	info := &Info{
		FontSize:     uint16(fontSize),
		BlockSize:    uint16(blockSize),
		UnicodeIndex: make([]uint16, 65536),
		UnicodeSize:  make([]CharSize, 65536),
	}
	// info.FontMap = make(map[rune]uint16)
	// info.IndexUnicode = make(map[int]rune)
	return info
}

// SetChars
//  Description 如果startIndex=0且allChar为空，则为仅重绘
//  Receiver i *Info
//  Param fontFile string 字体文件
//  Param allChar string 全字符串，若第一个字符不是空格，会自动补充为空格
//  Param startIndex int 开始位置，即字库上面跳过多少字符
//  Param reDraw bool 是否用新字体重绘startIndex之前的字符
//
func (i *Info) SetChars(fontFile, allChar string, startIndex int, reDraw bool) {

	glog.V(6).Infof("SetChars font:%v addCharNum:%v starIndex:%v\n", fontFile, len(allChar), startIndex)
	// 加载字体
	data, err := os.ReadFile(fontFile)
	if err != nil {
		glog.Fatalln(err)
	}
	font, err := opentype.Parse(data)
	if err != nil {
		glog.Fatalln(err)
	}
	i.FontFace, err = opentype.NewFace(font, &opentype.FaceOptions{
		Size: float64(i.FontSize),
		DPI:  72,
	})
	if err != nil {
		glog.Fatalln(err)
	}

	// 处理字符
	// 去重、排序
	chars := []rune(allChar)

	noReDraw := false
	if len(chars) == 0 && startIndex == 0 {
		if i.CharNum == 0 {
			glog.Fatalln("需要载入字体")
		}
		if reDraw {
			chars = make([]rune, 0, i.CharNum)
			for _, char := range i.IndexUnicode {
				if char == 0 {
					chars = append(chars, '□')
				} else {
					chars = append(chars, char)
				}
			}
		} else {
			noReDraw = true
		}
	} else {
		for startIndex > int(i.CharNum) {
			i.DrawSize = append(i.DrawSize, DrawSize{})
			i.IndexUnicode = append(i.IndexUnicode, rune(0))
			i.CharNum++
		}

		tempDrawSize := make([]DrawSize, len(chars))
		i.DrawSize = append(i.DrawSize[:startIndex], append(tempDrawSize)...)

		tempIndexUnicode := make([]rune, len(chars))
		i.IndexUnicode = append(i.IndexUnicode[:startIndex], append(tempIndexUnicode)...)

		i.CharNum = uint16(len(i.DrawSize))
	}
	if !noReDraw {
		for index := 0; index < int(i.CharNum); index++ {
			var char rune
			if index < startIndex {
				if reDraw {
					char = i.IndexUnicode[index]
				} else {
					continue
				}
			} else {
				char = chars[index-startIndex]
			}
			// i.FontMap[char] = uint16(index)
			i.UnicodeIndex[char] = uint16(index)
			i.IndexUnicode[index] = char

			bounds, advance, ok := i.FontFace.GlyphBounds(char)

			if !ok {
				glog.Fatalf("字体文件中不存在的字符 %v %v\n", string(char), index)
				panic("字体文件中不存在的字符")
			}

			// fmt.Println(string(char), " ", bounds.Min.X.Floor(), " ", bounds.Min.Y.Floor()+int(i.FontSize))
			w := uint8(advance.Ceil())
			if char == 32 || w == 0 {
				w = uint8(i.FontSize)
			}
			i.DrawSize[index].X = uint8(bounds.Min.X.Floor())
			i.DrawSize[index].W = w
			i.DrawSize[index].Y = uint8(bounds.Min.Y.Floor())
			i.UnicodeSize[char].W = w

		}
	}

}
func (i *Info) Import(r io.Reader, opt ...interface{}) error {

	return nil
}

// Export
//  Description
//  Receiver i *Info
//  Param w io.Writer
//  Param opt ...interface{}
//  Return error
//
func (i *Info) Export(w io.Writer, opt ...interface{}) error {

	var err error
	for _, char := range i.IndexUnicode {

		if char == 0 {
			_, err = w.Write([]byte(string('□')))
		} else {
			_, err = w.Write([]byte(string(char)))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Write
//  Description
//  Receiver i *Info
//  Param w io.Writer
//  Param opt ...interface{}
//  Return error
//
func (i *Info) Write(w io.Writer, opt ...interface{}) error {

	data, err := restruct.Pack(binary.LittleEndian, i)
	if err != nil {
		glog.Fatalln("restruct.Pack", err)
	}
	_, err = w.Write(data)

	return err

}
