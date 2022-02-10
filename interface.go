package main

import "io"

type Interface interface {
	// Export 解包用，提取资源
	//  Description
	//  Param w
	//  Param opt
	//  Return error
	Export(w io.Writer, opt ...interface{}) error
	// Import 打包用，导入、替换资源
	//  Description
	//  Param r
	//  Param opt
	//  Return error
	Import(r io.Reader, opt ...interface{}) error
	// Write 打包用，输出到文件
	//  Description
	//  Param w
	//  Param opt
	//  Return error
	Write(w io.Writer, opt ...interface{}) error
}
