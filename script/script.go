package script

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"lucksystem/charset"
	"lucksystem/game/enum"
	"lucksystem/pak"

	"github.com/go-restruct/restruct"
)

// Script 从文件中直接读取到的代码结构∂
// 可用作运行时，不可直接导出，需要先转化为Entry
type Script struct {
	Info  `struct:"-"`
	Entry `struct:"-"`
	Codes []*CodeLine `struct:"while=!_eof"`
}

// CodeLine
// 导出逻辑
// 1.读入脚本文件，解析为[]CodeLine
// 2.vm运行，解析ParamBytes得到完整参数列表
// 3.将制定导出参数加入到Params，转为字符串导出
// 导入逻辑
// 1.读入脚本文件，解析为[]CodeLine
// 2.解析文本，添加进对应的Params中
// 3.vm运行，将Params中的参数替换到完整参数列表中，并转为ParamBytes
type CodeLine struct {
	CodeInfo   `struct:"-"`  // CodeLine信息
	FixedParam []uint16      `struct:"-"` // RawBytes数据，根据FixedFlag来解析
	ParamBytes []byte        `struct:"-"` // RawBytes数据，需要使用VM运行来解析
	Params     []interface{} `struct:"-"` // Export: ParamBytes解析而来; Import: 解析文本得来
	OpStr      string        `struct:"-"` // Opcode解析而来
	Len        uint16        // 代码长度
	Opcode     uint8         // 指令码
	FixedFlag  uint8         // 固定参数标记
	RawBytes   []byte        `struct:"size=Len - 4"` //参数数据 `struct:"size=((Len+ 1)& ^1)- 4"`
	Align      []byte        `struct:"size=Len & 1"` // 对齐，只读
}

func LoadScript(opts *LoadOptions) *Script {
	script := &Script{
		Info: Info{
			Name: opts.Name,
		},
	}
	script.InitEntry()
	if opts.Entry != nil {
		script.Name = opts.Entry.Name
		err := script.ReadData(opts.Entry.Data)
		if err != nil {
			panic(err)
		}
	} else {
		err := script.Read(opts.Filename)
		if err != nil {
			panic(err)
		}
	}

	return script
}

func (s *Script) ReadByEntry(entry *pak.Entry) error {
	s.Name = entry.Name
	return s.ReadData(entry.Data)
}

func (s *Script) Read(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		glog.V(8).Infoln("os.ReadFile", err)
		return err
	}
	return s.ReadData(data)
}

