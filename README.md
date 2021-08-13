# Luck System

LucaSystem engine galgame **Emulator**  
LucaSystem Gal引擎的模拟器

## LucaSystem解析

### Luca Pak 封包文件


### Luca CZImage 图片文件

#### CZ0
- [ ] LucaSystemTools中完成
#### CZ1
- [x] 完成
#### CZ2
- [ ] 遇到问题，未解决
#### CZ3
- [x] 完成
#### CZ4
- [ ] LucaSystemTools中完成
#### CZ5
- [ ] 未遇到


### Luca Script 脚本文件
根据时间，可以LucaSystem的脚本类型分为三个版本，目前仅研究V3版本，即最新版本。LucaSystemTools支持V2版本的脚本解析

| 类型  |  长度 | 名称 |  说明 | 
|------|-------|-----|-----|
| uint16 |  2  | len | 代码长度|
| uint8 |  1  | opcode | 指令索引|
| uint8 |  1  | flag | 一个标志，值0~3|
| []uint16 |  2 * n  | data0 | 未知参数，其中n=flag(flag<3),n=2(flag==3)|
| params |  len -4 -2*n  | params | 参数|
| uint8 |  k  | align | 补齐位，其中k=len%2|

### Luca Font 字体文件
#### info文件

### Luca OggPak 音频封包


## 目前支持的游戏

1. SP:《Summer Pockets》 Nintendo Switch
2. LB_EN:《Little Busters! English Edition》 Steam

## 目前支持的指令

- MESSAGE (LB_EN、SP)
- SELECT (LB_EN、SP)
- IMAGELOAD (LB_EN、SP)

- UNKNOW0 (仅LB_EN出现)
- EQU
- EQUN
- ADD
- RANDOM
- IFN
- IFY
- GOTO
- JUMP
- FARCALL
- MOVE

## 更新日志

### 8.12
- 支持字库的加载
  - 字库info文件的解析与应用
  - 字库CZ1图像的解析
- 现已支持根据文字内容，按指定字体生成文字图像

### 8.11
- 支持动态加载pak中的文件
  - 加载pak仅加载pak文件头，内部文件需要时读取
- 支持音频文件的oggpak的解包
- 开始编写CZ图像解析
  - 完成通用lzw解压
  - 支持CZ3图像的加载


### 8.7
- 完美支持脚本导出为文本、导入为脚本
- 开始设计与编写模拟器主体

### 8.3
- 支持pak文件的加载

### 8.1
- 完成大部分导出模式功能
  - 解析文本
  - 合并导出参数和原脚本参数
  - 将文本中的数据合并到原脚本，并转为字节数据

### 7.28
- 完善导出模式，支持更多指令

### 7.27 累积
- 为虚拟机增加导入模式和导出模式
  - 导出模式：不执行引擎层代码，将脚本转为字符串并导出
  - 导入模式：开始设计与编写



### 7.13
- 增加engine结构，即引擎层，与虚拟机做区分
  - 虚拟机：执行脚本内容，保存、计算变量等逻辑相关操作
  - 引擎：执行模拟器的显示、交互等

### 7.12
- 支持表达式计算
  - 表达式的读取以及中缀表达式转后缀表达式
  - 后缀表达式的计算
- 引擎中使用内置数据类型，不在使用包装数据类型


### 7.11
- 重构代码结构，使用vm来处理脚本执行相关
- 增加context，在执行中传递变量表等数据
- 增加变量表，储存运行时变量
- 优化参数的读取
- 统一接口代码，虚拟机与引擎前端交互接口

### 6.30
- 支持多游戏
- 设计参数、函数等结构

### 6.28
- 框架设计与编写
- 第三方包的选择与测试
- 支持LB_EN基本解析

### 计划
- 支持更多LucaSystem引擎的游戏脚本解析
- 完善引擎函数
- 引擎层交互的初步实现