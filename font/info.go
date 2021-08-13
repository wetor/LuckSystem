package font

import (
	"encoding/binary"
	"fmt"
	"lucksystem/utils"

	"github.com/go-restruct/restruct"
)

type DrawSize struct {
	X uint8
	W uint8
	Y uint8
}
type CharSize struct {
	X uint8
	W uint8
}

type FontInfo struct {
	FontSzie    uint16
	CharSize    uint16
	CharNum     uint16
	DrawSize    []DrawSize      `struct:"size=CharNum"`
	UnicodeChar []uint16        `struct:"size=65536"`
	UnicodeSize []CharSize      `struct:"size=65536"`
	FontMap     map[rune]uint16 // unicode -> imgindex
}

func LoadFontInfo(data []byte) *FontInfo {
	info := new(FontInfo)
	err := restruct.Unpack(data, binary.LittleEndian, info)
	if err != nil {
		utils.Log("restruct.Unpack", err.Error())
		panic(err)
	}
	info.FontMap = make(map[rune]uint16)
	// 6 + 3*7112
	// fmt.Println(info.FontSzie, info.CharSize, info.CharNum)
	for i, ch := range info.UnicodeChar {
		if ch != 0 || i == 32 {
			info.FontMap[rune(i)] = ch
		}
	}
	fmt.Println("font info", info.CharNum, len(info.FontMap))
	return info
	// for unicode, index := range info.FontMap {

	// 	uni := make([]byte, 2)
	// 	binary.LittleEndian.PutUint16(uni, uint16(unicode))
	// 	str, _ := charset.ToUTF8(charset.Unicode, uni)

	// 	fmt.Println(index, unicode, str, info.DrawSize[index], info.UnicodeSize[unicode])

	// }
}

func (i *FontInfo) Get(unicode rune) (int, DrawSize) {
	index, has := i.FontMap[unicode]
	if !has {
		panic("不存在此字符 " + string(unicode))
	}
	return int(index), i.DrawSize[index]
}