// SetOperateParams 设置需要导入导出的变量数据
//
//	1.index 当前代码行号，即Codes的下标
//	2.mode VMRun/VMRunExport/VMRunImport 模拟器/导出/导入 模式
//	  模拟器模式：
//	    直接返回，不作操作
//	  导出模式：
//	    将脚本中解析需要导出的参数列表paramList对Params赋值，对OpStr赋值
//	  导入模式：
//	    1.将从文件直接解析来的Params（[]string）根据从脚本中解析需要导出的paramList数据类型，转换Param类型
//	    2.将正确类型的Param(导出参数列表)与不导出参数列表合并成allParamList
//	    3.调用CodeParamsToBytes，将完整参数列表转为RawBytes
//	3.params 参数/设置列表，如长度为len
//	  参数列表：
//	    支持的特殊类型为：[]uint16（会解析为数个uint16参数）、*JumpParam（暂不支持ScriptName参数）、*StringParam
//	  设置列表：从后向前解析，如果为个别设置项空，其他设置项顺位后移，如opcode为空，则coding为params[len-1]
//	    1.params[len-1] opcode string 可选
//	    2.params[len-2] coding Charset 可选，所有字符串参数的默认编码
//	    3.params[len-3] export []bool 可选，导出列表，需要与参数列表(不含设置)数量相同，若少于有效参数数量，则默认补充false，不导出
func (s *Script) SetOperateParams(index int, mode enum.VMRunMode, params ...interface{}) error {
	if mode == enum.VMRun {
		return nil
	}

	paramNum := len(params)
	var paramsExport []bool // 导出参数列表
	strCharset := charset.UTF_8

	code := s.Codes[index]
	op := code.OpStr

	if paramNum > 0 && op == "UNDEFINED" {
		if val, ok := params[paramNum-1].(string); ok { // 最后一个参数为编码类型
			op = val // opcode
			paramNum--
		}
	}
	if paramNum > 0 {
		if val, ok := params[paramNum-1].(charset.Charset); ok { // 最后一个参数为编码类型
			strCharset = val // 编码
			paramNum--
		}
	}
	if paramNum > 0 {
		if val, ok := params[paramNum-1].([]bool); ok { // 最后一个参数为导出列表
			paramsExport = val // 导出参数列表
			paramNum--
		} else {
			// 默认全部导出
			paramsExport = make([]bool, paramNum)
			for i := 0; i < paramNum; i++ {
				paramsExport[i] = true
			}
		}
	}
	for len(paramsExport) < paramNum { // 参数数量大于导出列表，补全，默认false不导出
		paramsExport = append(paramsExport, false)
	}

	paramList := make([]interface{}, 0, paramNum)

	for i := 0; i < paramNum; i++ {
		if paramsExport[i] {
			switch param := params[i].(type) {
			case []uint16:
				for _, val := range param {
					paramList = append(paramList, val)
				}
			case *JumpParam:
				paramList = append(paramList, param)
				if param.GlobalIndex == 0 {
					s.AddExportGotoLabel(index, param.Position)
				}
			case *StringParam:
				paramList = append(paramList, param.Data)
			default:
				paramList = append(paramList, param)
			}
		}

	}

	if mode == enum.VMRunExport {
		code.Params = paramList
	} else if mode == enum.VMRunImport {
		// 导入模式
		if len(paramList) != len(code.Params) {
			panic("导入参数数量不匹配 " + strconv.Itoa(index))
		}
		// 导入数据类型转化
		for i := 0; i < len(paramList); i++ {
			switch paramList[i].(type) {
			case byte:
				val, _ := strconv.ParseUint(code.Params[i].(string)[2:], 16, 8)
				code.Params[i] = byte(val)
			case uint16:
				val, _ := strconv.ParseUint(code.Params[i].(string), 10, 16)
				code.Params[i] = uint16(val)
			case uint32:
				val, _ := strconv.ParseUint(code.Params[i].(string), 10, 32)
				code.Params[i] = uint32(val)
			}
			//fmt.Println(code.OpStr, code.Params[i], paramList[i])
		}

		// 将导入数据合并为列表
		allParamList := make([]interface{}, 0, paramNum)
		pi := 0
		for i := 0; i < paramNum; i++ {
			if paramsExport[i] { // 导出标记为true，则替换数据
				switch param := params[i].(type) {
				case []uint16:
					for j := range param {
						allParamList = append(allParamList, code.Params[pi])
						pi++
						_ = j
					}
				case *StringParam:
					param.Data = code.Params[pi].(string)
					allParamList = append(allParamList, param)
					pi++
				default:
					allParamList = append(allParamList, code.Params[pi])
					pi++
				}
			} else {
				switch param := params[i].(type) {
				case []uint16:
					for _, val := range param {
						allParamList = append(allParamList, val)
					}
				case charset.Charset:
					// 忽略
				case *StringParam:
					allParamList = append(allParamList, param)
				default:
					allParamList = append(allParamList, param)
				}
			}
		}
		// 将完整参数列表转为[]byte
		s.CodeParamsToBytes(code, strCharset, allParamList)
	}
	return nil
}

// Export 导出可编辑脚本
func (s *Script) Export(w io.Writer) error {
	glog.V(2).Infoln("Export: ", s.Name)
	var err error
	for i, code := range s.Codes {
		// 内部跳转
		labelIndex, has := s.ELabelMap[code.Pos]
		if has {
			code.LabelIndex = labelIndex
		}
		gotoIndex, has := s.EGotoMap[i]
		if has {
			code.GotoIndex = gotoIndex
		}
		// 跨文件跳转
		labelIndex, has = s.EGlobalLabelMap[code.Pos]
		if has {
			code.GlobalLabelIndex = labelIndex
		}
		gotoIndex, has = s.EGlobalGotoMap[i]
		if has {
			code.GlobalGotoIndex = gotoIndex
		}
	}
	bw := bufio.NewWriter(w)
	for _, code := range s.Codes {
		str := ToStringCodeParams(code)
		str = strings.Replace(str, "\n", "\\n", -1)
		_, err = fmt.Fprintln(bw, str)
		if err != nil {
			return err
		}
	}
	return bw.Flush()

}

// Import 导入可编辑脚本
func (s *Script) Import(r io.Reader) error {
	glog.V(2).Infoln("Import: ", s.Name)
	br := bufio.NewReader(r)
	for i, code := range s.Codes {
		line, err := br.ReadString('\n')
		if len(line) <= 1 {
			return errors.New("文本行不能为空 " + strconv.Itoa(i))
		} else if err == io.EOF {
			return errors.New("文本行数不匹配 " + strconv.Itoa(i))
		} else if err != nil {
			return err
		}
		line = strings.Replace(line, "\\n", "\n", -1)
		ParseCodeParams(code, line)

		glog.V(6).Info(i)
		glog.V(6).Infof("%v", line)
		if code.LabelIndex > 0 {
			glog.V(6).Infof(" label%d:", code.LabelIndex)
		}
		glog.V(6).Infof("%s %v", code.OpStr, code.Params)
		if code.GotoIndex > 0 {
			glog.V(6).Infof(" {goto label%d}", code.GotoIndex)
		}

	}
	return nil
}

