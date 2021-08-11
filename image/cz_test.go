package czimage

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/go-restruct/restruct"
)

func TestCZ3(t *testing.T) {
	restruct.EnableExprBeta()
	for i := 0; i < 10; i++ {
		data, _ := os.ReadFile("../data/LB_EN/IMAGE/" + strconv.Itoa(i) + ".cz3")
		cz, err := LoadCzImage(i, data)
		if err != nil {
			panic(err)
		}
		cz.Save("../data/LB_EN/IMAGE/" + strconv.Itoa(i) + ".png")
	}

}

func TestLineDiff(t *testing.T) {
	fmt.Println(1)
}
