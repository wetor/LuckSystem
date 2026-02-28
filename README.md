# LuckSystem 2.3.2 — Yoremi Fork

Fork de [LuckSystem](https://github.com/wetor/LuckSystem) avec corrections de bugs, support de nouveaux formats, et interface graphique pour la traduction de visual novels Visual Art's/Key.

Fork of [LuckSystem](https://github.com/wetor/LuckSystem) with bug fixes, new format support, and graphical interface for Visual Art's/Key visual novel translation.

---

## Supported engines / Moteurs supportés

ProtoDB / LUCA System — AIR, CLANNAD, Kanon, Little Busters, Summer Pockets, Harmonia, LOOPERS, LUNARiA, Planetarian, etc.

---

## GUI

A graphical interface is available in a separate repository:
**[LuckSystem-2.3.2-Yoremi-Update + GUI](https://github.com/yoremi-trad-fr/LuckSystem-2.3.2-Yoremi-Update)** — Built with Wails (Go + Svelte), calls `lucksystem.exe` via subprocess.

### GUI Features
- **Game presets** / Auto-detect available games from data/ folder (OPCODE + plugin auto-fill)
- Dialogue Extract / Extract translatable dialogue from decompiled scripts to TSV (single file or batch)
- Dialogue Import / Reimport translated dialogue from TSV back into scripts (single file or batch)
- Script Decompile / Compile
- PAK Extract / Replace (CG and Font workflows separated)
- Font Extract / Edit (append, insert, redraw modes)
- Image Export / Import (single file + batch folder mode)
- Real-time console output
- **Stop button** to cancel any running operation
- No CMD popup window during batch operations
- Auto-detection of `lucksystem.exe`

> Place `LuckSystemGUI.exe` in the same folder as `lucksystem.exe` to use.

### Linux

Une version Linux est disponibleen binaires séparés (GUI + CLI). Voir les releases pour le téléchargement.

A Linux version is available as separate binaries (GUI + CLI). See the releases for download.

---

## Patches

### Version 3.1.3 — Patch 3 *(latest)*

20. **GUI: Game preset auto-scan from data/ folder** — `app.go`, `frontend/src/App.svelte`, `wailsjs/go/main/App.js`, `App.d.ts`
    - `ScanGameData()` scans `data/` next to lucksystem, discovers all OPCODE `.txt` files (recursive, excluding `base/`)
    - Dynamic "Game preset" dropdown replaces static LB_EN/SP selector — auto-fills Opcode, Plugin and Game fields
    - 9 presets detected: AIR, CartagraHD, HARMONIA, KANON, LB_EN, LOOPERS, LUNARiA, PlanetarianSG, SP

### Version 3.1.3 — Patch 2

19. **`--game`/`-g` flag for forced game type (CLI + GUI)** — `cmd/script.go`, `cmd/scriptDecompile.go`, `cmd/scriptImport.go`, `app.go`, `frontend/src/App.svelte`
    - CLI: new persistent flag `--game`/`-g` overrides auto-detection; `resolveGameName()` priority chain: flag > auto-detect > fallback
    - CLI: `detectGameName()` improved with 2 strategies: parent dir match + search anywhere in path
    - GUI: `gameName` parameter added to ScriptDecompile/ScriptCompile, passed as `-g` to lucksystem

### Version 3.1.3 — Patch 1

18. **Script decompile GameName auto-detection** — `cmd/scriptDecompile.go`, `cmd/scriptImport.go`
    - Added `detectGameName()`: extracts game name from OPCODE path (e.g. `data/LB_EN/OPCODE.txt` → `LB_EN`)
    - Ensures `LB_EN` operator is used instead of generic fallback → MESSAGE/SELECT/BATTLE decoded as text, not raw codepoints
    - Priority: plugin `.py` > auto-detect from `-O` path > `"Custom"` generic fallback

### Version 3.1.2 — Patch 1

17. **PAK Import/Export path separator fix (Windows)** — `pak/pak.go`
    - `path.Base()` → `filepath.Base()` in Import dir mode: fixed crash `strconv.Atoi` on full Windows paths
    - `path.Join()` → `filepath.Join()` in Export: fixed mixed `/`+`\` separators in list files causing CZ corruption on re-import
    - Fixed error variable leak (`err` scope) and file handle leak on skipped files

### Version 3.1.1 — Patch 1

16. **Undefined opcode warning verbosity reduction** — `game/operator/undefined_operate.go`, `cmd/scriptDecompile.go`, `cmd/scriptImport.go`
    - Replaced per-opcode `glog.V(5).Infoln()` (1,461 lines for LB_EN) with silent `opcodeTracker` accumulator
    - Single sorted summary block printed after `RunScript()` completes
    - Eliminates false "infinite loop" appearance on slow machines and in GUI

### Version 3.1 — Patch 1

15. **Little Busters EN script decompile fix** — `game/VM/vm.go`, `game/game.go`, `game/operator/generic.go` (new), `cmd/scriptDecompile.go`, `cmd/scriptImport.go`
    - `NewVM()`: nil pointer crash when no game-specific operator matched (e.g., `GameName: "Custom"`) — added nil guard + generic fallback operator
    - `game.go:load()`: SEEN8500/SEEN8501 (baseball mini-game data tables with `firstLen=0`) caused underflow panic in `restruct.Unpack` — added `isValidScript()` pre-check + `safeLoadScript()` panic recovery
    - `scriptDecompile.go` / `scriptImport.go`: auto-detection of `GameName` from OPCODE path (e.g., `data\LB_EN\OPCODE.txt` → `LB_EN`)
    - Generic operator handles common opcodes (IFN, IFY, GOTO, JUMP, FARCALL, GOSUB, EQU, ADD, RANDOM); unknown opcodes → `UNDEFINED` dump

### Version 3 — Patch 3

13. **CZ2 font import resize fix** — `czimage/cz2.go`, `font/font.go`
    - `Import()`: update `CzHeader` dimensions instead of silent `nil` return when image is resized (append/insert modes)
    - Added `SetDimensions()` method on `Cz2Image`
    - `Write()`: sync header before `Import()` call
14. **GUI improvements** — hidden CMD window, Stop button, free-text output fields for font edit, PAK Font Replace list mode

### Version 3 — Patch 2

- **Graphical Interface** — Wails + Svelte GUI (separate repository)

### Version 3 — Patch 1

8. **CZ1 32-bit Import/Export rewrite** — `czimage/cz1.go`
9. **CZ1 8-bit palette support** (Colorbits > 32 normalization) — `czimage/cz1.go`
10. **Non-CZ files graceful handling** — `czimage/cz.go`
11. **CZ0 logging visibility** — `czimage/cz0.go`

### Merged upstream (PR #35)

12. **CZ2 font decompressor crash fix** — `czimage/lzw.go` (boundary check in `decompressLZW2`)

### Version 2 (7 patches)

1. **Variable-length script import** — `script/script.go`
2. **CZ3 pipeline fixes** (magic byte, NRGBA, buffer aliasing) — `czimage/cz3.go`, `imagefix.go`
3. **LZW decompressor memory corruption** — `czimage/lzw.go`
4. **RawSize carry-over + UTF-8 length** — `czimage/util.go`
5. **CZ4 format support** (new) — `czimage/cz4.go`
6. **PAK block alignment padding** — `pak/pak.go`
7. **AIR.py module resolution** — `data/AIR.py`

---

## Features / Fonctionnalités

| Feature | Format | Status |
|---------|--------|--------|
| Script decompile/compile | SCRIPT.PAK | ✅ |
| CZ0 image export | CZ0 | ✅ |
| CZ1 image export/import (32-bit + 8-bit palette) | CZ1 | ✅ |
| CZ3 image export/import | CZ3 | ✅ |
| CZ4 image export/import | CZ4 | ✅ |
| Font extract/edit (append, insert, redraw) | FONT.PAK (CZ2) | ✅ |
| PAK extract/replace | *.PAK | ✅ |
| Graphical interface (GUI) | — | ✅ |

---

## Documentation

| Document | Description |
|----------|-------------|
| [CHANGELOG.md](CHANGELOG.md) | Full changelog — all versions (EN + FR) |
| [TECHNICAL.md](TECHNICAL.md) | Technical analysis — all patches |
| [Usage.md](Usage.md) | CLI command reference |

---

## Usage

```bash
# Decompile scripts
lucksystem script decompile -s SCRIPT.PAK -c UTF-8 -O data/AIR.txt -p data/AIR.py -o Export

# Import translated scripts
lucksystem script import -s SCRIPT.PAK -c UTF-8 -O data/AIR.txt -p data/AIR.py -i Export -o SCRIPT_FR.PAK

# Export CZ image to PNG
lucksystem image export -i image.cz3 -o image.png

# Extract FONT.PAK
lucksystem pak extract -i FONT.PAK -o list.txt --all ./fonts/

# Edit font with TTF (append French accents)
lucksystem font edit -s 明朝32 -S info32 -f Arial.ttf -o 明朝32_out -O info32_out -c accents_fr.txt -a
```

---

## Tested games / Jeux testés

- **AIR** (Steam) — French translation complete (scripts + CG + UI)
- **Summer Pockets** — RawSize fix confirmed
- **Kanon** — CZ2 font fix confirmed
- **Little Busters English** — CZ4 confirmed, script decompile confirmed (161 scripts, 102k+ MESSAGE lines, text properly decoded)

---

## Credits

- **[wetor](https://github.com/wetor)** — LuckSystem original
- **masagrator** — RawSize bug identification (CZ3 layers)
- **[G2-Games](https://github.com/G2-Games)** — CZ4 reference ([lbee-utils](https://github.com/G2-Games/lbee-utils))
- **Yoremi** — patches 1-20, AIR French translation, GUI
--------------------
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