func (s *Script) ReadData(data []byte) error {

	err := restruct.Unpack(data, binary.LittleEndian, s)
	if err != nil {
		return err
	}
	s.CodeNum = len(s.Codes)
	// s.FormatCodes = make([]string, s.CodeNum)
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
	return nil
}

// Write
func (s *Script) CodeParamsToBytes(code *CodeLine, coding charset.Charset, params []interface{}) {
	buf := &bytes.Buffer{}
	size := 0
	position := 0
	for _, param := range code.FixedParam {
		l, _ := SetParam(buf, param, coding, false)
		size += l
	}
	for _, param := range params {
		l, j := SetParam(buf, param, coding, false)
		size += l
		if j {
			position = size - 4
		}

	}
	if buf.Len() != len(code.RawBytes) {
		glog.V(4).Infof("%v %v\n\t%v %v\n\t%v %v\n", code.OpStr, params, buf.Len(), buf.Bytes(), len(code.RawBytes), code.RawBytes)
	}
	code.Len = uint16(size + 4)
	position += 4
	code.Align = make([]byte, code.Len&1)
	code.RawBytes = buf.Bytes()
	if code.LabelIndex > 0 {
		s.AddImportLabel(code.LabelIndex, s.CurPos)
	}
	if code.GotoIndex > 0 {
		s.AddImportGoto(s.CurPos+position, code.GotoIndex)
	}
	if code.GlobalLabelIndex > 0 {
		s.AddImportGlobalLabel(code.GlobalLabelIndex, s.CurPos)
	}
	if code.GlobalGotoIndex > 0 {
		s.AddImportGlobalGoto(s.CurPos+position, code.GlobalGotoIndex)
	}

	s.CurPos += int(code.Len + code.Len&1)

}

func (s *Script) Write(w io.Writer) error {
	data, err := restruct.Pack(binary.LittleEndian, s)
	if err != nil {
		return err
	}

	for gotoPos, labelIndex := range s.IGotoMap {
		jumpPos, has := s.ILabelMap[labelIndex]
		if !has {
			return errors.New("Goto-Label不匹配 " + strconv.Itoa(labelIndex))
		}
		pos := make([]byte, 4)
		binary.LittleEndian.PutUint32(pos, uint32(jumpPos))
		copy(data[gotoPos:gotoPos+4], pos)
		glog.V(6).Infoln(data[gotoPos : gotoPos+4])
	}

	for gotoPos, labelIndex := range s.IGlobalGotoMap {
		jumpPos, has := s.IGlobalLabelMap[labelIndex]
		if !has {
			return errors.New("Global Goto-Label不匹配 " + strconv.Itoa(labelIndex))
		}
		pos := make([]byte, 4)
		binary.LittleEndian.PutUint32(pos, uint32(jumpPos))
		copy(data[gotoPos:gotoPos+4], pos)
		glog.V(6).Infoln(data[gotoPos : gotoPos+4])
	}

	_, err = w.Write(data)
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
		writeSize := uint16(0)
		switch coding {
		case charset.UTF_8:
			writeSize = uint16(0x10000 - size)
		case charset.ShiftJIS:
			fallthrough
		case charset.Unicode:
			fallthrough
		default:
			writeSize = uint16(size / 2)
		}
		binary.Write(w, binary.LittleEndian, writeSize)
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
//
//	1.data[0] param 数据
//	2.data[1] coding 可空，默认Unicode
//	3.data[2] len 可空，是否为lstring类型
//	return size 字节长度
//	return position 跳转标记位置
func SetParam(buf *bytes.Buffer, param interface{}, coding charset.Charset, hasLen bool) (size int, jump bool) {
	if len(coding) == 0 {
		coding = charset.Unicode
	}
	switch value := param.(type) {
	case string:
		size = CodeString(buf, value, hasLen, coding)
	case *StringParam:
		size = CodeString(buf, value.Data, value.HasLen, value.Coding)
	case *JumpParam:
		// 填充跳转 pos
		//if value.ScriptName != "" {
		//	size += CodeString(buf, value.ScriptName, false, coding)
		//}
		if value.Position > 0 || value.GlobalIndex > 0 { // 现在为labelIndex
			jump = true
			binary.Write(buf, binary.LittleEndian, uint32(value.Position))
			size += 4
		}
	default:
		tmp := buf.Len()
		binary.Write(buf, binary.LittleEndian, param)
		size = buf.Len() - tmp
	}
	return size, jump
}
