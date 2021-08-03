package pak

type FileEntry struct {
	Offset uint32
	Length uint32
	Data   []byte `struct:"-"`
	Name   string `struct:"-"`
	Index  int    `struct:"-"`
}

func (e *FileEntry) OpenScript() {

}
