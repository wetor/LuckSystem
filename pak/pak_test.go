package pak

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/go-restruct/restruct"
	"lucksystem/charset"
	"lucksystem/voice"
)

func TestPak(t *testing.T) {
	restruct.EnableExprBeta()
	pak := LoadPak("D:\\Game\\LOOPERS\\LOOPERS\\files\\src\\SCRIPT.PAK", charset.UTF_8)

	fmt.Printf("%v %v\n", pak.Header, pak.Files[0].Name)
	for _, f := range pak.Files {
		fmt.Println(f.ID, f.Name, f.Offset, f.Length)
		fs, _ := os.Create("D:\\Game\\LOOPERS\\LOOPERS\\files\\src\\Unpak\\" + f.Name)
		err := pak.Export(fs, "name", f.Name)
		if err != nil {
			panic(err)
		}
		err = fs.Close()
	}
}

func TestImportPak(t *testing.T) {
	restruct.EnableExprBeta()
	pak := LoadPak("D:\\Game\\LOOPERS\\LOOPERS\\files\\SCRIPT.PAK", charset.UTF_8)

	fmt.Printf("%v %v\n", pak.Header, pak.Files[0].Name)
	for _, f := range pak.Files {
		fmt.Println(f.ID, f.Name, f.Offset, f.Length)
		fs, _ := os.Create("D:\\Game\\LOOPERS\\LOOPERS\\files\\Unpak\\" + f.Name)
		err := pak.Export(fs, "name", f.Name)
		if err != nil {
			panic(err)
		}
		err = fs.Close()
	}
}

func TestVoicePak(t *testing.T) {
	restruct.EnableExprBeta()
	pak := LoadPak(
		"/Volumes/NTFS/Download/Little.Busters.English.Edition/Little Busters! English Edition/files/VOICE0.PAK",
		charset.ShiftJIS,
	)

	fmt.Printf("%v\n", pak.Header)
	// for _, f := range pak.Files {
	// 	fmt.Println(f.ID, f.Name, f.Offset, f.Length)
	// }
	for i := 0; i < 3; i++ {

		e, _ := pak.GetById(i)
		ogg, _ := voice.LoadOggPak(i, e.Data)
		for j, oggf := range ogg.Files {
			fmt.Println(i, j, oggf.SampleRate, oggf.Length)
			f, _ := os.Create("../data/LB_EN/VOICE/" + strconv.Itoa(i) + "_" + strconv.Itoa(j) + ".ogg")
			f.Write(oggf.Data)
			f.Close()
		}

	}

}
func TestCGPak(t *testing.T) {
	restruct.EnableExprBeta()
	pak := LoadPak(
		"../data/LB_EN/BGCG.PAK",
		charset.ShiftJIS,
	)

	fmt.Printf("%v\n", pak.Header)
	// for _, f := range pak.Files {
	// 	fmt.Println(f.ID, f.Name, f.Offset, f.Length)
	// }
	for i := 0; i < 10; i++ {

		e, _ := pak.GetById(i)
		f, _ := os.Create("../data/LB_EN/IMAGE/" + e.Name + ".cz3")
		f.Write(e.Data)
		f.Close()

	}

}
func TestFontPak(t *testing.T) {
	restruct.EnableExprBeta()
	pak := LoadPak(
		"../data/LB_EN/FONT.PAK",
		charset.UTF_8,
	)

	fmt.Printf("%v\n", pak.Header)
	for _, f := range pak.Files {
		fmt.Println(f.ID, f.Name, f.Offset, f.Length)
	}
	//list := []string{"info32", "info36"}
	//for _, name := range list {
	//
	//	e, _ := pak.Get(name)
	//	f, _ := os.Create("../data/LB_EN/IMAGE/" + e.Name)
	//	f.Write(e.Data)
	//	f.Close()
	//
	//}

}
func TestPakReplace(t *testing.T) {
	restruct.EnableExprBeta()
	pak := LoadPak(
		"../data/LB_EN/SCRIPT.PAK",
		charset.ShiftJIS,
	)

	fmt.Printf("%v\n", pak.Header)
	for i, f := range pak.Files {
		if i < 160 {
			continue
		}
		fmt.Println(f.ID, f.Name, f.Offset, f.Length, f.Replace)
	}
	fmt.Printf("==============\n")

	f, _ := os.Open("/Users/wetor/GoProjects/LuckSystem/LuckSystem/data/LB_EN/SCRIPT/_VARSTR")
	pak.SetById(166, f)
	f.Close()
	fs, _ := os.Create("../data/LB_EN/SCRIPT.PAK.out")
	pak.Write(fs)
	fs.Close()
	fmt.Printf("%v\n", pak.Rebuild)
	for i, f := range pak.Files {
		if i < 160 {
			continue
		}
		fmt.Println(f.ID, f.Name, f.Offset, f.Length, f.Replace)
	}
}

