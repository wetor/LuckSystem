package pak

import (
	"encoding/binary"
	"errors"
	"github.com/golang/glog"
	"io"
	"lucksystem/charset"
	"os"
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
//  Param setFileName string
//  Return error
//
func (p *PakFile) Set(name, setFileName string) error {
	id, has := p.NameMap[name]
	if !has {
		return errors.New("文件不存在")
	}
	return p.SetById(id, setFileName)
}
func (p *PakFile) SetById(id int, setFileName string) error {
	return p.SetByIndex(id-int(p.IDStart), setFileName)
}
func (p *PakFile) SetByIndex(index int, setFileName string) error {
	if index < 0 || index >= int(p.FileCount) {
		return errors.New("文件id错误")
	}
	entry := p.Files[index]
	newData, err := os.ReadFile(setFileName)
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

func (p *PakFile) Write(filename string) error {

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
		glog.V(8).Infoln("os.Open", err)
		return err
	}
	defer oldFile.Close()

	file, err := os.Create(filename)
	if err != nil {
		glog.V(8).Infoln("os.Create", err)
		return err
	}
	defer file.Close()

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
