package pak

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"io"
	"lucksystem/charset"
	"lucksystem/utils"
	"os"
	"path"
	"strconv"

	"github.com/go-restruct/restruct"
)

type PakFileOptions struct {
	FileName string
	Coding   charset.Charset
}

type PakHeader struct {
	HeaderLength uint32
	FileCount    uint32
	IDStart      uint32 // 图像包中，id段开始
	BlockSize    uint32

	Unk2 uint32
	Unk3 uint32
	Unk4 uint32
	Unk5 uint32

	Flags uint32
	// 4 * 9 = 36
}
type PakType string

type PakFile struct {
	PakHeader `struct:"-"`
	Files     []*FileEntry    `struct:"size=FileCount"`
	NameMap   map[string]int  `struct:"-"`
	FileName  string          `struct:"-"`
	Coding    charset.Charset `struct:"-"`
	OffsetPos int64           `struct:"-"` // Files 数据开始位置
	DataPos   int64           `struct:"-"` // Files.Data 数据开始位置
	Rebuild   bool            `struct:"-"` // 替换数据后，是否需要重构pak
}

func NewPak(opt *PakFileOptions) *PakFile {

	pakFile := &PakFile{
		Rebuild: false,
	}
	pakFile.FileName = opt.FileName
	if opt.Coding != "" {
		pakFile.Coding = opt.Coding
	} else {
		pakFile.Coding = charset.UTF_8
	}
	return pakFile
}

func (p *PakFile) Open() error {
	f, err := os.Open(p.FileName)

	if err != nil {
		glog.V(8).Infoln("os.Open", err)
		return err
	}
	defer f.Close()

	headerLenBytes := make([]byte, 4)
	f.ReadAt(headerLenBytes, 0)
	headerLen := binary.LittleEndian.Uint32(headerLenBytes)

	data := make([]byte, headerLen)
	f.ReadAt(data, 0)

	err = restruct.Unpack(data, binary.LittleEndian, &p.PakHeader)
	if err != nil {
		glog.V(8).Infoln("restruct.Unpack1", err)
		return err
	}

	temp := make([]byte, 4)
	tempPos := int64(32)
	f.ReadAt(temp, tempPos)
	for binary.LittleEndian.Uint32(temp) != p.HeaderLength/p.BlockSize {
		tempPos += 4
		f.ReadAt(temp, tempPos)
	}

	// 文件偏移 长度读取
	p.OffsetPos = tempPos
	offData := data[tempPos : tempPos+int64(8*p.FileCount)]
	err = restruct.Unpack(offData, binary.LittleEndian, p)
	if err != nil {
		glog.V(8).Infoln("restruct.Unpack2", err)
		return err
	}
	// 读取文件名
	named := (p.Flags & 512) != 0
	if named {

		f.ReadAt(temp, tempPos-4)
		offset := int(binary.LittleEndian.Uint32(temp))
		size := 0
		for _, file := range p.Files {
			for data[offset+size] != 0x00 {
				size++
			}
			file.Name, err = charset.ToUTF8(p.Coding, data[offset:offset+size])

			if err != nil {
				return err
			}
			//fmt.Println(file.Name, file.Offset, file.Length)
			offset += size + 1
			size = 0

		}
	}

	p.NameMap = make(map[string]int)
	// 读取文件数据
	for i, file := range p.Files {
		if !named {
			file.Name = strconv.Itoa(i)
		}

		file.ID = int(p.IDStart) + i

		p.NameMap[file.Name] = file.ID
		file.Offset *= p.BlockSize
		file.Replace = false
		// file.Data = data[file.Offset : file.Offset+file.Length]

		if i == 0 {
			p.DataPos = int64(file.Offset) // 第一个文件数据的位置
		}
	}

	return nil
}
func (p *PakFile) ReadAll() []*FileEntry {
	f, err := os.Open(p.FileName)
	if err != nil {
		glog.V(8).Infoln("os.Open", err)
		return nil
	}
	defer f.Close()
	for _, e := range p.Files {
		e.Data = make([]byte, e.Length)
		f.ReadAt(e.Data, int64(e.Offset))
	}
	return p.Files
}

func (p *PakFile) Get(name string) (*FileEntry, error) {

	id, has := p.NameMap[name]
	if !has {
		return nil, errors.New("文件不存在")
	}
	return p.GetById(id)
}

