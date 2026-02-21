# Technical Analysis — LuckSystem Patches Yoremi

Technical document detailing the 7 patches applied to LuckSystem 2.3.2 for visual novel translation support.

---

## Patch 1 — Variable-length script import

### Modified file
`script/script.go` — lines 172-243

### Problem
The `VMRun()` function in import mode strictly checked that the imported parameter count (`code.Params`) matched the expected count (`expectedExportCount`). This check failed with translations of different length, because variable-length `StringParam` entries were not properly accounted for.

```go
// BEFORE: panic if lengths differ
if expectedExportCount != len(code.Params) {
    panic("导入参数数量不匹配...")
}
```

### Fix
- Removed the strict parameter count check (lines 175-194 of original)
- Replaced `for i := 0; i < len(paramList)` with `for i := 0; i < maxLen` where `maxLen = min(len(paramList), len(code.Params))`
- Added bounds checking (`if pi < len(code.Params)`) in the merge of `StringParam`, `JumpParam` and `[]uint16`

### Impact
Translations can now be longer or shorter than the original. Jump offsets are automatically recalculated by existing downstream code.

---

## Patch 2 — CZ3 pipeline fixes (PNG export/import)

### Modified files
`czimage/cz3.go`, `czimage/imagefix.go`

### Problem 1 — Magic byte overwritten (cz3.go, line 185)
The `Write()` function let the `CzHeader.Magic` field get corrupted from "CZ3" to "CZ0", making the file unreadable by the game engine.

```go
// FIX: force magic before write
cz.CzHeader.Magic = []byte{'C', 'Z', '3', 0}
```

### Problem 2 — NRGBA format not guaranteed (cz3.go, lines 84-99, 120-125)
The CZ3 format encodes pixels as BGRA 32-bit. Go's PNG library can decode to RGB, RGBA, NRGBA or paletted depending on the source file. Without forced conversion, a 24-bit RGB PNG is processed as 32-bit, shifting all pixels.

```go
// FIX: systematic conversion to NRGBA 32-bit
pic = ImageToNRGBA(cz.PngImage)
```

### Problem 3 — Buffer aliasing in DiffLine/LineDiff (imagefix.go)
The original code created a slice alias (`currLine = pic.Pix[i:...]`) instead of a copy. The delta operation (`currLine[x] -= preLine[x]`) modified `pic.Pix` directly, corrupting the source data for subsequent lines.

```go
// BEFORE (bugged): alias, modifies pic.Pix
currLine = pic.Pix[i : i+lineByteCount]

// AFTER (fixed): copy into separate buffer
copy(currLine, pic.Pix[i:i+lineByteCount])
```

Same issue in `LineDiff()`: `preLine = currLine` created an alias instead of a copy.

---

## Patch 3 — LZW decompressor memory corruption

### Modified file
`czimage/lzw.go` — functions `decompressLZW()` and `decompressLZW2()`

### Problem
The LZW dictionary added entries using `dictionary[dictionaryCount] = append(w, entry[0])`. In Go, `append()` can return the same underlying slice if capacity allows. When `w` was later reassigned (`w = entry`), old dictionary entries could point to modified data.

### Fix
Explicit allocation of a new slice before adding to dictionary:

```go
// BEFORE (bugged):
dictionary[dictionaryCount] = append(w, entry[0])

// AFTER (fixed):
newEntry := make([]byte, len(w)+1)
copy(newEntry, w)
newEntry[len(w)] = entry[0]
dictionary[dictionaryCount] = newEntry
```

---

## Patch 4 — Incorrect RawSize in CZ block table

### Modified file
`czimage/util.go` — functions `Compress()` and `Compress2()`

### Context
The CZ3 format stores pixels as LZW-compressed blocks. Each block declares a `CompressedSize` and `RawSize` in a header table. The game engine relies strictly on the declared `RawSize` to allocate buffers and position data during decompression.

### Problem 1 — LZW carry-over not compensated

The `compressLZW()` function operates in blocks of `size` maximum codes. When the limit is reached, it keeps a `lastElement` (the entry being built in the dictionary) that is carried over to the next block. The `count` value returned includes the bytes of this element:

```
Block N: reads 502 bytes, produces 500 codes, keeps 1 byte as carry-over
  → count = 502, lastElement = 1 byte
  → RawSize SHOULD be 501 (502 read - 1 carried over)
  → RawSize WAS 502 (bugged)
```

Effect: the first block declares a `RawSize` too large by 1, the last block too small by 1. The game engine reads 1 byte too many from the first block and 1 byte too few from the last, causing a cascading offset shift.

### Problem 2 — Go UTF-8 encoding

LuckSystem uses Go `string` values as LZW dictionary keys. The carry-over element is built using `element = string(c)` where `c` is a `byte` (0-255).

In Go, `string(byte(c))` performs a `byte → rune → UTF-8` conversion. For bytes 0-127, the result has length 1. For bytes 128-255, Go produces a **2-byte** UTF-8 string:

```go
string(byte(127)) // len = 1 (ASCII)
string(byte(128)) // len = 2 (UTF-8: 0xC2 0x80)
string(byte(255)) // len = 2 (UTF-8: 0xC3 0xBF)
```

An early fix attempt used `len(last)` to count carry-over data bytes. For a carry value of 200, `len(last) = 2` in Go while it represents only 1 data byte. This caused ±1 errors on blocks whose carry-over fell on a byte > 127.

### Final fix

The carry-over from `compressLZW()` is **always 0 or 1 data byte**, regardless of `len(last)` in Go:

```go
// BEFORE (bugged) — original version:
RawSize: uint32(count) // includes carry-over

// ATTEMPT 1 (partially bugged):
rawSize := prevCarryLen + count - len(last) // len(last) ≠ 1 for bytes > 127

// AFTER (fixed):
carry := 0
if len(last) > 0 {
    carry = 1  // always 1 DATA byte, regardless of len(last) in Go
}
rawSize := prevCarry + count - carry
```

### Verification
Round-trip test on original AIR CZ3 (`title1a`, 1280×720, 32-bit, 10 blocks): all 10 `RawSize` values produced by the fixed version match **exactly** those of the original file created by Visual Art's tools.

### Scope
This bug affects **all games** supported by LuckSystem that use multi-block CZ images (i.e., any image whose compressed data exceeds 0xFEFD LZW codes — the vast majority of CGs). Single-block CZ images (small UI elements) are not affected because there is no carry-over.

---

## Patch 5 — CZ4 image format support

### New/modified files
`czimage/cz4.go` (new), `czimage/imagefix.go`, `czimage/cz.go`

### CZ4 format specification

CZ4 is used in newer LUCA System games (Little Busters English Edition, LOOPERS, Harmonia, Kanon 2024). It differs from CZ3 in the pixel data layout:

| Aspect | CZ3 | CZ4 |
|--------|-----|-----|
| Pixel storage | Interleaved RGBA (4 bytes/pixel) | Separated channels: [RGB w×h×3] [Alpha w×h] |
| Delta encoding | Full RGBA lines together | RGB and Alpha independently |
| LZW compression | Same | Same |
| Block height | `(h+2)/3` | `(h+2)/3` |
| Header | CzHeader + Cz3Header (28 bytes) | Identical to CZ3 |

### File structure
```
[Magic "CZ4\x00"] [CzHeader 15 bytes] [Cz3Header 13 bytes]
[Block count (uint32)] [Block table: {CompressedSize, RawSize} × N]
[LZW compressed data]
```

### Decode algorithm (LineDiff4)
1. LZW decompress → raw data `[RGB section: w×h×3 bytes][Alpha section: w×h bytes]`
2. Delta decode RGB section: for each line, if `y % blockHeight != 0`, `curr[x] += prev[x]`
3. Delta decode Alpha section: same algorithm, independent
4. Interleave RGB + Alpha → NRGBA output

### Encode algorithm (DiffLine4)
1. Split NRGBA → `[RGB section][Alpha section]`
2. Delta encode RGB: for each line, if `y % blockHeight != 0`, `curr[x] -= prev[x]`; `prev[x] += curr[x]`
3. Delta encode Alpha: same algorithm
4. LZW compress concatenated `[RGB][Alpha]`

