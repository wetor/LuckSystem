package main

import (
	"flag"
	"github.com/go-restruct/restruct"
	"io"
	"lucksystem/charset"
	"lucksystem/czimage"
	"lucksystem/font"
	"lucksystem/game/VM"
	"lucksystem/game/enum"
	"lucksystem/pak"
	"lucksystem/script"
	"os"
	"strconv"
	"strings"
)

func Cmd() {

	var err error
	var pType, pSrc, pInput, pOutput, pMode, pParams string

	flag.StringVar(&pType, "type", "", `[required] Source file type
  pak	: *.PAK file
  cz	: CZ0,CZ1,CZ3 file
  info	: font info file in FONT.PAK
  font	: font file in FONT.PAK (cz file)
  script: script file in SCRIPT.PAK`)
	flag.StringVar(&pSrc, "src", "", "[required] Source file")
	flag.StringVar(&pInput, "input", "", "[optional] Import mode Input file")
	flag.StringVar(&pOutput, "output", "", "[required] Output file")
	flag.StringVar(&pMode, "mode", "", `[required] "export" or "import"`)

	flag.StringVar(&pParams, "params", "", `Parameter list. Use "," to split. Example: -params="all,/path/to/save"
export.type:
  pak:
    params[0]		mode	string	可选 [all,index,id,name] (Optional [all,index,id,name])
      params[0]==all: output为每行一个文件路径的txt文件 (Output is a TXT file with one file path per line)
        params[1]	dir	string	保存文件夹路径 (Save folder path)
      params[0]==index:
        params[1]	index	int	从0开始的文件序号 (File sequence number from 0)
      params[0]==id:
        params[1]	id	int	文件唯一ID (File unique ID)
      params[0]==name:
        params[1]	name	string	文件名 (File name)
  cz:
  info:
  font:
    params[0] 	txtFilename	string 	导出的全字符文件 (Exported full character file)
  script:
import.type:
  pak:
    params[0]		mode	string	可选 [file,list,dir] (Optional [file,list,dir])
      params[0]==file: input为导入文件，根据类型自动判断是文件名还是id，不支持index (Input is an imported file. It automatically determines whether it is a file name or ID according to the type. Index is not supported)
        params[1]	name	string	替换指定文件名文件 (Replace the specified file name)
          or
        params[1]	id	int	替换指定id文件 (Replace the specified ID file)
      params[0]==list: input为包含多个文件路径的txt文件，按照export.pak.params[0]==all输出的txt格式 (Input is a TXT file containing multiple file paths, in the TXT format output by export.pak.params[0]==all)
      params[0]==dir:  input忽略 (Input ignored)
        params[1]	dir	string	导入文件目录，按照Export导出的文件名进行匹配。若存在文件名则用文件名匹配，不存在则用ID匹配 (Import the file directory and match the file name exported by export. If there is a file name, match it with the file name; if there is no file name, match it with the ID)
  cz:
    params[0]		bool	是否填充为原cz图像大小 (Fill to original CZ image size)
  info:
    params[0]	onlyRedraw	bool 	仅使用新字体重绘，不增加新字符 (Redraw with new fonts only, without adding new characters)
      or
    params[0]	allChar	string	增加的全字符。不能包含空格和半角逗号 (Added full characters. Cannot contain spaces and half width commas)
    params[1]	startIndex	int	开始位置。前面跳过字符数量 (Start position. Number of characters skipped before)
    params[2]	redraw	bool	是否用新字体重绘startIndex之前的字符 (Redraw characters before startIndex with new font)
  font:
    params[0]	onlyRedraw	bool	仅使用新字体重绘，不增加新字符 (Redraw with new fonts only, without adding new characters)
      or
    params[0]	allCharFile	string	增加的全字符文件，若startIndex==0，且第一个字符不是空格，会自动补充为空格 (Added full characters file. If params[1]==0 and the first character is not a space, it will be automatically supplemented with a space)
    params[1]	startIndex	int	开始位置。前面跳过字符数量，-1为添加到最后 (Start position. The number of characters skipped before, - 1 is added to the last)
    params[2]	redraw	bool	是否用新字体重绘startIndex之前的字符 (Redraw characters before startIndex with new font)
  script:
`)

	var pCoding string
	flag.StringVar(&pCoding, "charset", "", `[pak.optional] Pak filename charset. Default UTF-8
    UTF-8: "UTF-8", "UTF_8", "utf8"
    SHIFT-JIS: "SHIFT-JIS", "Shift_JIS", "sjis"`)

	var pFontInfo string
	flag.StringVar(&pFontInfo, "info", "", "[font.required] Font info file in FONT.PAK")

	var pGame string
	var pScriptVersion int = 3
	flag.StringVar(&pGame, "game", "", `[script.required] Game name (support LB_EN, SP)`)
	//flag.IntVar(&pScriptVersion, "ver", 3, `[script] Script version, only support 3. Default value 3`)

	flag.Parse()
	restruct.EnableExprBeta()
	var object Interface

	switch pType {
	case "pak":
		coding := charset.UTF_8
		switch pCoding {
		case "utf8", "UTF-8", "UTF_8":
			coding = charset.UTF_8
		case "sjis", "SHIFT-JIS", "Shift_JIS":
			coding = charset.ShiftJIS
		default:
			panic("Unknown charset")
		}
		object = pak.LoadPak(pSrc, coding)
	case "cz":
		object = czimage.LoadCzImageFile(pSrc)
	case "info":
		object = font.LoadFontInfoFile(pSrc)
	case "font":
		object = font.LoadLucaFontFile(pSrc, pFontInfo)
	case "script":
		scr := script.LoadScriptFile(pSrc, pGame, pScriptVersion)
		err = scr.Read()
		if err != nil {
			panic(err)
		}
		var vm *VM.VM
		var output io.Writer
		var input io.Reader
		output, err = os.Create(pOutput)
		if err != nil {
			panic(err)
		}

		if pMode == "export" {
			vm = VM.NewVM(scr, enum.VMRunExport)
			err = vm.LoadOpcode("data/" + pGame + "/OPCODE.txt")
			if err != nil {
				panic(err)
			}
			vm.Run()
			err = scr.Export(output)
		} else if pMode == "import" {
			vm = VM.NewVM(scr, enum.VMRunImport)
			err = vm.LoadOpcode("data/" + pGame + "/OPCODE.txt")
			if err != nil {
				panic(err)
			}
			input, err = os.Open(pInput)
			if err != nil {
				panic(err)
			}
			err = scr.Import(input)
			if err != nil {
				panic(err)
			}
			vm.Run()
			err = scr.Write(output)
		}

		if err != nil {
			panic(err)
		}
		return // type=script直接结束
	default:
		panic("Unknown type")

	}

	var param []interface{}
	for _, p := range strings.Split(pParams, ",") {
		if n, err := strconv.Atoi(p); err == nil {
			param = append(param, n)
		} else if p == "true" {
			param = append(param, true)
		} else if p == "false" {
			param = append(param, false)
		} else {
			param = append(param, p)
		}
	}
	var output io.Writer
	var input io.Reader
	output, err = os.Create(pOutput)
	if err != nil {
		panic(err)
	}

	switch pMode {
	case "export":
		err = object.Export(output, param...)
	case "import":
		input, err = os.Open(pInput)
		if err != nil {
			panic(err)
		}
		err = object.Import(input, param...)
		if err != nil {
			panic(err)
		}
		err = object.Write(output)
	default:
		panic("Unknown mode")
	}
	if err != nil {
		panic(err)
	}
}
