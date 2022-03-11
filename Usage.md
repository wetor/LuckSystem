# Usage

```
  -v value
  		0 隐藏所有输出
		2 提示信息
		6 调试信息，一些运行时输出
		7 错误信息，不影响运行
		8 错误信息（不panic，完全可忽略）
        log level for V logs   
        
  -type string
        [required] Source file type
          pak   : *.PAK file
          cz    : CZ0,CZ1,CZ3 file
          info  : font info file in FONT.PAK
          font  : font file in FONT.PAK (cz file)
          script: script file in SCRIPT.PAK
  -src string
        [required] Source file
  -mode string
        [required] "export" or "import"
  -output string
        [required] Output file
  -input string
        [optional] Import mode Input file
        
  -charset string
        [pak.optional] Pak filename charset. Default UTF-8
            UTF-8: "UTF-8", "UTF_8", "utf8"
            SHIFT-JIS: "SHIFT-JIS", "Shift_JIS", "sjis"
  -game string
        [script.required] Game name (support LB_EN, SP)
  -info string
        [font.required] Font info file in FONT.PAK
```

## -params=params...

Parameter list. Use "," to split. Example: -params="all,/path/to/save"

### -mode=export

#### -type=pak

| params[0] | mode | string | 可选[all,index,id,name] | Optional [all,index,id,name] |
| --------- | ---- | ------ | ----------------------- | ---------------------------- |

##### params[0]=all

| params[1] | dir  | string | 保存文件夹路径 | Save folder path |
| --------- | ---- | ------ | -------------- | ---------------- |

##### params[0]=index

| params[1] | index | int  | 从0开始的文件序号 | File sequence number from 0 |
| --------- | ----- | ---- | ----------------- | --------------------------- |

##### params[0]=id

| params[1] | id   | int  | 文件唯一ID | File unique ID |
| --------- | ---- | ---- | ---------- | -------------- |

##### params[0]=name

| params[1] | name | string | 文件名 | File name |
| --------- | ---- | ------ | ------ | --------- |

##### Example

```shell
LuckSystem -mode=export -type=pak -charset=utf8 -src=/source/pak/file -output=/output/file/name/list.txt -params=all,/save/dir/path

LuckSystem -mode=export -type=pak -charset=utf8 -src=/source/pak/file -output=/output/file/name -params=index,10

LuckSystem -mode=export -type=pak -charset=utf8 -src=/source/pak/file -output=/output/file/name -params=id,10051

LuckSystem -mode=export -type=pak -charset=utf8 -src=/source/pak/file -output=/output/file/name -params=name,SEEN2005

```

#### -type=script

##### Example

```shell
LuckSystem -v=0 -mode=export -type=script -game=LB_EN -src=/source/script/file -output=/output/file/name.txt
```

#### -type=cz

##### Example

```shell
LuckSystem -mode=export -type=cz -src=/source/cz/file -output=/output/file/name.png

go run . -mode=export -type=cz -src=/source/cz/file -output=/output/file/name.png
```

#### -type=info

##### Example

```shell
LuckSystem -mode=export -type=info -src=/source/info/file -output=/output/file/name/chars.txt
```

#### -type=font

| params[0] | txtFilename | string | 导出的全字符文件 | Exported full character file |
| --------- | ----------- | ------ | ---------------- | ---------------------------- |

##### Example

```shell
LuckSystem -mode=export -type=font -info=/font/info/file -src=/source/font/file -output=/output/file/name.png -params=/output/file/name/chars.txt
```

### -mode=import

#### -type=pak

| params[0] | mode | string | 可选[file,list,dir] | Optional [file,list,dir] |
| --------- | ---- | ------ | ------------------- | ------------------------ |

##### params[0]=file

- 此时input为导入文件，根据类型自动判断是文件名还是id，不支持index

- Input is an imported file. It automatically determines whether it is a file name or ID according to the type. Index is not supported

| params[1] | name | string | 替换指定文件名文件 | Replace the specified file name |
| --------- | ---- | ------ | ------------------ | ------------------------------- |

or

| params[1] | id   | int  | 替换指定id文件 | Replace the specified ID file |
| --------- | ---- | ---- | -------------- | ----------------------------- |

##### params[0]=list

- input为包含多个文件路径的txt文件，按照export.pak.params[0]==all输出的txt格式
- Input is a TXT file containing multiple file paths, in the TXT format output by export.pak.params[0]==all

##### params[0]=dir

- input忽略。导入文件目录，按照Export导出的文件名进行匹配。若存在文件名则用文件名匹配，不存在则用ID匹配

- Input ignored. Import the file directory and match the file name exported by export. If there is a file namis no file name, match it with the ID

