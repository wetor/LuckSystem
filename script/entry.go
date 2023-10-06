package script

// Entry 脚本导入导出实体
// 不可用作运行时，需先转话为Script
type Entry struct {
	// 递增序列，从1开始
	IndexNext int

	// 导出：
	// 1.执行脚本，遇到GOTO等跳转指令，将jumpPos作为key，递增序列Index作为value，存入ELabelMap，即一个地址只能存在一个LabelIndex
	// 2.同时，将当前指令的CodeIndex作为key，1中和jumpPos对应的LabelIndex作为value，存入EGotoMap，即标记次条语句包含跳转到LabelIndex的指令
	ELabelMap       map[int]int // Pos(跳转地址) -> LabelIndex(标签序号) ，Pos 通过GOTO (pos) 生成，Index 为序列
	EGotoMap        map[int]int // CodeIndex(代码序号) -> LabelIndex(标签序号) ，CodeIndex为当前语句序列，此语句含有跳转指令，跳转到LabelIndex
	EGlobalLabelMap map[int]int // Pos(跳转地址) -> LabelIndex(标签序号)
	EGlobalGotoMap  map[int]int // CodeIndex(代码序号) -> LabelIndex(标签序号)
	// 导入当前位置
	CurPos int
	// 导入：
	// 1.解析文本，同时开始序列化脚本，转为二进制数据并写入。
	// 2.遇到Label标签，将LabelIndex作为key，当前语句开始位置的文件偏移Pos作为value，存入ILabelMap，即标签对应的跳转地址
	// 3.遇到GOTO等跳转指令时，将要跳转到的LabelIndex作为key，[jumpPos参数所在的文件偏移]作为value存入IGotoMap，即暂时留空，后续再补充数据
	// 4.数据写入完成，遍历IGotoMap，根据ILabelMap的key，即LabelIndex，在ILabelMap中取得语句偏移Pos，写入[jumpPos参数所在的文件偏移]位置，填充数据。
	ILabelMap map[int]int // LabelIndex(标签序号) -> CodeStartPos(代码开头地址，跳转目标地址)
	IGotoMap  map[int]int // GotoParamPos(跳转参数地址) -> LabelIndex(标签序号)

	IGlobalLabelMap map[int]int // LabelIndex(标签序号) -> CodeStartPos(代码开头地址，跳转目标地址)，需要统合所有script的标签
	IGlobalGotoMap  map[int]int // GotoParamPos(跳转参数地址) -> LabelIndex(标签序号)
}

func (e *Entry) InitEntry() {
	e.ELabelMap = make(map[int]int)
	e.EGotoMap = make(map[int]int)
	e.EGlobalLabelMap = make(map[int]int)
	e.EGlobalGotoMap = make(map[int]int)

	e.ILabelMap = make(map[int]int)
	e.IGotoMap = make(map[int]int)
	e.IGlobalLabelMap = make(map[int]int)
	e.IGlobalGotoMap = make(map[int]int)

	e.IndexNext = 1
	e.CurPos = 0
}

func (e *Entry) AddExportGotoLabel(codeIndex, pos int) int {

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
func (e *Entry) AddImportGoto(pos, labelIndex int) {
	e.IGotoMap[pos] = labelIndex
}

// labelIndex, 当前代码位置
func (e *Entry) AddImportLabel(labelIndex, pos int) {
	e.ILabelMap[labelIndex] = pos
}

func (e *Entry) SetGlobalLabel(labels map[int]int) {
	e.EGlobalLabelMap = labels
}

func (e *Entry) SetGlobalGoto(gotos map[int]int) {
	e.EGlobalGotoMap = gotos
}

func (e *Entry) AddImportGlobalGoto(pos, labelIndex int) {
	e.IGlobalGotoMap[pos] = labelIndex
}

func (e *Entry) AddImportGlobalLabel(labelIndex, pos int) {
	e.IGlobalLabelMap[labelIndex] = pos
}

func (e *Entry) SetImportGlobalLabel(labels map[int]int) {
	e.IGlobalLabelMap = labels
}
