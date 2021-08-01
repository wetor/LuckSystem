package script

import (
	"bytes"
	"encoding/binary"
	"io"
	"lucascript/charset"
	"os"

	"github.com/go-restruct/restruct"
)

func (s *ScriptFile) Write() error {

	pos := 0
	// 预处理 FixedParam
	for i, code := range s.Codes {
		code.Index = i
		code.Pos = pos
		pos += ((int(code.Len) + 1) & ^1) // 向上对齐2
		//(4 + len(code.ParamBytes))
		if code.FixedFlag > 0 {
			if code.FixedFlag >= 2 {
				code.FixedParam = make([]uint16, 2)
				code.FixedParam[0] = binary.LittleEndian.Uint16(code.RawBytes[0:2])
				code.FixedParam[1] = binary.LittleEndian.Uint16(code.RawBytes[2:4])
				code.ParamBytes = make([]byte, len(code.RawBytes)-4)
				copy(code.ParamBytes, code.RawBytes[4:])

			} else {
				code.FixedParam = make([]uint16, 1)
				code.FixedParam[0] = binary.LittleEndian.Uint16(code.RawBytes[0:2])
				code.ParamBytes = make([]byte, len(code.RawBytes)-2)
				copy(code.ParamBytes, code.RawBytes[2:])
			}
		} else {
			code.ParamBytes = make([]byte, len(code.RawBytes))
			copy(code.ParamBytes, code.RawBytes)
		}
	}

	data, err := restruct.Pack(binary.LittleEndian, s)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.FileName+".out", data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func CodeString(w io.Writer, data string, hasLen bool, coding charset.Charset) int {

	dst, err := charset.UTF8To(coding, []byte(data))
	if err != nil {
		panic(err)
	}
	buf := []byte(dst)
	size := len(buf)
	if hasLen {
		binary.Write(w, binary.LittleEndian, uint16(size/2))
		size += 2
	}
	w.Write(buf)

	switch coding {
	case charset.ShiftJIS:
		fallthrough
	case charset.UTF_8:
		w.Write([]byte{0x00})
		size += 1
	case charset.Unicode:
		fallthrough
	default:
		w.Write([]byte{0x00, 0x00})
		size += 2
	}

	return size
}

// SetParam 参数转为字节
//   1.data[0] param 数据
//   2.data[1] coding 可空，默认Unicode
//   3.data[2] len 可空，是否为lstring类型
//   return size 字节长度
func SetParam(buf *bytes.Buffer, data ...interface{}) int {

	size := 0
	lenStr := false
	var coding charset.Charset
	if len(data) >= 2 {
		coding = data[1].(charset.Charset)
	} else {
		coding = charset.Unicode
	}
	if len(data) >= 3 {
		lenStr = data[2].(bool)
	}
	switch value := data[0].(type) {
	case string:
		size = CodeString(buf, value, lenStr, coding)
	case *StringParam:
		size = CodeString(buf, value.Data, value.HasLen, value.Coding)
	case *JumpParam:
		// 填充跳转 pos
		if value.ScriptName != "" {
			size += CodeString(buf, value.ScriptName, false, coding)
		}
		if value.Position > 0 { // 现在为labelIndex
			binary.Write(buf, binary.LittleEndian, uint32(value.Position))
			size += 4
		}
	default:
		tmp := buf.Len()
		binary.Write(buf, binary.LittleEndian, data[0])
		size = buf.Len() - tmp
	}
	return size
}