func (p *PakFile) GetById(id int) (*FileEntry, error) {
	return p.GetByIndex(id - int(p.IDStart))
}
func (p *PakFile) GetByIndex(index int) (*FileEntry, error) {

	if index < 0 || index >= int(p.FileCount) {
		return nil, errors.New("文件id错误")
	}

	entry := p.Files[index]
	if entry.Offset == 0 && entry.Data != nil && len(entry.Data) > 0 {
		// 外部数据
		return entry, nil
	}
	f, err := os.Open(p.FileName)
	if err != nil {
		glog.V(8).Infoln("os.Open", err)
		return nil, err
	}
	defer f.Close()

	entry.Data = make([]byte, entry.Length)
	f.ReadAt(entry.Data, int64(entry.Offset))

	return entry, nil
}

// Set 设置外部文件替换pak文件
//  Description
//  Receiver p *PakFile
//  Param name string
//  Param r io.Reader
//  Return error
//
func (p *PakFile) Set(name string, r io.Reader) error {
	id, has := p.NameMap[name]
	if !has {
		return errors.New("文件不存在")
	}
	return p.SetById(id, r)
}
func (p *PakFile) SetById(id int, r io.Reader) error {
	return p.SetByIndex(id-int(p.IDStart), r)
}
func (p *PakFile) SetByIndex(index int, r io.Reader) error {
	if index < 0 || index >= int(p.FileCount) {
		return errors.New("文件id错误")
	}
	entry := p.Files[index]
	newData, err := io.ReadAll(r)
	if err != nil {
		glog.V(8).Infoln("os.ReadFile", err)
		return err
	}
	alignLength := entry.Length
	if alignLength/p.BlockSize*p.BlockSize != alignLength {
		alignLength = (alignLength/p.BlockSize + 1) * p.BlockSize
	}

	p.Rebuild = alignLength < uint32(len(newData)) // 无法装下新数据内容，需要重构

	entry.Replace = true
	entry.Data = newData
	entry.Length = uint32(len(newData))
	return nil
}

// Write
//  Description
//  Receiver p *PakFile
//  Param w io.Writer 必须是*File类型
//  Param opt ...interface{}
//  Return error
//
func (p *PakFile) Write(w io.Writer, opt ...interface{}) error {

	oldOffset := make(map[int]uint32, p.FileCount)
	if p.Rebuild {
		offset := uint32(0)
		flag := false
		for i, f := range p.Files {
			if flag {
				if offset/p.BlockSize*p.BlockSize != offset {
					offset = (offset/p.BlockSize + 1) * p.BlockSize
				}
				oldOffset[i] = f.Offset //保存旧的offset
				f.Offset = offset       // 更新Offset
				offset += f.Length
			}

			if f.Replace { // 资源发生替换，从下一个开始，重新计算offset
				flag = true
				offset = f.Offset + f.Length
			}
		}
	}
	oldFile, err := os.Open(p.FileName)
	if err != nil {
		glog.V(8).Infoln("os.Open.oldFile", err)
		return err
	}
	defer oldFile.Close()

	file, ok := w.(*os.File)
	if !ok {
		glog.V(8).Infoln("w.(*os.File)", err)
		return err
	}

	// 1. 复制文件全部内容
	_, err = io.Copy(file, oldFile)
	if err != nil {
		glog.V(8).Infoln("io.Copy", err)
		return err
	}
	// 2. 写入偏移和长度和数据
	temp := make([]byte, 4)
	flag := false
	if p.Rebuild {
		for i, f := range p.Files {
			if f.Replace { // 资源发生替换，从这里开始写入数据
				flag = true
			}
			if flag {
				binary.LittleEndian.PutUint32(temp, f.Offset/p.BlockSize)
				file.WriteAt(temp, p.OffsetPos+int64(i*8))
				binary.LittleEndian.PutUint32(temp, f.Length)
				file.WriteAt(temp, p.OffsetPos+int64(i*8+4))
				if f.Replace {
					file.WriteAt(f.Data, int64(f.Offset))
				} else {
					data := make([]byte, f.Length)
					oldFile.ReadAt(data, int64(oldOffset[i]))
					file.WriteAt(data, int64(f.Offset))
				}
			}
		}
	} else {
		for i, f := range p.Files {
			if f.Replace {
				binary.LittleEndian.PutUint32(temp, f.Length)
				file.WriteAt(temp, p.OffsetPos+int64(i*8+4))
				file.WriteAt(f.Data, int64(f.Offset))
			}
		}
	}

	return nil
}

