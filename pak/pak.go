package pak

import (
	"encoding/binary"
	"errors"
	"lucascript/charset"
	"lucascript/utils"
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
	Unk1         uint32
	BlockSize    uint32

	Unk2 uint32
	Unk3 uint32
	Unk4 uint32
	Unk5 uint32

	Flags uint32
	// 4 * 9 = 36
}

type PakFile struct {
	PakHeader `struct:"-"`
	Files     []*FileEntry    `struct:"size=FileCount"`
	NameMap   map[string]int  `struct:"-"`
	FileName  string          `struct:"-"`
	Coding    charset.Charset `struct:"-"`
}

func NewPak(opt *PakFileOptions) *PakFile {

	pakFile := &PakFile{}
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
		utils.Log("os.Open", err.Error())
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
		utils.Log("restruct.Unpack1", err.Error())
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
	offData := data[tempPos : tempPos+int64(8*p.FileCount)]
	err = restruct.Unpack(offData, binary.LittleEndian, p)
	if err != nil {
		utils.Log("restruct.Unpack2", err.Error())
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
		file.Index = i
		p.NameMap[file.Name] = i
		file.Offset *= p.BlockSize
		// file.Data = data[file.Offset : file.Offset+file.Length]
	}

	return nil
}

func (p *PakFile) Get(name string) (*FileEntry, error) {

	index, has := p.NameMap[name]
	if !has {
		return nil, errors.New("文件不存在")
	}

	f, err := os.Open(p.FileName)
	if err != nil {
		utils.Log("os.Open", err.Error())
		return nil, err
	}
	defer f.Close()
	entry := p.Files[index]
	entry.Data = make([]byte, entry.Length)
	f.ReadAt(entry.Data, int64(entry.Offset))

	return entry, nil
}
func (p *PakFile) GetById(id int) (*FileEntry, error) {

	if id < 0 || id >= int(p.FileCount) {
		return nil, errors.New("文件id错误")
	}
	f, err := os.Open(p.FileName)
	if err != nil {
		utils.Log("os.Open", err.Error())
		return nil, err
	}
	defer f.Close()
	entry := p.Files[id]
	entry.Data = make([]byte, entry.Length)
	f.ReadAt(entry.Data, int64(entry.Offset))

	return entry, nil
}
