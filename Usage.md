## 使用help能获取详细指令信息

## Example
```shell
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
