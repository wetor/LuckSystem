package script

type ScriptFile struct {
	FileName string      `struct:"-"`
	GameName string      `struct:"-"`
	Version  uint8       `struct:"-"`
	CodeNum  int         `struct:"-"`
	Codes    []*CodeLine `struct:"while=true"`
}

// Export 导出可编辑脚本
func (s *ScriptFile) Export() {

}

// Import 导入可编辑脚本
func (s *ScriptFile) Import() {

}

// Save 保存为脚本文件
func (s *ScriptFile) Save() {

}