| params[1] | dir  | string | 导入文件目录 | Import the file directory |
| --------- | ---- | ------ | ------------ | ------------------------- |

##### Example

```shell
# 将/input/file/name导入到/source/pak/file中替换文件名为SEEN2005的文件并保存到/output/file/name.pak
LuckSystem -mode=import -type=pak -charset=utf8 -src=/source/pak/file -input=/input/file/name -output=/output/file/name.pak -params=file,SEEN2005

# 将/input/file/name导入到/source/pak/file中替换ID为10051的文件并保存到/output/file/name.pak
LuckSystem -mode=import -type=pak -charset=utf8 -src=/source/pak/file -input=/input/file/name -output=/output/file/name.pak -params=file,10051

# 将/input/file/name/list.txt列表文件中的所有文件导入到/source/pak/file中替换同名文件并保存到/output/file/name.pak
LuckSystem -mode=import -type=pak -charset=utf8 -src=/source/pak/file -input=/input/file/name/list.txt -output=/output/file/name.pak -params=list

# 将/import/input/dir文件夹中的所有文件导入到/source/pak/file中替换同名文件并保存到/output/file/name.pak，其中-input参数必须指定已存在文件，但实际并不会使用
LuckSystem -mode=import -type=pak -charset=utf8 -src=/source/pak/file -input=/source/pak/file -output=/output/file/name.pak -params=dir,/import/input/dir
```

#### -type=script

##### Example

```shell
LuckSystem -v=0 -mode=import -type=script -game=LB_EN -src=/source/script/file -input=/input/file/name.txt -output=/output/file/name

```

#### -type=cz

| params[0] | fillSize | bool | 否填充为原cz图像大小 | Fill to original CZ image size |
| --------- | -------- | ---- | -------------------- | ------------------------------ |

##### Example

```shell
LuckSystem -mode=import -type=cz -src=/source/cz/file -input=/input/file/name.png -output=/output/file/name.png -params=true
```

#### -type=info

- input为ttf字体文件
- input is ttf font file

| params[0] | onlyRedraw | bool | 仅使用新字体重绘，不增加新字符 | Redraw with new fonts only, without adding new characters |
| --------- | ---------- | ---- | ------------------------------ | --------------------------------------------------------- |

or

| params[0] | allChar | string | 增加的全字符。不能包含空格和半角逗号 | Added full characters. Cannot contain spaces and half width commas |
| --------- | ------- | ------ | ------------------------------------ | ------------------------------------------------------------ |

| params[1] | startIndex | int  | 开始位置。前面跳过字符数量 | Start position. Number of characters skipped before |
| --------- | ---------- | ---- | -------------------------- | --------------------------------------------------- |

| params[1] | redraw | bool | 是否用新字体重绘startIndex之前的字符 | Redraw characters before startIndex with new font |
| --------- | ------ | ---- | ------------------------------------ | ------------------------------------------------- |

##### Example

```shell
LuckSystem -mode=import -type=info -src=/source/info/file -input=/ttf/font/file.ttf -output=/output/file/name -params=true

LuckSystem -mode=import -type=info -src=/source/info/file -input=/ttf/font/file.ttf -output=/output/file/name -params=测试汉字TestADdChar,6000,false
```

#### -type=font

- input为ttf字体文件

- input is ttf font file

| params[0] | onlyRedraw | bool | 仅使用新字体重绘，不增加新字符 | Redraw with new fonts only, without adding new characters |
| --------- | ---------- | ---- | ------------------------------ | --------------------------------------------------------- |

or

| params[0] | allCharFile | string | 增加的全字符文件，若startIndex==0，且第一个字符不是空格，会自动补充为空格 | Added full characters file. If params[1]==0 and the first character is not a space, it will be automatically supplemented with a space |
| --------- | ----------- | ------ | ------------------------------------------------------------ | ------------------------------------------------------------ |

| params[1] | startIndex | int  | 开始位置。前面跳过字符数量 | Start position. Number of characters skipped before |
| --------- | ---------- | ---- | -------------------------- | --------------------------------------------------- |

| params[1] | redraw | bool | 是否用新字体重绘startIndex之前的字符 | Redraw characters before startIndex with new font |
| --------- | ------ | ---- | ------------------------------------ | ------------------------------------------------- |

##### Example

```shell
LuckSystem -mode=export -type=font -info=/font/info/file -src=/source/font/file -input=/ttf/font/file.ttf -output=/output/file/name -params=true

LuckSystem -mode=export -type=font -info=/font/info/file -src=/source/font/file -input=/ttf/font/file.ttf -output=/output/file/name -params=/chars/file/name.txt,100,false
```

