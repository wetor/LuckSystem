package script

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func (e *ScriptEntry) AddExportGotoLabel(codeIndex, pos int) int {

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

func (s *ScriptFile) AddCodeParams(index int, op string, params ...interface{}) {

	paramList := make([]interface{}, 0, len(params))
	for i := 0; i < len(params); i++ {
		switch param := params[i].(type) {
		case []uint16:
			for _, val := range param {
				paramList = append(paramList, val)
			}
		case *JumpParam:
			paramList = append(paramList, param)
			s.AddExportGotoLabel(index, param.Position)
		default:
			paramList = append(paramList, param)

		}
	}
	s.Codes[index].Params = paramList
	s.Codes[index].OpStr = op

}
func (s *ScriptFile) ToStringCodeParams(code *CodeLine) string {
	paramStr := make([]string, 0, len(code.Params))
	for i := 0; i < len(code.Params); i++ {
		switch param := code.Params[i].(type) {
		case []uint16:
			for _, val := range param {
				paramStr = append(paramStr, strconv.FormatInt(int64(val), 10))
			}
		case byte:
			paramStr = append(paramStr, fmt.Sprintf("0x%X", param))
		case string:
			paramStr = append(paramStr, `"`+param+`"`)
		case *JumpParam:
			if code.GotoIndex > 0 {
				if len(param.ScriptName) > 0 {
					paramStr = append(paramStr, fmt.Sprintf(`{goto "%s" label%d}`, param.ScriptName, code.GotoIndex))
				} else {
					paramStr = append(paramStr, fmt.Sprintf("{goto label%d}", code.GotoIndex))
				}
			} else {
				if len(param.ScriptName) > 0 {
					paramStr = append(paramStr, fmt.Sprintf(`{goto "%s" %d}`, param.ScriptName, param.Position))
				} else {
					paramStr = append(paramStr, fmt.Sprintf("{goto %d}", param.Position))
				}
			}
		default:
			paramStr = append(paramStr, fmt.Sprintf("%v", param))

		}
	}
	str := strings.Join(paramStr, ", ")

	if code.LabelIndex > 0 {
		return fmt.Sprintf(`label%d: %s (%s)`, code.LabelIndex, code.OpStr, str)
	} else {
		return fmt.Sprintf(`%s (%s)`, code.OpStr, str)
	}
}

func (s *ScriptFile) ParseCodeParams(index int, codeStr string) {

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
		str := s.ToStringCodeParams(code)
		fmt.Println(i, str)
		fmt.Fprintln(w, str)
	}
	return w.Flush()

}

// Import 导入可编辑脚本
func (s *ScriptFile) Import() {

}

// Save 保存为脚本文件
func (s *ScriptFile) Save() {

}
