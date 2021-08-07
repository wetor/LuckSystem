package script

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"lucascript/charset"
	"lucascript/game/enum"
	"lucascript/pak"
	"lucascript/utils"
	"os"
	"path"
	"strconv"

	"github.com/go-restruct/restruct"
)

type ScriptFileOptions struct {
	FileName string
	GameName string
	Version  uint8
}
type JumpParam struct {
	ScriptName string
	Position   int
}

type StringParam struct {
	Data   string
	Coding charset.Charset
	HasLen bool
}

// ScriptFile 从文件中直接读取到的代码结构
// 可用作运行时，不可直接导出，需要先转化为ScriptEntry
type ScriptFile struct {
	ScriptInfo  `struct:"-"`
	ScriptEntry `struct:"-"`
	Codes       []*CodeLine `struct:"while=true"`
}

type ScriptInfo struct {
	FileName string
	Name     string
	GameName string
	Version  uint8
	CodeNum  int
}

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

type CodeInfo struct {
	Index      int // 序号
	Pos        int // 文件偏移，和Index绑定
	LabelIndex int // 跳转目标标记，从1开始
	GotoIndex  int // 跳转标记，为0则不使用
}

// ScriptEntry 脚本导入导出实体
// 不可用作运行时，需先转话为ScriptFile
type ScriptEntry struct {
	// 递增序列，从1开始
	IndexNext int

	// 导出：
	// 1.执行脚本，遇到GOTO等跳转指令，将jumpPos作为key，递增序列Index作为value，存入ELabelMap，即一个地址只能存在一个LabelIndex
	// 2.同时，将当前指令的CodeIndex作为key，1中和jumpPos对应的LabelIndex作为value，存入EGotoMap，即标记次条语句包含跳转到LabelIndex的指令
	ELabelMap map[int]int // Pos(跳转地址) -> LabelIndex(标签序号) ，Pos 通过GOTO (pos) 生成，Index 为序列
	EGotoMap  map[int]int // CodeIndex(代码序号) -> LabelIndex(标签序号) ，CodeIndex为当前语句序列，此语句含有跳转指令，跳转到LabelIndex

	// 导入当前位置
	CurPos int
	// 导入：
	// 1.解析文本，同时开始序列化脚本，转为二进制数据并写入。
	// 2.遇到Label标签，将LabelIndex作为key，当前语句开始位置的文件偏移Pos作为value，存入ILabelMap，即标签对应的跳转地址
	// 3.遇到GOTO等跳转指令时，将要跳转到的LabelIndex作为key，[jumpPos参数所在的文件偏移]作为value存入IGotoMap，即暂时留空，后续再补充数据
	// 4.数据写入完成，遍历IGotoMap，根据ILabelMap的key，即LabelIndex，在ILabelMap中取得语句偏移Pos，写入[jumpPos参数所在的文件偏移]位置，填充数据。
	ILabelMap map[int]int // LabelIndex(标签序号) -> CodeStartPos(代码开头地址，跳转目标地址)
	IGotoMap  map[int]int // GotoParamPos(跳转参数地址) -> LabelIndex(标签序号)
}

func (e *ScriptEntry) InitEntry() {
	e.ELabelMap = make(map[int]int)
	e.EGotoMap = make(map[int]int)

	e.ILabelMap = make(map[int]int)
	e.IGotoMap = make(map[int]int)

	e.IndexNext = 1
	e.CurPos = 0
}

func (e *ScriptEntry) addExportGotoLabel(codeIndex, pos int) int {

	val, has := e.ELabelMap[pos]
	if has {
		e.EGotoMap[codeIndex] = val
		return val
	}
	e.ELabelMap[pos] = e.IndexNext
	e.EGotoMap[codeIndex] = e.IndexNext
	e.IndexNext++
	return e.ELabelMap[pos]
}

// labelIndex, Goto参数位置
func (e *ScriptEntry) addImportGoto(pos, labelIndex int) {
	e.IGotoMap[pos] = labelIndex
}

// labelIndex, 当前代码位置
func (e *ScriptEntry) addImportLabel(labelIndex, pos int) {
	e.ILabelMap[labelIndex] = pos
}

func NewScriptFile(opt ScriptFileOptions) *ScriptFile {
	script := new(ScriptFile)
	script.FileName = opt.FileName
	_, script.Name = path.Split(opt.FileName)
	script.GameName = opt.GameName
	script.Version = opt.Version
	script.InitEntry()
	return script
}

