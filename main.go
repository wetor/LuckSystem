package main

import "flag"

// glog level:
//   1
//   2 提示信息
//   3 vm运行信息
//   4 vm调试信息
//   5 vm错误信息，不影响运行
//   6 调试信息，一些运行时输出
//   7 错误信息，不影响运行
//   8 错误信息（不panic，完全可忽略）

func main() {

	// Export all to folder
	// go run . -type pak -i data/LB_EN/SCRIPT.PAK -o data/LB_EN/test  -charset=sjis  -mode export -config all
	// Export file by index
	// go run . -type pak -i data/LB_EN/SCRIPT.PAK -o data/LB_EN/test/001.scr  -charset=sjis  -mode export -config index,1
	// Export file by name
	// go run . -type pak -i data/LB_EN/SCRIPT.PAK -o data/LB_EN/test/001.scr  -charset=sjis  -mode export -config name,_BUILD_COUNT
	// Folder import Pak
	// go run . -type pak -i data/LB_EN/SCRIPT.PAK -o data/LB_EN/SCRIPT.PAK.out  -charset=sjis  -mode import -config data/LB_EN/test
	// File import Pak
	// go run . -type pak -i data/LB_EN/SCRIPT.PAK -o data/LB_EN/SCRIPT.PAK.out  -charset=sjis  -mode import -config data/LB_EN/test/_CGMODE

	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	Cmd()

}
