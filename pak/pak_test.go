package pak

import (
	"fmt"
	"lucascript/charset"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestPak(t *testing.T) {
	restruct.EnableExprBeta()
	pak := NewPak(&PakFileOptions{
		FileName: "../data/SP/SCRIPT.PAK",
		Coding:   charset.ShiftJIS,
	})
	err := pak.Open()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v %v\n", pak.PakHeader, pak.Files[0].Name)

}