// Export Pak文件解包接口
//  Description
//  Receiver p *PakFile
//  Param w io.Writer 保存到的文件名
//  Param opt ...interface{}
//    opt[0] 		mode 	string	可选 all,index,id,name
//      mode=="all": w为每行一个文件路径的txt文件
//        opt[1] 	dir 	string 	保存文件夹路径
//      mode=="index":
//        opt[1] 	index 	int 	从0开始的文件序号
//      mode=="id":
//        opt[1]	id		int		文件唯一ID
//      mode=="name":
//        opt[1]	name	string	文件名
//  Return error
//
// TODO 未测试
func (p *PakFile) Export(w io.Writer, opt ...interface{}) error {
	var err error
	switch opt[0].(string) {
	case "all":
		dir := opt[1].(string)
		fes := p.ReadAll()
		for _, e := range fes {
			file := ""
			line := ""
			if len(e.Name) != 0 {
				file = path.Join(dir, e.Name)
				line = fmt.Sprintf("name:%s,%s", e.Name, file)
			} else {
				file = path.Join(dir, strconv.Itoa(e.ID))
				line = fmt.Sprintf("id:%d,%s", e.ID, file)
			}
			_, err = w.Write([]byte(line))
			if err != nil {
				return err
			}
			err = os.WriteFile(file, e.Data, 0666)
			if err != nil {
				return err
			}
		}
	case "index":
		fallthrough
	case "id":
		fallthrough
	case "name":
		var e *FileEntry
		switch opt[0].(string) {
		case "index":
			e, err = p.GetByIndex(opt[1].(int))
		case "id":
			e, err = p.GetById(opt[1].(int))
		case "name":
			e, err = p.Get(opt[1].(string))
		}
		if err != nil {
			return err
		}
		_, err = w.Write(e.Data)
		if err != nil {
			return err
		}
	default:
		return errors.New("pak.mode error")
	}
	return nil
}

// Import
//  Description
//  Receiver p *PakFile
//  Param r io.Reader 导入文件
//  Param opt ...interface{}
//    opt[0]		mode	string	可选file,list,dir
//      mode=="file": r为导入文件，opt[1]和opt[2]输入其一即可，若都有值，优先opt[1]
//        opt[1]	name	string	替换指定文件名文件
//        opt[2]	id		int		替换指定id文件
//      mode=="list": r为包含多个文件路径的txt文件，按照Export mode=="all"输出的txt格式
//      mode=="dir":  r为空
//        opt[1]	dir		string	导入文件目录，按照Export导出的文件名进行匹配。若存在文件名则用文件名匹配，不存在则用ID匹配
//  Return error
//
// TODO 未测试
func (p *PakFile) Import(r io.Reader, opt ...interface{}) error {

	var err error
	switch opt[0].(string) {
	case "file":
		if opt[1] != nil && len(opt[1].(string)) != 0 {
			name := opt[1].(string)
			err = p.Set(name, r)
		} else if opt[2] != nil {
			index := opt[2].(int)
			err = p.SetById(index, r)
		} else {
			glog.Fatalln("输入参数有误")
		}
	case "list":
		var t, name, file string
		var id int
		var fs *os.File

		scan := bufio.NewScanner(r)
		for scan.Scan() {
			line := scan.Text()
			_, err = fmt.Sscanf(line, "%s:%s,%s", &t, &name, &file)
			if err != nil {
				glog.V(2).Infoln(err)
				continue
			}
			fs, err = os.Open(file)
			if err != nil {
				glog.V(2).Infof("%v %s\n", err, file)
				continue
			}
			if t == "id" {
				_, err = fmt.Sscanf(line, "id:%d,%s", &id, &file)
				err = p.SetById(id, fs)
			} else {
				err = p.Set(name, fs)
			}
			if err != nil {
				return err
			}
			err = fs.Close()
			if err != nil {
				return err
			}
		}
		if err = scan.Err(); err != nil {
			glog.Fatalln("pak.Import.list 输入文件格式有误")
		}

	case "dir":
		var files []string
		var id int
		files, err = utils.GetDirFileList(opt[1].(string))
		if err != nil {
			return err
		}
		for _, file := range files {
			name := path.Base(file)
			_, has := p.NameMap[name]
			fs, _ := os.Open(file)
			if has {
				err = p.Set(name, fs)
			} else {
				id, err = strconv.Atoi(name)
				if err != nil {
					continue
				}
				err = p.SetById(id, fs)
			}
			if err != nil {
				glog.V(2).Infof("%v %s\n", err, file)
				return err
			}
			err = fs.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
