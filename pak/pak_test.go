package pak

import (
	"fmt"
	"lucascript/charset"
	"lucascript/voice"
	"os"
	"strconv"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestPak(t *testing.T) {
	restruct.EnableExprBeta()
	pak := NewPak(&PakFileOptions{
		FileName: "../data/LB_EN/SCRIPT.PAK",
		Coding:   charset.ShiftJIS,
	})
	err := pak.Open()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v %v\n", pak.PakHeader, pak.Files[0].Name)
	for _, f := range pak.Files {
		fmt.Println(f.Index, f.Name, f.Offset, f.Length)
	}
}

func TestVoicePak(t *testing.T) {
	restruct.EnableExprBeta()
	pak := NewPak(&PakFileOptions{
		FileName: "/Volumes/NTFS/Download/Little.Busters.English.Edition/Little Busters! English Edition/files/VOICE0.PAK",
		Coding:   charset.ShiftJIS,
	})
	err := pak.Open()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", pak.PakHeader)
	// for _, f := range pak.Files {
	// 	fmt.Println(f.Index, f.Name, f.Offset, f.Length)
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
	pak := NewPak(&PakFileOptions{
		FileName: "../data/LB_EN/BGCG.PAK",
		Coding:   charset.ShiftJIS,
	})
	err := pak.Open()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", pak.PakHeader)
	// for _, f := range pak.Files {
	// 	fmt.Println(f.Index, f.Name, f.Offset, f.Length)
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
	pak := NewPak(&PakFileOptions{
		FileName: "../data/LB_EN/FONT.PAK",
		Coding:   charset.UTF_8,
	})
	err := pak.Open()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", pak.PakHeader)
	for _, f := range pak.Files {
		fmt.Println(f.Index, f.Name, f.Offset, f.Length)
	}
	// list := []string{"info20", "info24", "明朝24", "明朝20"}
	// for _, name := range list {

	// 	e, _ := pak.Get(name)
	// 	f, _ := os.Create("../data/LB_EN/IMAGE/" + e.Name + ".cz1")
	// 	f.Write(e.Data)
	// 	f.Close()

	// }

}
