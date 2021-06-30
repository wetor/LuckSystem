package charset

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

type Charset string

//中文
const (
	GBK     Charset = "GBK"
	GB18030 Charset = "GB18030"
	GB2312  Charset = "GB2312"
	Big5    Charset = "Big5"
)

//日文
const (
	EUCJP     Charset = "EUCJP"
	ISO2022JP Charset = "ISO2022JP"
	ShiftJIS  Charset = "Shift_JIS"
)

//韩文
const (
	EUCKR Charset = "EUCKR"
)

//Unicode
const (
	UTF_8    Charset = "UTF-8"
	UTF_16   Charset = "UTF-16"
	UTF_16BE Charset = "UTF-16BE"
	UTF_16LE Charset = "UTF-16LE"
	Unicode  Charset = "UTF-16LE"
)

//其他编码
const (
	Macintosh Charset = "macintosh"
	IBM       Charset = "IBM*"
	Windows   Charset = "Windows*"
	ISO       Charset = "ISO-*"
)

var charsetAlias = map[string]string{
	"HZGB2312": "HZ-GB-2312",
	"hzgb2312": "HZ-GB-2312",
	"GB2312":   "HZ-GB-2312",
	"gb2312":   "HZ-GB-2312",
}

func Convert(dstCharset Charset, srcCharset Charset, src []byte) (dst string, err error) {
	if dstCharset == srcCharset {
		return string(src), nil
	}
	dst = string(src)
	// Converting <src> to UTF-8.
	if srcCharset != "UTF-8" {
		if e := getEncoding(srcCharset); e != nil {
			tmp, err := ioutil.ReadAll(
				transform.NewReader(bytes.NewReader(src), e.NewDecoder()),
			)

			if err != nil {
				return "", fmt.Errorf("%s to utf8 failed. %v", srcCharset, err)
			}
			src = tmp
		} else {
			return dst, errors.New(fmt.Sprintf("unsupport srcCharset: %s", srcCharset))
		}
	}
	// Do the converting from UTF-8 to <dstCharset>.
	if dstCharset != "UTF-8" {
		if e := getEncoding(dstCharset); e != nil {
			tmp, err := ioutil.ReadAll(
				transform.NewReader(bytes.NewReader(src), e.NewEncoder()),
			)
			if err != nil {
				return "", fmt.Errorf("utf to %s failed. %v", dstCharset, err)
			}
			dst = string(tmp)
		} else {
			return dst, errors.New(fmt.Sprintf("unsupport dstCharset: %s", dstCharset))
		}
	} else {
		dst = string(src)
	}
	return dst, nil
}

func ToUTF8(srcCharset Charset, src []byte) (dst string, err error) {
	return Convert("UTF-8", srcCharset, src)
}

func UTF8To(dstCharset Charset, src []byte) (dst string, err error) {
	return Convert(dstCharset, "UTF-8", src)
}

func getEncoding(charset Charset) encoding.Encoding {
	if c, ok := charsetAlias[string(charset)]; ok {
		charset = Charset(c)
	}
	if e, err := ianaindex.MIB.Encoding(string(charset)); err == nil && e != nil {
		return e
	}
	return nil
}
