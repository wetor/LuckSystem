## 使用此工具进行翻译工作时，请注意

- 最好需要提取OPCODE，从可执行文件中获取。如无法提取，需要进行以下工作：
  - 反编译脚本，得到全为uint16的脚本文件，找出其中可能包含字符串的操作（一般特别长，并且存在连续的、较大的值）
  - 使用Plugin中的`opcode_dict`功能，映射为插件函数并尝试进行解析，解析出字符串并进行翻译
  - 此时的字符串不能进行超长，且需要使用**全角空格**(半角字符串则使用半角空格)进行填充为原始长度，否则会导致导入后的脚本无法正常使用
  - 如想进行任意字符串的修改，需要识别出[base](data/base)下插件文件中的跳转操作，如`IFN` `FARCALL` `JUMP`等操作，并使用`read_jump`准确的解析出跳转的类型、跳转值等
- 如已经有完整的OPCODE:
  - 反编译脚本，得到全为uint16的脚本文件，找出其中可能包含字符串的操作（一般特别长，并且存在连续的、较大的值）
  - 使用Plugin进行尝试解析
  - 此时的字符串不能进行超长，且需要使用**全角空格**(半角字符串则使用半角空格)进行填充为原始长度，否则会导致导入后的脚本无法正常使用
  - 如想进行任意字符串的修改，需要解析[base](data/base)下插件文件中的跳转操作，如`IFN` `FARCALL` `JUMP`，并使用`read_jump`准确的解析出跳转的类型、跳转值等
- 导入和导出必须使用同一份相同的原SCRIPT.PAK、OPCODE和插件，对插件的任何修改都需要重新进行反编译，才可以导入
- 建议：编写额外的工具，从反编译后的脚本中提取需要翻译的文本，并在翻译完成过使用工具替换到反编译后的脚本中，然后再导入，防止游戏数值被意外的修改
## 使用help能获取详细指令信息

## Example
```shell
# 反编译SCRIPT.PAK，额外忽略_build_time, TEST脚本
lucksystem script decompile \
  -s D:/Game/LOOPERS/files/SCRIPT.PAK \
  -c UTF-8 \
  -O data/LOOPERS.txt \
  -p data/LOOPERS.py \
  -o D:/Game/LOOPERS/files/Export
  -b _build_time,TEST

# 导出CZ2图片
lucksystem  image export \
  -i C:/Users/wetor/Desktop/Prototype/CZ2/32/明朝32 \
  -o C:/Users/wetor/Desktop/Prototype/CZ2/32/明朝32.png

# 反编译SCRIPT.PAK
lucksystem script decompile \
  -s D:/Game/LOOPERS/files/SCRIPT.PAK \
  -c UTF-8 \
  -O data/LOOPERS.txt \
  -p data/LOOPERS.py \
  -o D:/Game/LOOPERS/files/Export

# lucksystem script decompile -s D:/Game/LOOPERS/LOOPERS/files/src/SCRIPT.PAK -c UTF-8 -O data/LOOPERS.txt -p data/LOOPERS.py -o D:/Game/LOOPERS/LOOPERS/files/Export

# 导入修改后的反编译脚本到SCRIPT.PAK
lucksystem script import \
  -s D:/Game/LOOPERS/files/SCRIPT.PAK \
  -c UTF-8 \
  -O data/LOOPERS.txt \
  -p data/LOOPERS.py \
  -i D:/Game/LOOPERS/files/Export \
  -o D:/Game/LOOPERS/files/Import/SCRIPT.PAK

# lucksystem script import -s D:/Game/LOOPERS/LOOPERS/files/src/SCRIPT.PAK -c UTF-8 -O data/LOOPERS.txt -p data/LOOPERS.py -i D:/Game/LOOPERS/LOOPERS/files/Export -o D:/Game/LOOPERS/LOOPERS/files/Import/SCRIPT.PAK


# 查看FONT.PAK文件列表
lucksystem pak \
  -s data/LB_EN/FONT.PAK \
  -L

# 提取FONT.PAK中所有文件到temp中
lucksystem pak extract \
  -i data/LB_EN/FONT.PAK \
  -o data/LB_EN/FONT.txt \
  --all data/LB_EN/temp

# 提起FONT.PAK中第6个(index从零开始)个
lucksystem pak extract \
  -i data/LB_EN/FONT.PAK \
  -o data/LB_EN/5 \
  --index 5

# 替换FONT.PAK中temp内对应文件
lucksystem pak replace \
  -s data/LB_EN/FONT.PAK \
  -o data/LB_EN/FONT.out.PAK \
  -i data/LB_EN/temp

# 替换FONT.PAK中文件名为info32的文件
lucksystem pak replace \
  -s data/LB_EN/FONT.PAK \
  -o data/LB_EN/FONT.out.PAK \
  -i data/LB_EN/temp/info32 \
  --name info32

# 替换FONT.PAK中存在FONT.txt列表的文件
lucksystem pak replace \
  -s data/LB_EN/FONT.PAK \
  -o data/LB_EN/FONT.out.PAK \
  -i data/LB_EN/FONT.txt \
  --list

# 提取cz到png图片
lucksystem image export \
  -i data/LB_EN/FONT/明朝32 \
  -o data/LB_EN/0.png

# 导入png图片到cz
lucksystem image import \
  -s data/LB_EN/FONT/明朝32 \
  -i data/LB_EN/0.png \
  -o data/LB_EN/0.cz1

# 提取32号"明朝"字体图片和字符集txt
lucksystem font extract \
  -s data/Other/Font/明朝32 \
  -S data/Other/Font/info32 \
  -o data/Other/Font/明朝32_f.png \
  -O data/Other/Font/info32_f.txt

# 用ttf重绘32号"明朝"字体
lucksystem font edit \
  -s data/LB_EN/FONT/明朝32 \
  -S data/LB_EN/FONT/info32 \
  -f data/Other/Font/ARHei-400.ttf \
  -o data/Other/Font/明朝32 \
  -O data/Other/Font/info32 \
  -r

# 将allchar.txt中字符追加到明朝32字体最后方
lucksystem font edit \
  -s data/LB_EN/FONT/明朝32 \
  -S data/LB_EN/FONT/info32 \
  -f data/Other/Font/ARHei-400.ttf \
  -o data/Other/Font/明朝32 \
  -O data/Other/Font/info32 \
  -c data/Other/Font/allchar.txt \
  -a

# 将allchar.txt中字符替换到到明朝32字体第7106个字符处
lucksystem font edit \
  -s data/LB_EN/FONT/明朝32 \
  -S data/LB_EN/FONT/info32 \
  -f data/Other/Font/ARHei-400.ttf \
  -o data/Other/Font/明朝32 \
  -O data/Other/Font/info32 \
  -c data/Other/Font/allchar.txt \
  -i 7105

```
