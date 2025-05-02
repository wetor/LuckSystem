# Important
This project only accepts **bug issues** and **pull requests**, and does not provide assistance in use  
此项目仅接受现有功能的BUG反馈和Pull requests，不提供使用上的帮助

# Luck System

LucaSystem 引擎解析工具

## 使用方法：[Usage](Usage.md)
## 插件手册：[Plugin](Plugin.md)

## LucaSystem解析完成进度

### Luca Pak 封包文件

- 导出完成
- 导入完成
    - 仅支持替换文件数据

### Luca CZImage 图片文件

#### CZ0

- 导出完成 32位
- 导入完成 32位

#### CZ1

- 导出完成 8位
- 导入完成 8位

#### CZ2

- 导出完成 8位
- 导入完成 8位

#### CZ3

- 导出完成 32位 24位
- 导入完成 32位 24位

#### CZ4

- LucaSystemTools中完成

#### CZ5

- 未遇到

### Luca Script 脚本文件

- 导出完成
- 导入完成
- ~~简单的模拟执行~~
- 支持插件扩展（gpython）
  - 非标准的Python，语法类似Python3.4，缺少大量的内置库和一些特性，基本使用没有问题
  - 插件手册 [Plugin](Plugin.md)

#### 笔记

根据时间，可以LucaSystem的脚本类型分为三个版本，目前仅研究V3版本，即最新版本。LucaSystemTools支持V2版本的脚本解析

| 类型  |  长度 | 名称 | 说明                                 | 
|-----|-------|-----|------------------------------------|
| uint16 |  2  | len | 代码长度                               |
| uint8 |  1  | opcode | 指令索引                               |
| uint8 |  1  | flag | 一个标志，值0~3                          |
| []uint16 |  2 * n  | data0 | 未知参数，其中n=flag(flag<3),n=2(flag==3) |
| params |  len -4 -2*n  | params | 参数                                 |
| uint8 |  k  | align | 补齐位，其中k=len%2                      |

### Luca Font 字体文件

- 解析完成
- 能够简单使用，生成指定文本的图像
- 导出完成
- 导入、制作完成

#### info文件

- 导出完成
- 导入完成

### Luca OggPak 音频封包

- 导出完成

## 目前支持的游戏
1. 《LOOPERS》 Steam
2. LB_EN:《Little Busters! English Edition》 Steam
3. SP:《Summer Pockets》 Nintendo Switch
4. CartagraHD
5. KANON
6. HARMONIA

## 目前支持的指令

- MESSAGE (LB_EN、SP、LOOPERS)
- SELECT (LB_EN、SP)
- IMAGELOAD (LB_EN、SP)

- BATTLE (LB_EN)
- EQU
- EQUN
- EQUV
- ADD
- RANDOM
- IFN
- IFY
- GOTO
- JUMP
- FARCALL
- GOSUB


## 更新日志

### 2.3.2
- 支持 LUNARiA Steam version [@thedanill](https://github.com/thedanill)
- 支持 AIR Steam version [@thedanill](https://github.com/thedanill)
- 支持 Planetarian SG Steam version [@thedanill](https://github.com/thedanill)

### 2.3.1
- 支持 Harmonia FULL HD Steam version [@Mishalac](https://github.com/MishaIac)

### 2.3.0
- 支持 Kanon [@Mishalac](https://github.com/MishaIac)

### 2.2.3
- 支持`-blacklist`命令，添加额外的脚本名黑名单

### 2.2.1 (2023.12.4)
- 支持[CartagraHD](https://vndb.org/r78712)脚本导入导出（未测试）

### 2.2.0 (2023.12.3)
- 支持CZ2的导入（未实际测试）

### 2.1.0 (2023.11.28)
- 支持CZ2的导出

### 2023.10.7
- 支持LOOPERS导入和导出(已测试)
- 支持Plugin扩展以支持任意游戏
- 内置SummerPockets(未测试)和LOOPERS默认Plugin插件和OPCODE
- 移除模拟器相关代码


### 6.26
- 完全重构cmd使用方式
  - 暂不支持script脚本的cmd调用
- 支持24位cz3图像，修复缺少Colorblock值导致的错误
- font插入新字符改为追加替换模式，总字符数增加或保持不变

### 3.15
- 修复cz图像导出时alpha通道异常的问题

### 3.11
- 修复script导入导出交互bug
- 测试部分交互
- 新增Usage文档

### 3.03
- 完整的控制台交互接口（未测试）
- 帮助文档

### 2.17
- 统一cz、info、font、pak、script的接口
- 完善测试用例

### 2.10 
- 统一接口规范

### 2.9
- 修复script导入导出中换行、空行的问题
- Merge AEBus pr
  - 1. Fixed situation when LuckSystem would stop parsing scripts after finding END opcode
  - 2. Added handling of TASK, SAYAVOICETEXT, VARSTR_SET opcodes, and fixed handling of BATTLE opcode.
  - 3. Added opcode names for LB_EN, changed first three opcodes to EQU, EQUN, EQUV as specified in LITBUS_WIN32.exe, added handling of these opcodes in LB_EN.go

### 1.25
- 完成pak导入导出交互

### 1.22
- 完成CZ1导入
- 完成CZ0导出导入
- 支持LB_EN BATTLE指令
- 修正PAK文件ID，与脚本中的ID对应
- 更换日志库为glog
- 引入tui库tview

### 1.21
- 完成LZW压缩
- 完成图像拆分算法
- 支持CZ3格式替换图像

### 2022.1.19

- 支持替换pak文件内容并打包
    - 不支持修改文件名和增加文件
- 不再以LucaSystem引擎模拟器为目标，现以替代LucaSystemTools项目为目标

### 8.13

- 项目更名为LuckSystem
    - 目标为实现LucaSystem引擎的模拟器

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
