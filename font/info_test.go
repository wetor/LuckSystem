package font

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}
func TestInfo(t *testing.T) {
	restruct.EnableExprBeta()
	list := []string{"info32"}
	for _, name := range list {
		file := "../data/LB_EN/IMAGE/" + name
		data, _ := os.ReadFile(file)
		info := LoadFontInfo(data)
		txtFile, _ := os.Create(file + ".txt")
		info.Export(txtFile)
		fmt.Println(info.CharNum, " ", len(info.IndexUnicode))
		infoFile, _ := os.Create(file + ".out")
		info.Write(infoFile)

		infoFile.Close()
		txtFile.Close()
	}
}
func TestInfo2(t *testing.T) {
	restruct.EnableExprBeta()
	file := "../data/Other/Font/info32e.info"
	data, _ := os.ReadFile(file)
	info := LoadFontInfo(data)
	txtFile, _ := os.Create(file + ".txt")
	info.Export(txtFile)
	fmt.Println(info.CharNum, " ", len(info.IndexUnicode))
	infoFile, _ := os.Create(file + ".out")
	info.Write(infoFile)

}
func TestStr(t *testing.T) {
	aaa := []int{1, 2, 3}
	index := 3
	bbb := make([]int, 10)
	copy(bbb, aaa[:index])

	fmt.Println(bbb)
}
func TestInfo_Export(t *testing.T) {
	restruct.EnableExprBeta()
	var err error
	savePath := "../data/LB_EN/FONT/"
	infoFiles := []string{"info32", "info24"}

	//============
	for _, name := range infoFiles {
		data, _ := os.ReadFile(savePath + name)
		info := LoadFontInfo(data)
		fmt.Println(name, info.CharNum, len(info.IndexUnicode))
		fs, _ := os.Create(savePath + name + "_export.txt")
		err = info.Export(fs)
		if err != nil {
			panic(err)
		}
		err = fs.Close()
		if err != nil {
			panic(err)
		}
	}

}

func TestInfo_Import(t *testing.T) {
	restruct.EnableExprBeta()
	var err error
	loadPath := "../data/LB_EN/FONT/"
	infoFiles := []string{"info32", "info24"}
	addChars := "!@#$%.,abcdefgAB CDEFG12345测试中文汉字"
	font, _ := os.Open("../data/Other/Font/ARHei-400.ttf")
	defer font.Close()
	//============
	fmt.Println()
	for _, name := range infoFiles {
		data, _ := os.ReadFile(loadPath + name)
		info := LoadFontInfo(data)
		fmt.Println(name, info.CharNum, len(info.IndexUnicode))

		err = info.Import(font, 0, true, "")
		if err != nil {
			panic(err)
		}
		fmt.Println(name, info.CharNum, len(info.IndexUnicode))
		fs, _ := os.Create(loadPath + name + "_byOnlyRedraw")
		err = info.Write(fs)
		if err != nil {
			panic(err)
		}
		err = fs.Close()
		if err != nil {
			panic(err)
		}
	}

	//============
	fmt.Println()
	for _, name := range infoFiles {
		data, _ := os.ReadFile(loadPath + name)
		info := LoadFontInfo(data)
		fmt.Println(name, info.CharNum, len(info.IndexUnicode))

		err = info.Import(font, int(info.CharNum), false, addChars)
		if err != nil {
			panic(err)
		}
		fmt.Println(name, info.CharNum, len(info.IndexUnicode))
		fs, _ := os.Create(loadPath + name + "_byAddChar")
		err = info.Write(fs)
		if err != nil {
			panic(err)
		}
		err = fs.Close()
		if err != nil {
			panic(err)
		}
	}

	//============
	fmt.Println()
	for _, name := range infoFiles {
		data, _ := os.ReadFile(loadPath + name)
		info := LoadFontInfo(data)
		fmt.Println(name, info.CharNum, len(info.IndexUnicode))

		err = info.Import(font, int(info.CharNum), true, addChars)
		if err != nil {
			panic(err)
		}
		fmt.Println(name, info.CharNum, len(info.IndexUnicode))
		fs, _ := os.Create(loadPath + name + "_byAddCharRedraw")
		err = info.Write(fs)
		if err != nil {
			panic(err)
		}
		err = fs.Close()
		if err != nil {
			panic(err)
		}
	}
}