// SetOperateParams 设置需要导入导出的变量数据
//   1.index 当前代码行号，即Codes的下标
//   2.mode VMRun/VMRunExport/VMRunImport 模拟器/导出/导入 模式
//     模拟器模式：
//       直接返回，不作操作
//     导出模式：
//       将脚本中解析需要导出的参数列表paramList对Params赋值，对OpStr赋值
//     导入模式：
//       1.将从文件直接解析来的Params（[]string）根据从脚本中解析需要导出的paramList数据类型，转换Param类型
//       2.将正确类型的Param(导出参数列表)与不导出参数列表合并成allParamList
//       3.调用CodeParamsToBytes，将完整参数列表转为RawBytes
//   3.params 参数/设置列表，如长度为len
//     参数列表：
//       支持的特殊类型为：[]uint16（会解析为数个uint16参数）、*JumpParam（暂不支持ScriptName参数）、*StringParam
//     设置列表：从后向前解析，如果为个别设置项空，其他设置项顺位后移，如opcode为空，则coding为params[len-1]
//       1.params[len-1] opcode string 可选
//       2.params[len-2] coding Charset 可选，所有字符串参数的默认编码
//       3.params[len-3] export []bool 可选，导出列表，需要与参数列表(不含设置)数量相同，若少于有效参数数量，则默认补充false，不导出
func (s *ScriptFile) SetOperateParams(index int, mode enum.VMRunMode, params ...interface{}) error {
	if mode == enum.VMRun {
		return nil
	}

	paramNum := len(params)
	var paramsExport []bool // 导出参数列表
	strCharset := charset.UTF_8
	op := GetOperateName() // runtime 向上两层的函数名

	if paramNum > 0 && op == "UNDEFINE" {
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

	code := s.Codes[index]

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
				s.addExportGotoLabel(index, param.Position)
			case *StringParam:
				paramList = append(paramList, param.Data)
			default:
				paramList = append(paramList, param)
			}
		}

	}

	if mode == enum.VMRunExport {
		code.Params = paramList
		code.OpStr = op
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
				val, _ := strconv.ParseUint(code.Params[i].(string), 10, 16)
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
func (s *ScriptFile) Export(file string) error {
	for i, code := range s.Codes {
		labelIndex, has := s.ELabelMap[code.Pos]
		if has {
			code.LabelIndex = labelIndex
		}
		gotoIndex, has := s.EGotoMap[i]
		if has {
			code.GotoIndex = gotoIndex
		}
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for i, code := range s.Codes {
		str := ToStringCodeParams(code)
		fmt.Fprintln(w, str)

		fmt.Println(i, str)
	}
	return w.Flush()

}

// Import 导入可编辑脚本
func (s *ScriptFile) Import(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	r := bufio.NewReader(f)
	for i, code := range s.Codes {
		line, err := r.ReadString('\n')
		if len(line) <= 1 {
			return errors.New("文本行不能为空 " + strconv.Itoa(i))
		} else if err == io.EOF {
			return errors.New("文本行数不匹配 " + strconv.Itoa(i))
		} else if err != nil {
			return err
		}
		ParseCodeParams(code, line)

		fmt.Print(i)
		if code.LabelIndex > 0 {
			fmt.Printf(" label%d:", code.LabelIndex)
		}
		fmt.Printf(" %s %v", code.OpStr, code.Params)
		if code.GotoIndex > 0 {
			fmt.Printf(" {goto label%d}", code.GotoIndex)
		}
		fmt.Print("\n")

	}
	return nil
}
func OpenScriptFile(entry *pak.FileEntry) (*ScriptFile, error) {
	script := &ScriptFile{}
	script.Name = entry.Name
	err := script.ReadData(entry.Data)
	if err != nil {
		return nil, err
	}
	return script, nil
}

// Read
func (s *ScriptFile) ReadByEntry(entry *pak.FileEntry) error {
	s.Name = entry.Name
	return s.ReadData(entry.Data)
}
func (s *ScriptFile) Read() error {

	data, err := os.ReadFile(s.FileName)
	if err != nil {
		utils.Log("os.ReadFile", err.Error())
		return err
	}
	return s.ReadData(data)
}

func (s *ScriptFile) ReadData(data []byte) error {
	fmt.Println(len(data))
	err := restruct.Unpack(data, binary.LittleEndian, s)
	if err != nil {
		utils.Log("restruct.Unpack", err.Error())
		// return err
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
func (s *ScriptFile) CodeParamsToBytes(code *CodeLine, coding charset.Charset, params []interface{}) {
	buf := &bytes.Buffer{}
	size := 0
	for _, param := range code.FixedParam {
		size += SetParam(buf, param, coding)
	}
	for _, param := range params {
		size += SetParam(buf, param, coding)
	}
	fmt.Printf("%v %v\n\t%v %v\n\t%v %v\n", code.OpStr, params, buf.Len(), buf.Bytes(), len(code.RawBytes), code.RawBytes)
	code.Len = uint16(size + 4)
	code.Align = make([]byte, code.Len&1)
	code.RawBytes = buf.Bytes()
	if code.LabelIndex > 0 {
		s.addImportLabel(code.LabelIndex, s.CurPos)
	}
	if code.GotoIndex > 0 {
		s.addImportGoto(s.CurPos+int(code.Len)-4, code.GotoIndex) // -4 最后一个参数
	}
	s.CurPos += int(code.Len + code.Len&1)

}
func (s *ScriptFile) Write() error {

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
		fmt.Println(data[gotoPos : gotoPos+4])
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
