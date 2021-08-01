package script

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"lucascript/charset"
	"lucascript/game/enum"
	"os"
	"strconv"
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

	// 导入：
	// 1.解析文本，同时开始序列化脚本，转为二进制数据并写入。
	// 2.遇到Label标签，将LabelIndex作为key，当前语句开始位置的文件偏移Pos作为value，存入ILabelMap，即标签对应的跳转地址
	// 3.遇到GOTO等跳转指令时，将要跳转到的LabelIndex作为key，[jumpPos参数所在的文件偏移]作为value存入IGotoMap，即暂时留空，后续再补充数据
	// 4.数据写入完成，遍历IGotoMap，根据ILabelMap的key，即LabelIndex，在ILabelMap中取得语句偏移Pos，写入[jumpPos参数所在的文件偏移]位置，填充数据。
	ILabelMap map[int]int // LabelIndex(标签序号) -> CodeStartPos(代码开头地址，跳转目标地址)
	IGotoMap  map[int]int // LabelIndex(标签序号) -> GotoParamPos(跳转参数地址)
}

func (e *ScriptEntry) InitEntry() {
	e.ELabelMap = make(map[int]int)
	e.EGotoMap = make(map[int]int)

	e.ILabelMap = make(map[int]int)
	e.IGotoMap = make(map[int]int)

	e.IndexNext = 1
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

func NewScript(opt ScriptFileOptions) *ScriptFile {
	script := new(ScriptFile)
	script.FileName = opt.FileName
	script.GameName = opt.GameName
	script.Version = opt.Version
	script.InitEntry()
	return script
}

// SetOperateParams 设置需要导出的变量值
// 若为导入模式，则会校验读取到的参数列表，与导出的变量值，转换类型并且替换为导入参数
func (s *ScriptFile) SetOperateParams(index int, mode enum.VMRunMode, op string, params ...interface{}) error {
	paramNum := len(params)
	var paramsExport []bool // 导出参数列表
	strCharset := charset.UTF_8
	if paramNum > 0 {
		if val, ok := params[paramNum-1].(charset.Charset); ok { // 最后一个参数为编码类型
			strCharset = val // 编码
			paramNum--
		}
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
				val, _ := strconv.ParseUint(code.Params[i].(string), 16, 8)
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
		CodeParamsToBytes(code, strCharset, allParamList)
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
