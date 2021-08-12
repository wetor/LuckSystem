package font

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestInfo(t *testing.T) {
	restruct.EnableExprBeta()
	list := []string{"info24"}
	for _, name := range list {

		data, _ := os.ReadFile("../data/LB_EN/IMAGE/" + name)
		LoadFontInfo(data)
		fmt.Println()
	}
}

func TestStr(t *testing.T) {
	fmt.Println(len("真人@「…戦いさ」"))
	str := "真人@「…戦いさ」"
	for _, r := range str {
		fmt.Println(r)
	}
}