### Verification
Round-trip test on all 7 CZ4 files from AIR SYSCG.pak: 7/7 perfect (0 bytes differ). Full SYSCG.pak: 51/51 files exported and re-imported without errors (CZ3 + CZ4 combined).

### Reference
Based on [lbee-utils](https://github.com/G2-Games/lbee-utils) by G2-Games (Rust implementation).

---

## Patch 6 — PAK block alignment padding

### Modified file
`pak/pak.go` — function `Write()`

### Problem
When rebuilding a PAK file with replaced entries (larger files), the total file size was not aligned to `BlockSize`. Some game engines check this alignment when loading PAK archives, causing read errors.

### Fix
After writing all file data, compute the position of the last byte written and pad with zero bytes to the next `BlockSize` boundary:

```go
if lastFileEnd%p.BlockSize != 0 {
    paddingSize := p.BlockSize - (lastFileEnd % p.BlockSize)
    padding := make([]byte, paddingSize)
    _, err = file.WriteAt(padding, int64(lastFileEnd))
}
```

---

## Patch 7 — AIR.py module resolution fix

### Modified file
`data/AIR.py`

### Problem
The original `AIR.py` used `from base.air import *` to import jump-related functions (IFN, IFY, FARCALL, GOTO, GOSUB, JUMP, etc.) from `data/base/air.py`. This import failed when running the `script import` command because LuckSystem's embedded Python runtime does not set the working directory to `data/`, so the relative path `base/air` cannot be resolved.

Reproduced with the command documented in usage.md:
```
lucksystem script import -p "data/AIR.py" ...
→ FileNotFoundError: 'Failed to resolve "base/air"'
→ panic: runtime error: invalid memory address or nil pointer dereference
```

Using `data/base/air.py` directly as the `-p` argument avoids the import error, but then the text-parsing functions (MESSAGE, VARSTR_SET, DIALOG, LOG_BEGIN, SELECT) are missing, causing the import to fail later when encountering these opcodes.

### Fix
- Merged all functions from `base/air.py` directly into `AIR.py` (IFN, IFY, FARCALL, GOTO, ONGOTO, GOSUB, JUMP, JUMPPOINT, RETURN, FARRETURN)
- Added the `ONGOTO` function which was missing from both files
- Removed the `from base.air import *` dependency
- The file is now fully self-contained and works regardless of working directory

---

## Modified files (summary)

| File | Patch | Description |
|------|-------|-------------|
| `script/script.go` | 1 | Variable-length script import |
| `czimage/cz3.go` | 2 | Magic byte, NRGBA conversion, logging |
| `czimage/imagefix.go` | 2, 5 | Buffer aliasing fix + CZ4 delta encode/decode |
| `czimage/lzw.go` | 3 | LZW dictionary memory corruption |
| `czimage/util.go` | 4 | RawSize carry-over + UTF-8 length |
| `czimage/cz4.go` | 5 | New: CZ4 format support |
| `czimage/cz.go` | 5 | CZ4 dispatcher in LoadCzImage |
| `pak/pak.go` | 6 | PAK block alignment padding |
| `data/AIR.py` | 7 | Self-contained module, no base/ dependency |

## Known remaining issues

- **CZ1 8-bit palette**: `GetOutputInfo()` does not account for the 1024-byte color palette between the header and block table for CZ1 images with `Colorbits=8` or `Colorbits=4`. This causes a panic on palette-indexed CZ1 files (e.g., `system_icon_*` in PARTS.pak).
- **CZ1 32-bit import**: `cz1.go Import()` only exports the alpha channel (`data[i] = A`), discarding RGB. Re-importing a CZ1 32-bit RGBA image produces a blank/white result (e.g., `systemmenu` in PARTS.pak).
- **Non-CZ files in PAK**: Some files extracted from PAK archives are not CZ images (e.g., トーンカーブ_夕/夜 are 768-byte RGB tone curve LUTs). These should be detected and skipped gracefully instead of crashing with "Unknown CZ image type".
