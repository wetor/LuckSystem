package pak

import (
	"encoding/binary"
	"errors"
	"fmt"
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

	Flags      uint32
	NameOffset uint32 `struct:"if=(Flags & 512) != 0 "`
	// 4 * 9 = 36
}

type PakFile struct {
	PakHeader
	Files    []*FileEntry    `struct:"size=FileCount"`
	NameMap  map[string]int  `struct:"-"`
	FileName string          `struct:"-"`
	Coding   charset.Charset `struct:"-"`
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
	data, err := os.ReadFile(p.FileName)
	if err != nil {
		utils.Log("os.ReadFile", err.Error())
		return err
	}
	err = restruct.Unpack(data, binary.LittleEndian, p)
	if err != nil {
		utils.Log("restruct.Unpack", err.Error())
		// return err
	}
	fmt.Println(p.PakHeader)
	// 读取文件名
	named := (p.Flags & 512) != 0
	fmt.Println(p.NameOffset)
	if named {
		offset := int(p.NameOffset)
		size := 0
		for _, file := range p.Files {
			for data[offset+size] != 0x00 {
				size++
			}

			file.Name, err = charset.ToUTF8(p.Coding, data[offset:offset+size])
			if err != nil {
				return err
			}
			// fmt.Println(file.Name, file.Offset, file.Length)
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
		file.Data = data[file.Offset : file.Offset+file.Length]
	}

	return nil
}

func (p *PakFile) Get(name string) (*FileEntry, error) {

	index, has := p.NameMap[name]
	if !has {
		return nil, errors.New("文件不存在")
	}
	return p.Files[index], nil
}
func (p *PakFile) GetById(id int) (*FileEntry, error) {

	if id < 0 || id >= int(p.FileCount) {
		return nil, errors.New("文件id错误")
	}
	return p.Files[id], nil
}
