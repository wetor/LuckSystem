package cmd

import (
	"fmt"
	"lucksystem/charset"
	"lucksystem/pak"
	"lucksystem/utils"
	"os"
	"path"
	"strconv"
)

type CmdPak struct {
	pak     *pak.PakFile
	charset charset.Charset
}

// Open
//  Description
//  Receiver c *CmdPak
//  Param argv ...interface{} filename
//
func (c *CmdPak) Open(argv []interface{}) {
	c.pak = pak.NewPak(&pak.PakFileOptions{
		FileName: argv[0].(string),
		Coding:   c.charset,
	})
	err := c.pak.Open()
	if err != nil {
		panic(err)
	}
}

// Export
//  Description
//  Receiver c *CmdPak
//  Param argv ...interface{} mode(all,index,id)
//    all: mode,dir
//    index: mode,index,outfile
//    id: mode,id,outfile
//    name: mode,name,outfile
func (c *CmdPak) Export(argv []interface{}) {
	fmt.Println(argv)
	var err error
	switch argv[0].(string) {
	case "all":
		path := argv[1].(string)
		fes := c.pak.ReadAll()
		for _, e := range fes {
			file := ""
			if len(e.Name) != 0 {
				file = path + "/" + e.Name
			} else {
				file = path + "/" + strconv.Itoa(e.ID)
			}
			err = os.WriteFile(file, e.Data, 0666)
		}
	case "index":
		fallthrough
	case "id":
		fallthrough
	case "name":
		var e *pak.FileEntry
		path := argv[2].(string)
		switch argv[0].(string) {
		case "index":
			i, _ := strconv.Atoi(argv[1].(string))
			e, err = c.pak.GetByIndex(i)
		case "id":
			i, _ := strconv.Atoi(argv[1].(string))
			e, err = c.pak.GetById(i)
		case "name":
			e, err = c.pak.Get(argv[1].(string))
		}
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(path, e.Data, 0666)
	default:
		panic("pak.mode error")
	}
	if err != nil {
		panic(err)
	}
}

// Import
//  Description
//  Receiver c *CmdPak
//  Param argv ...interface{} file(or dir),outfile
//
func (c *CmdPak) Import(argv []interface{}) {
	filename := argv[0].(string)
	outfile := argv[1].(string)
	s, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	if s.IsDir() {
		// folder
		files, err := utils.GetDirFileList(filename)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			name := path.Base(file)
			_, has := c.pak.NameMap[name]
			fs, _ := os.Open(file)
			if has {
				err = c.pak.Set(name, fs)
			} else {
				i, err := strconv.Atoi(name)
				if err != nil {
					continue
				}
				err = c.pak.SetById(i, fs)
			}
			fs.Close()
			if err != nil {
				panic(err)
			}
		}

	} else {
		// file
		name := path.Base(filename)
		_, has := c.pak.NameMap[name]
		fs, _ := os.Open(filename)
		if has {
			err = c.pak.Set(name, fs)
		} else {
			i, err := strconv.Atoi(name)
			if err != nil {
				panic(err)
			}
			err = c.pak.SetById(i, fs)
		}
		fs.Close()
		if err != nil {
			panic(err)
		}
	}
	ofs, _ := os.Create(outfile)
	err = c.pak.Write(ofs)
	ofs.Close()
	if err != nil {
		panic(err)
	}
}
