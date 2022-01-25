package cmd

import (
	"flag"
	"github.com/go-restruct/restruct"
	"lucksystem/charset"
	"strings"
)

type CmdInterface interface {
	Open([]interface{})
	Export([]interface{})
	Import([]interface{})
}

func Cmd() {

	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	//var runType, infile, outfile, mode, config string
	//
	//flag.StringVar(&runType, "type", "", `[required] Input file type
	// game   : Folder of Full Game Giles
	// script : SCRIPT.PAK file
	// font   : FONT.PAK file
	// pak    : *.pak file
	// cz     : cz0 cz1 cz3 file
	//`)
	//flag.StringVar(&infile, "i", "", "[required] Input file or folder")
	//flag.StringVar(&outfile, "o", "", "[required] Output file or folder")
	//
	//flag.StringVar(&mode, "mode", "", `[required] "export" or "import"`)
	//
	//flag.StringVar(&mode, "mode", "", `[required] "export" or "import"`)
	var runType, infile, outfile, mode, coding, config string

	flag.StringVar(&runType, "type", "", `[required] Input file type
  pak    : *.pak file`)
	flag.StringVar(&infile, "i", "", "[required] Input file or folder")
	flag.StringVar(&outfile, "o", "", "[required] Output file or folder")
	flag.StringVar(&mode, "mode", "", `[required] "export" or "import"`)

	flag.StringVar(&config, "config", "", `Parameter
  -type=pak -mode=export [required] : 
    all                   // -o=[Folder]
    index,[Index from 0]  // -o=[File]
    id,[ID in script]     // -o=[File]
    name,[Filename]       // -o=[File]
    // example: "-config=all", "-config=id,10045", "-config=name,_VARSTR"
  -type=pak -mode=import [required] : 
    [import file or folder]
    // example: "-config=/path/to/import", "-config=/path/to/file"
`)

	flag.StringVar(&coding, "charset", "", `[pak] Pak filename charset(utf8 or sjis)`)

	flag.Parse()
	restruct.EnableExprBeta()
	var cmd CmdInterface
	switch runType {
	case "pak":
		c := charset.UTF_8
		switch coding {
		case "utf8", "UTF-8", "UTF_8":
			c = charset.UTF_8
		case "sjis", "SHIFT-JIS", "Shift_JIS":
			c = charset.ShiftJIS
		default:
			panic("Unknown charset")
		}
		cmd = &CmdPak{
			charset: c,
		}
	}

	cmd.Open([]interface{}{infile})

	var param []interface{}
	for _, p := range strings.Split(config, ",") {
		param = append(param, p)
	}
	param = append(param, outfile)
	switch mode {
	case "export":
		cmd.Export(param)
	case "import":
		cmd.Import(param)
	default:
		panic("Unknown mode")
	}
}