func TestPakFindCZ2(t *testing.T) {
	//BGCG 266
	//SYSCG 33
	//SYSCG2 122
	//CHARCG 1910
	restruct.EnableExprBeta()
	pak := LoadPak(
		"/Volumes/NTFS/Download/Little.Busters.English.Edition/Little Busters! English Edition/files/CHARCG.PAK",
		charset.UTF_8,
	)

	fmt.Printf("%v\n", pak.Header)
	for _, f := range pak.Files {

		e, _ := pak.GetById(f.ID)
		fmt.Println(string(e.Data[:3]), f.ID, f.Name, f.Offset, f.Length, f.Replace)

	}
	fmt.Printf("==============\n")
	//e, _ := pak.GetById(10)
	//os.WriteFile("../data/LB_EN/IMAGE/10.cz0", e.Data, 0666)
}
func TestPakFindImage(t *testing.T) {

	restruct.EnableExprBeta()
	pak := LoadPak(
		"/Volumes/NTFS/Download/Little.Busters.English.Edition/Little Busters! English Edition/files/CHARCG.PAK",
		charset.UTF_8,
	)

	fmt.Printf("%v\n", pak.Header)
	for _, f := range pak.Files {

		fmt.Println(f.ID, f.Name, f.Offset, f.Length, f.Replace)

	}
	fmt.Printf("==============\n")
	//for _, f := range pak.Files {
	//	e, _ := pak.GetById(f.ID)
	//
	//	cz, err := czimage.LoadCzImage(e.Data)
	//	if err != nil {
	//		panic(err)
	//	}
	//	cz.Export("../data/LB_EN/IMAGE/char/" + strconv.Itoa(f.ID) + ".png")
	//}

}
func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}
func TestPakFile_Export(t *testing.T) {
	restruct.EnableExprBeta()
	var err error
	savePath := "../data/LB_EN/FONT/"
	nameFiles := []string{"明朝32", "info32"}
	indexFiles := []int{0, 33}
	idFiles := []int{1, 34}
	listFile := "list_byAll.txt"

	pak := LoadPak(
		"../data/LB_EN/FONT.PAK",
		charset.UTF_8,
	)

	//============
	for _, name := range nameFiles {
		fs, _ := os.Create(savePath + name + "_byName")
		err = pak.Export(fs, "name", name)
		if err != nil {
			panic(err)
		}
		err = fs.Close()
		if err != nil {
			panic(err)
		}
	}
	//============
	for _, index := range indexFiles {
		fs, _ := os.Create(savePath + strconv.Itoa(index) + "_byIndex")
		err = pak.Export(fs, "index", index)
		if err != nil {
			panic(err)
		}
		err = fs.Close()
		if err != nil {
			panic(err)
		}
	}
	//============
	for _, id := range idFiles {
		fs, _ := os.Create(savePath + strconv.Itoa(id) + "_byId")
		err = pak.Export(fs, "id", id)
		if err != nil {
			panic(err)
		}
		err = fs.Close()
		if err != nil {
			panic(err)
		}
	}
	//============
	fs, _ := os.Create(savePath + listFile)
	err = pak.Export(fs, "all", savePath)
	if err != nil {
		panic(err)
	}
	err = fs.Close()
	if err != nil {
		panic(err)
	}
}

func TestPakFile_Import(t *testing.T) {
	restruct.EnableExprBeta()

	loadPath := "../data/LB_EN/FONT/"
	nameFiles := []string{"明朝32", "info32"}
	idFiles := []int{1, 34}
	listFile := "list_byAll.txt"

	//===========
	// 混淆
	err := os.Rename(loadPath+"明朝32", loadPath+"明朝32_t")
	if err != nil {
		fmt.Println(err)
	}
	err = os.Rename(loadPath+"明朝24", loadPath+"明朝32")
	if err != nil {
		fmt.Println(err)
	}
	err = os.Rename(loadPath+"明朝32_t", loadPath+"明朝24")
	if err != nil {
		fmt.Println(err)
	}
	//===========

	pak := LoadPak(
		"../data/LB_EN/FONT.PAK",
		charset.UTF_8,
	)

	//============
	for _, name := range nameFiles {
		fs, _ := os.Open(loadPath + name + "_byName")
		err = pak.Import(fs, "file", name)
		if err != nil {
			panic(err)
		}
		err = fs.Close()
		if err != nil {
			panic(err)
		}
	}
	out, _ := os.Create(loadPath + "../FONT_byFileName.PAK")
	err = pak.Write(out)
	if err != nil {
		panic(err)
	}
	err = out.Close()
	if err != nil {
		panic(err)
	}

	//============
	for _, id := range idFiles {
		fs, _ := os.Open(loadPath + strconv.Itoa(id) + "_byId")
		err = pak.Import(fs, "file", id)
		if err != nil {
			panic(err)
		}
		err = fs.Close()
		if err != nil {
			panic(err)
		}
	}
	out, _ = os.Create(loadPath + "../FONT_byFileId.PAK")
	err = pak.Write(out)
	if err != nil {
		panic(err)
	}
	err = out.Close()
	if err != nil {
		panic(err)
	}

	//============
	fs, _ := os.Open(loadPath + listFile)
	err = pak.Import(fs, "list", "")
	if err != nil {
		panic(err)
	}
	err = fs.Close()
	if err != nil {
		panic(err)
	}
	out, _ = os.Create(loadPath + "../FONT_byList.PAK")
	err = pak.Write(out)
	if err != nil {
		panic(err)
	}
	err = out.Close()
	if err != nil {
		panic(err)
	}
	//============
	err = pak.Import(nil, "dir", loadPath)
	if err != nil {
		panic(err)
	}
	out, _ = os.Create(loadPath + "../FONT_byDir.PAK")
	err = pak.Write(out)
	if err != nil {
		panic(err)
	}
	err = out.Close()
	if err != nil {
		panic(err)
	}

}
