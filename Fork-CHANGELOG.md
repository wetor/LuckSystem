# V3.20 — Script plugin auto-selection + Dialogue GUI LOG_BEGIN hardening + AIR empty string fix

13/06/2026

## Added: safer script decompile/import when a plugin is omitted

### Context

A CartagraHD tester reported that translated `LOG_BEGIN` lines appeared in the decompiled `.txt` scripts but did not seem to apply after repacking. The supplied case showed that the edited scripts already contained the English `LOG_BEGIN` text, including the first visible line:

```text
LOG_BEGIN ("The roar of water fills my ears.")
```

Round-trip testing confirmed that `LOG_BEGIN` imports correctly when the CartagraHD plugin is loaded. The issue was therefore not a confirmed engine import bug; the likely failure mode was repacking with the OPCODE file but without the matching Python plugin, causing the generic fallback parser to be used.

### Fix / hardening

**CLI**

- `script decompile` and `script import` now auto-select a sibling Python plugin when `-O` points to a standard repository layout and `-p` was left empty:
  - `data/GAME.txt` -> `data/GAME.py`
  - `data/GAME/OPCODE.txt` -> `data/GAME.py`
- This keeps plugin-backed opcodes such as `MESSAGE`, `LOG_BEGIN`, and `SELECT` on the proper game-specific import/export path even if the user forgets to browse the plugin manually.
- The CLI prints an explicit info line, for example:

```text
[INFO] Auto-selected plugin: data\CartagraHD.py
```

**GUI**

- Dialogue extract/import opcode detection is now stricter:
  - accepts `LOG_BEGIN` behind `labelN:` and `globalN:` prefixes;
  - avoids treating `MESSAGE_CLEAR`, `MESSAGE_WAIT`, and similar non-dialogue opcodes as `MESSAGE`;
  - keeps `MESSAGE`, `LOG_BEGIN`, and `SELECT` handling aligned between backend code and visible GUI help text.
- Added a focused GUI backend test covering extraction and import of normal, labelled, and global-labelled `LOG_BEGIN` lines.
- Frontend npm scripts now invoke Vite through `node ./node_modules/vite/bin/vite.js` instead of executing the `.bin/vite` shim directly. This avoids `sh: 1: vite: Permission denied` on Linux builds when `node_modules/.bin/vite` lost its executable bit.
- Updated CLI and GUI version labels to `v3.20`.

**AIR**

- AIR `MESSAGE` lines with an empty UTF-8 translated string now decompile correctly. Some AIR Steam scripts encode that field as `00 00 00` (zero length plus string terminator); the reader now consumes the terminator and the writer preserves it during import.
- Fixes the `seen203` panic: `runtime error: slice bounds out of range` during AIR script extraction.

### Testing

- `go test ./... -count=1` from `SourcesGUI-wails`: OK.
- `go test ./cmd ./script`: OK.
- `go test ./cmd ./game/operator ./script`: OK.
- `go test ./... -run '^$'`: OK.
- `npm run build` from `SourcesGUI-wails/frontend`: OK.
- Repacked the supplied CartagraHD case without passing `-p`; auto-selection picked `data\CartagraHD.py`.
- Redecoded the resulting PAK and confirmed `0000-op0_HD.txt:144` still contains:

```text
LOG_BEGIN ("The roar of water fills my ears.")
```

- Decompile of the supplied AIR Steam `BAK-SCRIPT.PAK` with `data\AIR.py`: OK, including `seen203`.
- Import of the exported AIR scripts followed by redecompile of the rebuilt PAK: OK.

---

# V3.1.9 — CartagraHD ONGOTO fix + multi-goto support + zero-length string dump fix

06/06/2026

## Fixed: CartagraHD choices broken after translation (ONGOTO offsets not recalculated)

### Problem

Translating CartagraHD scripts produced broken in-game choices: branch offsets were exported as raw numbers and remained frozen when a dialogue line before a choice grew in length. The root causes were two independent bugs:

1. `ONGOTO` was not defined in `base/cartagrahd.py` — the opcode fell through to `UNDEFINED()`, so its jump targets were dumped as raw `uint16` values instead of labelled `{goto label...}` references. Since the import machinery only recalculates offsets expressed as labels, any line size change after an ONGOTO left all downstream branches pointing to wrong positions.

2. The engine's label parser handled only a single `{goto ...}` token per line. ONGOTO carries N branch targets on the same line; all targets after the first were silently ignored.

A third minor bug was also found: `operator/util.go` emitted a spurious extra character in the string dump when it encountered a zero-length string entry, producing noise in the exported text.

### Fix (4 files — CLI)

**`data/base/cartagrahd.py`**
- Added `ONGOTO` handler: reads the branch count N, then reads N `uint16` offsets and emits them as `{goto label_NNNN}` references, matching the pattern already used by `IFN`/`IFY`.

**`script/model.go`**
- Extended the `JumpParam` / label model to store a slice of jump targets per line instead of a single target, enabling one script line to hold multiple `{goto ...}` tokens.

**`script/script.go`**
- Updated `Export()` and `Import()` to iterate over all jump targets on a line (not just the first) when building and consuming label references.
- `Import()`: recalculates every target offset independently, so all N branches of an ONGOTO are correctly repointed after a size change.

**`game/operator/util.go`**
- Fixed zero-length string edge case in the string dump helper: a zero-length entry no longer emits a spurious character before the closing delimiter.

### Testing

- `go test ./script ./game/operator`: OK.
- Round-trip CartagraHD original (no translation): all 277 internal PAK entries are byte-identical; PAK hash matches original.
- Regression test — line 46 extended: ONGOTO targets shift from `3730 / 8104` to `3758 / 8132`, matching the two new absolute positions exactly.
- `go test ./...`: remaining failures are pre-existing fixture-only failures (absent `FONT.PAK`, `SP.py`, LOOPERS `SCRIPT.PAK`) — unrelated to this patch.

### Note

Existing CartagraHD dumps that contain raw ONGOTO numbers (e.g. `65530, 12835, …`) must be re-extracted with the corrected plugin before reimport. Old dumps lack the `{goto label}` tokens required for offset recalculation.

### Games affected

CartagraHD — any script containing ONGOTO (choice branches).

---

# V3.1.8 — Dedicated Vietnamese font GUI patcher + Latin redraw test mode

01/06/2026

## Added: AIR / Planetarian SG Vietnamese font generation from the GUI

### Problem

The v3.1.7 font patcher fixed AIR's PAK/CZ2/font-info round-trip issues, but it was still possible for a tester to use an old standalone helper executable and generate broken PAKs. It also left one visual question open: mixed text could still combine original engine Latin glyphs with injected Vietnamese glyphs drawn from the selected TTF.

### Fix

**GUI**

- Added/updated the dedicated `VIET FONT -> AIR / SG Patch` workflow.
- The GUI now calls the corrected Vietnamese font patch code directly instead of requiring a separate `vietnamesefont.exe` / `vietfontpatch.exe`.
- Added an experimental checkbox: `Redraw Latin alphabet from TTF`.
  - Disabled: safe mode, inject only missing Vietnamese glyphs.
  - Enabled: redraw existing `A-Z/a-z` cells and already-present Vietnamese glyphs from the selected TTF, then inject only the missing Vietnamese glyphs into tail cells.
- Experimental outputs include `_LATIN` in the folder name, for example `Arial_en_GOTHIC1_LATIN_Y+2`, so safe and test builds cannot overwrite each other.
- Kept the recommended first test as English slot + `GOTHIC1` + `Y+2`.
- Updated GUI and CLI version labels to `v3.1.8`.

### Technical notes

The Latin test mode does not append duplicate ASCII letters to the end of the charset. The engine already has mappings for `A-Z/a-z`, so appended duplicates would likely be ignored. Instead, the GUI redraws those existing mapped cells in place using the selected TTF.

Already-present Vietnamese glyphs from the requested charset are also redrawn in place in experimental mode. This avoids mixing original accented glyphs with newly injected TTF glyphs inside the same Vietnamese sentence.

### Testing

- `go test ./... -run '^$'`
- `go test ./... -run '^$'` from `SourcesGUI-wails`
- `npm run build` from `SourcesGUI-wails/frontend`
- `wails build` from `SourcesGUI-wails`

---

# V3.1.7 — AIR Vietnamese font rebuild / compact PAK fix

22/05/2026

## Fixed: AIR no longer crashes after font round-trip or Vietnamese charset injection

### Problem

AIR Steam accepted edited `FONT__INFO.PAK` files on their own, but crashed on startup as soon as the large English/Japanese `FONT_GOTHIC1.PAK` was rewritten, even with a no-op round-trip. The same edits also produced visual menu glitches when several font sizes or families were touched.

When Vietnamese glyphs were injected by replacing the tail of the original charset, missing characters started to appear, but early builds had two rendering problems:

- AIR's font PAK preload path was sensitive to rewritten PAK layout.
- Newly drawn Vietnamese glyphs used TTF vertical metrics directly, producing negative Y offsets and visible glyphs floating too high.

### Root cause

- `pak.Write()` copied the whole source PAK first. When replacements were smaller than the original slots and `Rebuild` was false, the output kept internal holes. When `Rebuild` was true, the file could still keep a stale copied tail unless it was explicitly truncated.
- AIR font info tables use the legacy layout `CharNum=100` plus `CharNum2=<real count>`. Loading normalized the count, but writing did not preserve that layout.
- Partial font replacement recomputed atlas dimensions from character count, which could shrink/reshape an atlas that AIR expected to stay byte-layout compatible.
- CZ2 import recompressed the alpha image using LuckSystem's generic block splitting. AIR tolerated some rewritten CZ2s, but the large Japanese Gothic atlas path proved stricter.
- Several Vietnamese characters already existed in AIR's charset (`á`, `ó`, `â`, etc.). Replacing them with newly drawn glyphs regressed their original metrics.

### Fix

**CLI / core**

- `font/info.go`
  - Preserve the `CharNum=100 + CharNum2` info layout on write when the source used it.
- `font/font.go`
  - Preserve original atlas dimensions during partial replacement.
  - Copy the old atlas first and redraw only the replaced cells.
- `czimage/cz2.go`, `czimage/util.go`
  - Preserve original CZ2 raw block boundaries when the edited image has the same total raw size.
- `pak/pak.go`
  - Support compact rebuilt PAKs with recalculated offsets.
  - Truncate rebuilt output to the aligned real end so stale bytes from the copied source PAK cannot remain.

**CLI helper tools**

- `tools/fontdiag`
  - New diagnostic tool to round-trip one font family PAK without changing the charset.
  - Forces compact rebuilds so AIR startup tests isolate CZ2/font writer behavior from PAK holes.
- `tools/vietfontpatch`
  - New AIR font patch helper.
  - Injects only characters missing from the original slot and keeps already-present Vietnamese glyphs mapped to their original cells.
  - Supports `-slot all|en|zc`, `-family all|GOTHIC1|...`, and `-yoffset N`.
  - Normalizes injected glyph vertical metrics against original Latin/accented glyphs. AIR English slot testing selected `-yoffset 2` as the best visual match.

**GUI**

- Updated GUI title/about/version text to `v3.1.7`.
- Removed stale duplicate `frontend/src/dialogue.go`; the maintained dialogue extract/import implementation is in `app.go`. This prevents `go test ./...` from treating the frontend folder as a broken Go package.
- No GUI workflow change is required: the GUI remains a subprocess wrapper around the CLI. Existing Font Extract/Edit forms continue to use the patched core code once rebuilt with the new CLI.

**Maintenance**

- `game/runtime/global_goto.go`: fixed two `%s`/integer log format strings so `go test ./...` no longer fails Go vet on that package.

### Testing

AIR Steam font tests:

- Info-only replacement starts successfully.
- Previous no-op `FONT_GOTHIC1.PAK` round-trip crash is fixed with compact rebuild.
- Full Vietnamese charset injection starts without menu visual corruption.
- Missing Vietnamese characters render in the English slot.
- Visual comparison of `Y+1`, `Y+2`, `Y+3` confirmed `Y+2` as the best default for the tested TTF.
- `FONTZC_*` generation remains supported by the helper, but final AIR validation focused on the English slot only.

Build checks:

- `go test ./tools/fontdiag ./tools/vietfontpatch ./czimage`
- GUI source/version strings checked; frontend dependencies are intentionally not committed.

---

# V3.1.6 — `font edit` / `font extract` panic fix for AIR & planetarian CZ2 fonts

17/05/2026

## Bug fixed: `Cz2Image.decompress` panics with `index out of range` on round-trip and silently corrupts pixels on load

### Problem

Running `font edit` on the Gothic fonts of AIR and planetarian (and re-extracting the resulting CZ2) crashed with:

```
panic: runtime error: index out of range [2143615] with length 2143615
goroutine 1 [running]:
lucksystem/czimage.(*Cz2Image).decompress
        czimage/cz2.go:69 +0x24d
```

The decompressed alpha buffer was shorter than `Width × Height`, so the inner `for y { for x { pic.SetNRGBA(x, y, cz.ColorPanel[buf[i]]) }}` loop indexed past the end of the slice. Reproduces on every Gothic size from 22 onwards (12 and 16 happened to fit under the threshold described below).

The same root cause was silently corrupting the *load* path too: extracting a Gothic font produced a PNG whose bytes drifted slightly from the source CZ2 (typically a few dozen bytes per file, visually undetectable but byte-different — confirmed by md5sum on the pre- vs. post-patch extract of `ゴシック32`).

### Root cause — two independent bugs in the LZW codec, both in `czimage/lzw.go`

**Bug 1 — slice-aliasing corruption in `decompressLZW2`:**

```go
dataSize := len(data)
data = append(data, []byte{0, 0}...)
```

`Decompress2` slices the full compressed buffer into per-block sub-slices via `compressed[offsetTemp:offset]`. Those sub-slices share the underlying array with the parent, and their *capacity* extends to the end of the parent buffer. So `append(data, 0, 0)` doesn't allocate; it writes two zero bytes directly into the parent at `parent[offset]` and `parent[offset+1]` — i.e. into the **first two bytes of the next block's bitstream**.

Every block after the first was reading a corrupted 16-bit prefix: the LZW code that should have been emitted as the first symbol of the block (the encoder's carry-over byte from the previous block) was overwritten with `00 00`. On `Decompress2(original_file)` the corruption was small (one byte per block boundary, plus dictionary drift of a few bytes). On the `Compress2` → `Decompress2` round-trip used by `font edit`, the corruption compounded across all five blocks and the total decoded length fell short of `W×H`, triggering the panic.

**Bug 2 — silent 18-bit truncation in `compressLZW2`:**

```go
writeBit := func(code uint64) {
    if code > 0x7FFF {
        bitIO.WriteBit(1, 1)
        bitIO.WriteBit(code, 18)  // <-- writes low 18 bits only
    } else { ... }
}
// ... main loop ...
dictionary[entry] = dictionaryCount
dictionaryCount++  // <-- never capped
```

The wire format encodes long codes in 18 bits, so the maximum representable code is `0x3FFFF` (= 262143). The encoder's dictionary, however, was incremented unconditionally. On large blocks (default `0x87BDF` ≈ 555 KiB target compressed size, the value `Compress2` uses when no hint is given) `dictionaryCount` overflows past `0x40000` and assigns codes that no longer fit in 18 bits. When such a code is later emitted, `WriteBit(code, 18)` silently truncates to the low 18 bits, producing a completely different value. The decoder reads that truncated value, finds it in *its* dictionary pointing to a much shorter byte sequence, and the per-block decoded output ends up several bytes shorter than the encoder's `RawSize` claimed.

Per-block traces on `ゴシック32` after Bug 1 was fixed:

```
block[0]: maxEmit=261664  overflows=0   loss=0
block[1]: maxEmit=261978  overflows=0   loss=0
block[2]: maxEmit=262532  overflows=16  loss=179
block[3]: maxEmit=262149  overflows=1   loss=5
block[4]: (last, small)                 loss=0
```

The 184-byte shortfall (179 + 5) matched the gap between the decoded buffer and `W×H` exactly.

### Fix (1 file — CLI)

**`czimage/lzw.go`**

- `decompressLZW2()`: copy the compressed sub-slice into a freshly allocated buffer with explicit `len = dataSize + 2` *before* feeding it to `NewBitIO`. The append never aliases the parent, so the next block's first two bytes stay intact.

  ```go
  padded := make([]byte, dataSize+2, dataSize+2)
  copy(padded, data)
  bitIO := NewBitIO(padded)
  ```

- `compressLZW2()`: freeze the dictionary once it reaches `0x40000` entries — codes already assigned keep being used for lookups, no new entries are added. This is the standard LZW behaviour when no clear-code is available in the wire format. Block sizes don't change in practice (the per-block byte target is unaffected); the only effect is that the encoder stops emitting un-decodable codes on very large blocks.

  ```go
  if dictionaryCount < 0x40000 {
      dictionary[entry] = dictionaryCount
      dictionaryCount++
  }
  ```

### Why this matters

`font edit -r` (redraw current font in a different TTF) and `font edit -c` (append/replace charset) were both unusable on the larger Gothic sizes for Visual Art's/Key Luck Engine games — i.e. exactly the sizes used for in-game dialogue. Any French/accented build that touched the font produced a CZ2 that silently re-loaded with a too-short buffer and crashed at extract or in-game.

Beyond the panic, the silent on-load pixel drift in Bug 1 also meant that any extracted PNG used as a reference (e.g. for diffing two patches, or for hand-pixel-art touch-ups) had a few wrong bytes near block boundaries. After this patch, `extract → compare` is byte-stable.

### Backward compatibility

- No exported function signatures changed.
- CZ2 files produced by the patched encoder remain decodable by the patched decoder; on the round-trip test (alpha bytes → `Compress2` → `Decompress2`) all 13 Gothic sizes (12 → 38) now decode byte-for-byte identical to the input.
- CZ2 files produced by the *unpatched* encoder are decoded correctly by the patched decoder, since the dictionary cap only changes what the encoder emits and the slice-aliasing fix is purely a defensive copy on read.
- The same overflow window exists in `compressLZW`/`Compress` (CZ1/CZ3): `dictionaryCount` is a `uint16` that wraps at 65536, and the default block size of `0xFEFD` (= 65277) keeps codes just under the wraparound for normal-sized images. It has not been observed to fire in practice on Key fonts but should be revisited if a CZ1/CZ3 round-trip ever produces short output.

### Testing

On AIR's Gothic font (`FONT_GOTHIC1.PAK`, 13 sizes from 12 to 38):

- `pak extract` of both `FONT__INFO.PAK` and `FONT_GOTHIC1.PAK`: OK.
- `font extract` on each of the 13 sizes: OK before *and* after the patch (no panic on load), but post-patch PNG md5 differs from pre-patch — the patched extract is the byte-correct one.
- `font edit -r` (redraw with DejaVu Sans) on sizes 12 / 22 / 32 / 38 followed by `font extract` on the result: OK on all four; pre-patch panicked from size 22 upward.
- Round-trip `Compress2(data) → Decompress2 → data'` on all 13 sizes: byte-exact on every size.
- Full `Cz2Image.Import(PNG) → Cz2Image.Write → LoadCzImage → GetImage` chain (the exact path `font edit` uses) on all 13 sizes: dimensions and decode succeed on every size.

---

# V3.1.5 — Improved error reporting for script import and silent raw-byte log removal

23/04/2026

## Bug fixed: `script import` crashes with cryptic panic on stray newlines in translated scripts

### Problem
When a translated script file contained an accidental newline inside a dialogue line (e.g. a line break before a closing `❞`), all subsequent opcodes were shifted by one line. The import continued silently until it hit a type mismatch deep in `SetOperateParams`, producing an unreadable Go panic:

```
panic: interface conversion: interface {} is *script.JumpParam, not string

goroutine 1 [running]:
lucksystem/script.(*Script).SetOperateParams(0xc00022e480, 0x691, 0x2, ...)
        script/script.go:218 +0x13c5
```

The panic gave no indication of which script file, which line, or which opcode was at fault. With 30+ script files and thousands of lines each, finding the stray newline required manual binary search across all translated files.

### Root cause
A single extra `\n` inside a MESSAGE text (e.g. line break between `autres` and `❞`) created one extra line in the `.txt` file. Since LuckSystem reads exactly N lines (one per opcode in the binary script), every line after the break was mapped to the wrong opcode. The mismatch eventually reached a SELECT/IFN opcode where a `*JumpParam` was cast as `string`, triggering the panic.

### Fix (3 files — CLI)

**`script/script.go`**
- `Import()`: error messages now include script name and line number; new end-of-file check detects extra lines and reports:
  `[seen110] file has 1 extra line(s) beyond expected 3206 (check for stray newlines in translated text)`
- `SetOperateParams()`: all `.(string)` type assertions replaced with safe `ok`-checked assertions; on mismatch, returns a clear error:
  `[seen110] line 1137 (MESSAGE): parameter 2 type mismatch: expected string, got *script.JumpParam (likely a stray newline shifted all lines)`
- `CodeParamsToBytes()`: raw-byte dump log moved from `V(4)` to `V(8)` — eliminates the massive console output that made real errors invisible (this log fired on every translated line since French text has different byte length than English/Japanese)

**`game/VM/vm.go`**
- `Run()`: added `defer/recover` that catches any panic during script processing and reformats it with context:
  `[seen110] line 1137 (MESSAGE): <original panic message>`

**`game/operator/opcode.go`**
- `SetOperateParams()`: error return from `script.SetOperateParams()` is now propagated (was silently discarded with `_ =`)

### Why this matters

A single misplaced newline in any of 30+ translated script files would crash the entire import with no indication of which file or line was responsible. The translator had to resort to binary search across thousands of lines to find the problem. With this fix:

1. **Extra-line detection** catches the most common cause (stray newlines) before the crash even happens
2. **Safe type assertions** turn the cryptic Go panic into a readable error with file + line + opcode
3. **VM-level recover** ensures that even unexpected panics include script context
4. **Silent log cleanup** removes the raw-byte dump that flooded the console and hid real errors

### Backward compatibility
- All function signatures unchanged
- No new external dependency
- No change to exported script format or PAK output
- Log behavior: `V(4)` users will see less noise; `V(8)` restores full debug output if needed

### Testing
Confirmed on Kanon Steam: import of 30+ script files completes cleanly; intentionally broken file (stray newline in seen110.txt) produces clear error message instead of panic.

---

# V3.1.4 — Patch 1: Plugin import resolution and nil-module crash fix

13/04/2026

## Bug fixed: `script decompile` panics on plugins using `from base.xxx import *` (Kanon, AIR, HARMONIA, LOOPERS, LUNARiA, PlanetarianSG, CartagraHD)

### Problem
Decompiling Kanon Steam scripts crashed with a Go panic:

```
> lucksystem.exe script decompile -s SCRIPT.PAK -c UTF-8 -o TRAD \
    -O data/KANON.txt -p data/KANON.py -g KANON
[INFO] Using game: KANON (from --game flag)
Traceback (most recent call last):
  File "data/KANON.py", line 2, in <module>
FileNotFoundError: 'Failed to resolve "base/kanon"'
panic: runtime error: invalid memory address or nil pointer dereference
[signal 0xc0000005 code=0x0 addr=0x8 pc=0x...]
goroutine 1 [running]:
lucksystem/game/operator.(*Plugin).Init(...)
        game/operator/plugin.go:40 +0x82
```

The error affected every plugin that uses package-style imports (`from base.kanon import *`) — i.e. all games except AIR (which had been manually inlined in patch 7) and SP/LB_EN (which use Go operators, no Python plugin).

### Root cause — two cumulative bugs in `game/operator/plugin.go`

**Bug 1 — gpython sys.path misconfigured.** `NewPlugin()` initialised the gpython context with `SysPaths: []string{"."}` and `CurDir: "/"`. `from base.kanon import *` caused gpython to look for `base/kanon.py` relative to the *process working directory* (and `/` is not a valid `CurDir` on Windows anyway). The import only succeeded when lucksystem was launched from the directory containing the plugin's `base/` folder — which was never the case in the GUI or in normal CLI usage.

**Bug 2 — no nil-check after import failure.** When `py.RunFile()` failed, the error was logged via `py.TracebackDump()` but `p.module` remained `nil`. The function still returned a non-nil `*Plugin`. The next time `Init()` or `UNDEFINED()` accessed `g.module.Globals[...]`, the runtime panicked with a nil pointer dereference, hiding the real (Python) error behind a Go stack trace.

### Fix (1 file — CLI)

**`game/operator/plugin.go`**
- `NewPlugin()`:
  - Resolve plugin file to an absolute path via `filepath.Abs()` (graceful fallback to the original path if resolution fails)
  - Add the plugin's directory to `SysPaths` *before* `"."` → package imports like `from base.xxx import *` now resolve against the plugin tree regardless of cwd or OS
  - Use the plugin directory as `CurDir` (replaces the hardcoded `"/"`)
  - On load failure, print a readable `[ERROR] Failed to load plugin "<path>": <err>` line in addition to the Python traceback
- `Init()`: nil-guard on `g.module` — skip Python `Init()` call instead of panicking
- `UNDEFINED()`: nil-guard on `g.module` — fall through to default "advance PC" behaviour instead of panicking

### Why this matters

Without this fix, **only AIR (Steam) and the LB_EN/SP Go operators worked**. Every other game in the data/ tree (Kanon, HARMONIA, LOOPERS, LUNARiA, PlanetarianSG, CartagraHD) crashed at the first decompile attempt — silently for the user, since the panic came after the `[INFO] Using game: ...` line and looked like a tool bug rather than a plugin path issue.

### Backward compatibility
- `SysPaths` keeps `"."` as a fallback → no regression for setups that worked before
- `NewPlugin()` / `Init()` / `UNDEFINED()` signatures unchanged
- No new external dependency (`path/filepath` is stdlib)
- The patch is upstream-ready (no fork-specific markers); intended for merge into the wetor repo

### Games affected
Kanon, HARMONIA, LOOPERS, LUNARiA, PlanetarianSG, CartagraHD — and any future plugin that uses `base/` shared modules.

### Testing
Confirmed working on Kanon Steam (`KANON.py` + `base/kanon.py`):
```
lucksystem.exe script decompile -s SCRIPT.PAK -c UTF-8 \
    -o TRAD -O data/KANON.txt -p data/KANON.py -g KANON
```
Decompilation completes; `from base.kanon import *` resolves correctly; no panic.

---

# LuckSystem — Yoremi Fork — CHANGELOG

---

# V3.1.3 — Patch 3: GUI — Auto-detect game presets from data/ folder

28/02/2026

## Improvement: automatic game preset scanning replaces manual file browsing

### Problem
The GUI's Game dropdown (added in Patch 2) only listed LB_EN and SP. Users working with other games (AIR, KANON, HARMONIA, LOOPERS, LUNARiA, PlanetarianSG, CartagraHD) had to manually browse for their OPCODE and plugin files every time.

### Fix (4 files — GUI only)

**`app.go`**
- Added `GamePreset` struct and `ScanGameData()` function
- Scans the `data/` folder next to the lucksystem executable
- Detects all `.txt` files (OPCODE definitions) recursively, excluding `data/base/`
- Matches each game with its `.py` plugin if present at `data/GAME.py`
- Returns sorted list of presets with absolute paths

**`frontend/src/App.svelte`**
- Dynamic "Game preset" dropdown in Decompile and Compile forms, populated from `ScanGameData()`
- Selecting a preset auto-fills Opcode file, Plugin file, and Game flag
- Manual browse buttons still available (reset preset to "— Manual —")
- Presets rescanned when lucksystem path is changed via "Locate"

**`frontend/wailsjs/go/main/App.js`** + **`App.d.ts`** — added `ScanGameData()` binding

### Detected presets (from standard data/ folder)
AIR (plugin), CartagraHD (plugin), HARMONIA (plugin), KANON (plugin), LB_EN, LOOPERS (plugin), LUNARiA (plugin), PlanetarianSG (plugin), SP (plugin)

### Games affected
All games — improves workflow for every supported game.

---

# V3.1.3 — Patch 2: `--game` / `-g` flag for forced game type (CLI + GUI)

28/02/2026

## Fix: explicit game type override for cross-platform reliability

### Problem
Under Linux, the auto-detection from Patch 1 could fail if the OPCODE file was placed in an arbitrary directory (e.g. `~/Bureau/OPCODE.txt` instead of `data/LB_EN/OPCODE.txt`). The parent directory name didn't match "LB_EN" → fell back to "Custom" → generic operator → MESSAGE as raw codepoints.

### Fix — CLI (3 files)

**`cmd/script.go`**
- Added `ScriptGameName` variable and `--game`/`-g` persistent flag on the `script` command

**`cmd/scriptDecompile.go`**
- New `resolveGameName()`: priority chain: `--game` flag → auto-detect from OPCODE path → "Custom" fallback
- Improved `detectGameName()` with 2 strategies: (1) parent directory match (original), (2) search anywhere in path (new, catches `/project_LB_EN_scripts/opcodes/OPCODE.txt`)

**`cmd/scriptImport.go`**
- Uses shared `resolveGameName()`, removed duplicate detection logic and unused `fmt` import

### Fix — GUI (4 files)

**`app.go`**
- Added `gameName` parameter to `ScriptDecompile()` and `ScriptCompile()` signatures
- Passes `-g gameName` to lucksystem when non-empty

**`frontend/src/App.svelte`**
- Added `gameName` variable and Game dropdown (Auto-detect / LB_EN / SP) in Decompile and Compile forms
- Passed as 6th/7th argument to Go backend

**`frontend/wailsjs/go/main/App.js`** + **`App.d.ts`** — updated bindings for new parameter

### Priority chain
```
--game flag (-g LB_EN)  →  highest priority (explicit override)
Auto-detect from -O     →  Strategy 1: parent dir, Strategy 2: anywhere in path
Fallback                →  "Custom" → NewGeneric()
```

### Usage
```bash
# Method 1: Explicit flag (most reliable, especially on Linux)
lucksystem script decompile -s SCRIPT.PAK -c UTF-8 -o output -O OPCODE.txt -g LB_EN

# Method 2: Auto-detect from path (works if LB_EN appears anywhere in path)
lucksystem script decompile -s SCRIPT.PAK -c UTF-8 -o output -O /path/LB_EN/OPCODE.txt
```

### Games affected
LB_EN and SP — any game using a Go operator without a Python plugin.

---

# V3.1.3 — Patch 1: Script decompile GameName auto-detection fix

27/02/2026

## Bug fixed: MESSAGE/SELECT/BATTLE opcodes exported as raw codepoints instead of text

### Problem
Decompiling LB_EN scripts with `lucksystem script decompile -s SCRIPT.PAK -O data/LB_EN/OPCODE.txt` produced MESSAGE lines with raw Unicode codepoints instead of readable text:

```
MESSAGE (0, 12502, 12523, 12523, 12523, 12523, 8230, 12288, ...)
```

Expected output:
```
MESSAGE (0, "ブルルルル…　ブルルルル…", "Burururururu...  Burururururu...", 0x5)
```

All 161 scripts were affected — no dialogue text was visible, only numeric sequences. The same issue affected `script import` (round-trip would fail).

### Root cause — `GameName: "Custom"` hardcoded in scriptDecompile.go / scriptImport.go

Both `scriptDecompile.go` and `scriptImport.go` passed `GameName: "Custom"` to `game.NewGame()`, regardless of the OPCODE path provided. The dispatch chain:

```
scriptDecompile.go:  GameName: "Custom"          ← always hardcoded
        ↓
vm.go NewVM():  switch "Custom" → no match (not "LB_EN" nor "SP")
        ↓
vm.Operate = nil → fallback NewGeneric()         ← patch 15 safety net
        ↓
Generic has no MESSAGE() method → dispatch to UNDEFINED()
        ↓
UNDEFINED() calls AllToUint16() → dumps codepoints as numbers
```

The `LB_EN` operator (`operator/LB_EN.go`) already fully implements MESSAGE, SELECT, BATTLE, TASK, SAYAVOICETEXT, and VARSTR_SET with proper `DecodeString()` calls — it was simply never instantiated because the GameName never matched `"LB_EN"` in the switch.

### Note on patch 15 (v3.1)
Patch 15 documented "auto-detection of GameName from OPCODE path" but the implementation was incomplete — the `detectGameName()` function was not present in the delivered files. The generic fallback added in patch 15 prevented the nil pointer crash but did not solve the text decoding issue. This patch completes the auto-detection.

### Fix (2 files)

**`cmd/scriptDecompile.go`**
- Added `detectGameName(opcodePath string) string` function: extracts parent directory from OPCODE path using `filepath.Dir()` + `filepath.Base()`, compares case-insensitive against known games (`LB_EN`, `SP`)
- Auto-detection only runs when no plugin file (`-p`) is provided (plugins take priority)
- Prints `[INFO] Auto-detected game: LB_EN (from OPCODE path)` when a match is found
- Replaced `GameName: "Custom"` with `GameName: gameName`

**`cmd/scriptImport.go`**
- Same auto-detection logic using `detectGameName()` (defined in `scriptDecompile.go`, same package `cmd`)
- Replaced `GameName: "Custom"` with `GameName: gameName`

### Priority chain
```
Plugin (-p file.py)  →  highest priority (always used if provided)
Auto-detect from -O  →  "data/LB_EN/OPCODE.txt" → "LB_EN"
Fallback             →  "Custom" → NewGeneric()
```

### Games affected
LB_EN and SP — any game that uses a Go operator (not a Python plugin) and has its OPCODE file in a subdirectory matching the game name.

---

# V3.1.2 — Patch 1: PAK Import/Export path separator fix (Windows)

26/02/2026

## Bug fixed: `pak replace` crash in directory mode + CZ corruption via mixed path separators

### Problem 1 — Crash in directory mode
`lucksystem pak replace -s OTHCG.PAK -i <folder> -o out.PAK` crashed with `strconv.Atoi: parsing "C:\Users\...\msg_01k_en": invalid syntax`. All files in the folder were skipped with "Skip File", then the last `strconv.Atoi` error propagated as a fatal crash.

### Problem 2 — CZ corruption in list mode (old list files)
Re-injecting CZ files into a PAK using a list file generated by an older version produced corrupted (non-replaced) CZ images in-game. The replaced files themselves were fine — only neighboring untouched files showed corruption.

### Root cause — `path.Base()` / `path.Join()` vs `filepath` (Windows)

The `pak.go` file used the `path` package (POSIX-only, `/` separator) instead of `path/filepath` (OS-native separator) in three locations:

**Import() line 532** — `path.Base(file)`: On Windows, `filepath.Walk` returns paths with `\` separators. `path.Base("C:\...\msg_01k_en")` treats `\` as regular characters and returns the **entire path** as the filename. Result: `CheckName` always fails → `strconv.Atoi(entire_path)` → crash.

**Export() lines 412/415** — `path.Join(dir, name)`: Generated list files with mixed separators (`C:\Users\...\OTHCG_extracted/aug_01`). When these list files were later used for `pak replace -l`, the mixed paths caused subtle matching failures during PAK rebuild, corrupting offset calculations for non-replaced files.

### Additional bugs fixed

**Error variable leak** — `id, err = strconv.Atoi(name)` wrote to the outer-scope `err`. After a `continue` on the last file, `return err` at function end returned the Atoi error instead of `nil`. Fix: local `parseErr` via `:=`.

**File handle leak** — `fs, _ := os.Open(file)` opened files before validation. When skipped via `continue`, the file descriptor was never closed. Fix: `fs.Close()` before every `continue` + error handling on `os.Open`.

### Fix (1 file)

**`pak/pak.go`**
- `Import()` mode `"dir"`: `path.Base()` → `filepath.Base()`; local `parseErr`/`openErr` variables; `fs.Close()` on all skip paths
- `Export()` mode `"all"`: `path.Join()` → `filepath.Join()` (lines 412, 415)
- Removed unused `"path"` import

### Games affected
All games — any PAK replace operation on Windows was affected.

---

# V3.1.1 — Patch 1: Undefined opcode warning verbosity reduction

25/02/2026

## Improvement: silent accumulation of undefined opcode warnings

### Problem
During `script decompile` on Little Busters EN, the 1,461 undefined opcode warnings (`Operation不存在 HAIKEI_SET`, etc.) were printed one-by-one via `glog.V(5).Infoln()`. On slower machines or when using the GUI, this created an apparent infinite loop — the console scrolled warnings for over 2 minutes, making it look like the tool was stuck. The decompilation itself only takes ~5 seconds.

### Fix (3 files)

**`game/operator/undefined_operate.go`** — Replaced the per-opcode `glog.V(5).Infoln()` call with a thread-safe `opcodeTracker` that silently accumulates counts in a `map[string]int`. Exposed `PrintUndefinedOpcodeSummary()` which prints a single sorted summary block after processing completes, then resets the tracker.

**`cmd/scriptDecompile.go`** — Added `operator.PrintUndefinedOpcodeSummary()` call after `g.RunScript()`.

**`cmd/scriptImport.go`** — Same summary call after `g.RunScript()`.

### Result
Instead of 1,461 individual warning lines:
```
[INFO] 1461 undefined opcodes skipped (15 unique types):
  HAIKEI_SET            x312
  WAIT                  x245
  DRAW                  x198
  ...
These are non-text opcodes (visual/audio/system) and can be safely ignored for translation work.
```

### Games affected
All games using the generic operator or any game with undefined opcodes (Little Busters EN, Kanon, Harmonia, LOOPERS, LUNARiA, Planetarian).

---

# V3.1 — Patch 1: Little Busters EN script decompile fix

24/02/2026

## Bug fixed: nil pointer crash + invalid script entries in `script decompile`

### Problem
`lucksystem script decompile -s SCRIPT.PAK -c UTF-8 -o WORK -O data\LB_EN\OPCODE.txt` crashed immediately with a nil pointer dereference. Two independent bugs prevented Little Busters EN (and any game without a `-p` plugin flag) from being decompiled.

### Root cause 1 — `NewVM()` nil pointer on unknown GameName
`scriptDecompile.go` hardcoded `GameName: "Custom"` when no `-p` plugin was specified. In `vm.go:NewVM()`, the switch only matched `"LB_EN"` and `"SP"` — any other value left `vm.Operate = nil`. Line 48 then called `vm.Operate.Init(vm.Runtime)` → panic nil pointer dereference.

### Root cause 2 — SEEN8500/SEEN8501 data tables parsed as scripts
SCRIPT.PAK contains 169 entries but only 167 are actual scripts. SEEN8500 and SEEN8501 are baseball mini-game data tables with `firstLen=0` (first 2 bytes of entry data). When `script.LoadScript()` called `restruct.Unpack` with `Len=0`, it computed `size = Len - 4` → unsigned underflow → crash.

### Fix (5 files)

**`game/operator/generic.go` (new)** — Generic fallback operator that handles common opcodes (IFN, IFY, GOTO, JUMP, FARCALL, GOSUB, EQU, EQUN, ADD, RANDOM) via embedded `LucaOperateDefault` + `LucaOperateExpr`. Unknown opcodes are dumped as `UNDEFINED` with uint16 params.

**`game/VM/vm.go`** — Nil guard after the GameName switch: if `vm.Operate` is still nil, instantiate `operator.NewGeneric()` with a warning. This prevents the nil pointer crash for any unsupported game.

**`cmd/scriptDecompile.go`** — Auto-detection of GameName from the OPCODE path: `data\LB_EN\OPCODE.txt` → parent dir `LB_EN` → uses the LB_EN-specific operator (MESSAGE, SELECT, BATTLE, TASK handlers). Falls back to `"Custom"` if the directory name doesn't match any known game.

**`cmd/scriptImport.go`** — Same auto-detection logic as scriptDecompile.go.

**`game/game.go`** — Added `isValidScript()` pre-check (rejects entries with `firstLen < 4`) and `safeLoadScript()` panic recovery wrapper. SEEN8500/SEEN8501 are now skipped with a warning instead of crashing.

### Result
- 161 scripts decompiled successfully (SEEN8500/8501 skipped with warning)
- 102,795 MESSAGE lines extracted (bilingual JP/EN format)
- 1,461 opcode warnings for unhandled visual/audio opcodes (HAIKEI_SET, INIT, DRAW, WAIT, BGM, SE…) — expected, does not affect text extraction

### Games affected
Any game without a dedicated Python plugin file

---

# V3 — Patch 3: CZ2 Font Import Fix + GUI improvements

22/02/2026

## Bug fixed: CZ2 font reimport crash (`font edit`)

### Problem
`lucksystem font edit` crashed with `invalid argument` when using **append** or **insert** modes. Root cause: `ReplaceChars()` increases the image height when adding characters, but `CzHeader.Width/Heigth` were never updated. `Cz2Image.Import()` then:
1. Cropped the new (taller) image to the old dimensions via `FillImage()`
2. Hit the size mismatch branch and returned `nil` — silently, because `err` was declared but never assigned
3. Left `Raw` and `OutputInfo` empty → `WriteStruct` called `restruct.Pack` on an empty struct → `invalid argument`

### Fix
- `czimage/cz2.go` — `Import()`: when dimensions differ (valid case after append/insert), update `CzHeader.Width/Heigth` to match the new image instead of silently returning. Also replaced the direct `*image.NRGBA` type assertion with a safe conversion for any PNG format.
- `czimage/cz2.go` — added `SetDimensions(w, h uint16)` method on `Cz2Image`
- `font/font.go` — `Write()`: calls `SetDimensions()` before `CzImage.Import()` to sync the header with the image produced by `ReplaceChars()`

### GUI improvements (v3 GUI)
- **Hidden CMD window**: subprocess calls no longer flash a console popup on Windows (`SysProcAttr{HideWindow: true}` — platform-specific build file)
- **Stop button**: a **■ Stop** button appears in the console header during any running operation; cancels the subprocess via `context.WithCancel`
- **Font Edit output fields**: Output CZ and Output info are now free-text (no extension enforcement from Windows SaveFileDialog)
- **PAK Font Replace**: added list file mode (same as CG Replace) to avoid the `strconv.Atoi` crash on directory mode

---

# V3 — Patch 2: LuckSystem GUI

02/21/2026

## New: Graphical Interface (separate repository)

A standalone GUI for LuckSystem built with **Wails** (Go + Svelte). The GUI calls `lucksystem.exe` via subprocess following [wetor's architectural recommendation](https://github.com/wetor/LuckSystem) — no LuckSystem source code is embedded.

### Features
- **8 operations**: Script Decompile/Compile, PAK Extract/Replace, Font Extract/Edit, Image Export/Import
- **Batch mode** for Image Export/Import (entire folders)
- Real-time console output with color-coded logs
- Auto-detection of `lucksystem.exe` (same directory, CWD, or PATH)
- Custom application icon
- SiglusTools-inspired layout

### Architecture
```
LuckSystemGUI.exe  ←→  lucksystem.exe (subprocess)
   (Wails/Go)              (CLI tool)
```

### Build
Requires: Go 1.23+, Node.js, Wails CLI
```bash
cd frontend && npm install && cd ..
wails build
```

---

# V3 — Patch 1: CZ1 32-bit Import/Export + CZ0 logging

02/20/2026

## Modified files
- `czimage/cz1.go` — Import/Export/Write rewrite
- `czimage/cz.go` — graceful handling of non-CZ files
- `czimage/cz0.go` — added V(0) logging in decompress()

## Bugs fixed

### 1. Missing extended header in Write()
The original `Write()` only wrote the 15 bytes of the `CzHeader` struct, ignoring the 13 bytes of extended header (offsets, crop, bounds). The output file had the block table at offset 15 instead of 28 → crash on reload.

**Fix**: Save raw bytes 15→HeaderLength into `ExtendedHeader` at `Load()`, rewrite them in `Write()`.

### 2. Import() only handled alpha
The 32-bit Import only compressed the A channel (`data[i] = pic.A`), discarding RGB. Result: white/transparent screen in-game.

**Fix**: Multi-mode Import based on Colorbits (4, 8, 24, 32). The 32-bit mode does a direct copy of `pic.Pix` (RGBA).

### 3. Colorbits > 32 (8-bit palette)
CZ1 palette files use Colorbits=248 (0xF8), a proprietary Visual Art's marker. LuckSystem didn't recognize it → palette ignored → `GetOutputInfo()` read the palette as block table → crash (`slice bounds out of range`).

**Fix**: Normalization `if Colorbits > 32 → Colorbits = 8` (same approach as lbee-utils).

### 4. Non-CZ files in PAK
Files without the "CZ" magic (e.g., トーンカーブ_夕/夜, 768-byte RGB tone curve LUTs) caused a `glog.Fatalln("Unknown Cz image type")`.

**Fix**: Check magic before unpacking, return `nil` with warning instead of crash.

### 5. BGRA palette in Write()
The palette is read in BGRA (file) and stored as NRGBA (Go). The old Write via `restruct` serialized as RGBA → inverted R↔B colors.

**Fix**: Manual write of each palette entry as [B,G,R,A].

### 6. CZ0 invisible in extraction logs
CZ0 files only had `V(6)` logging (deep debug), while CZ4 logs at `V(0)` (always visible). When extracting a PAK containing a mix of CZ0/CZ4, the last visible lines before a CZ0 came from the previous CZ4, giving the impression that CZ0 files were processed as CZ4.

**Fix**: Added `glog.V(0).Infof("Decompress CZ0: %dx%d, Colorbits=%d")` in `cz0.go:decompress()` (line 78).

## CZ1 format confirmed
- 32-bit: pixels stored as **RGBA** (not BGRA like CZ3)
- 8-bit palette: entries stored as **BGRA**, data = 1 byte/pixel (index)
- Extended header: 13 bytes mandatory (same structure as Cz3Header)

## Status
- ✅ CZ1 32-bit: round-trip OK, tested in-game (systemmenu FR)
- ✅ CZ1 8-bit palette: OK (system_icon, NUM files)
- ✅ Non-CZ files: warning instead of crash
- ✅ CZ0: correctly identified in extraction logs

### Merged upstream (PR #35)
- **CZ2 font decompressor crash fix** — `czimage/lzw.go` (boundary check in `decompressLZW2`)

---

# Version 2 (7 patches)

02/18/2026

## Patch 1 — Variable-length script import
**File:** `script/script.go`

Importing translated scripts crashed with a panic when the translation had a different length than the original. The code strictly checked `len(paramList) == len(code.Params)`, blocking any longer or shorter translation.

- Removed strict parameter count check
- Added bounds checking in the conversion loop and parameter merge
- Jump offsets (GOTO, IFN, IFY…) are automatically recalculated

## Patch 2 — CZ3 pipeline fixes (PNG export/import)
**Files:** `czimage/cz3.go`, `czimage/imagefix.go`

CZ3 export and import silently corrupted pixel data.

- **Magic byte**: `Write()` overwrote the magic from "CZ3" to "CZ0", making the file unreadable by the game
- **NRGBA format**: Automatic conversion of any PNG format to NRGBA 32-bit before encoding
- **Buffer aliasing**: `DiffLine()` and `LineDiff()` shared slices instead of copying, causing delta data corruption

## Patch 3 — LZW decompressor memory corruption
**File:** `czimage/lzw.go`

The LZW decompressor added dictionary entries that directly referenced the `w` slice instead of making a copy. Old dictionary entries pointed to corrupted data.

- Explicit allocation of `newEntry` with copy of `w` before adding to dictionary

## Patch 4 — Incorrect RawSize in CZ block table
**File:** `czimage/util.go`

Critical bug causing visual CG corruption in-game (color artifacts). `Compress()` and `Compress2()` computed incorrect `RawSize` for each LZW block.

1. **Uncompensated carry-over**: The last LZW element carried to the next block was not deducted from the byte counter.
2. **Go UTF-8 encoding**: `len(string(byte(200)))` returns 2 instead of 1 for bytes > 127, causing ±1 errors on RawSize.

## Patch 5 — CZ4 image format support
**Files:** `czimage/cz4.go` (new), `czimage/imagefix.go`, `czimage/cz.go`

Added CZ4 format decoding and encoding, used in newer games (Little Busters English, LOOPERS, Harmonia, Kanon 2024).

CZ4 differs from CZ3 by storing RGB (w×h×3) and Alpha (w×h) channels separately, each with independent delta line encoding. LZW compression and blockHeight calculation are identical to CZ3.

## Patch 6 — PAK block alignment padding
**File:** `pak/pak.go`

After writing a rebuilt PAK (when replaced files are larger than originals), the file was not aligned to block size, potentially causing read errors.

- Added zero padding at end of file to align to `BlockSize`

## Patch 7 — AIR.py module resolution fix
**File:** `data/AIR.py`

The AIR.py definition script used `from base.air import *` to import functions from `data/base/air.py`. This import consistently failed in `script import` mode because LuckSystem's working directory is not `data/`.

- Merged all `base/air.py` functions directly into `AIR.py`
- Added the missing `ONGOTO` opcode handler
- Removed the `from base.air import *` dependency

## Tested games
- AIR (Steam) — full French translation pipeline, SYSCG.pak 51/51 (CZ3+CZ4), SCRIPT.pak import/export
- Summer Pockets — RawSize fix confirmed (masagrator report)
- Kanon — CZ2 font fix confirmed

## Credits
- **wetor** — original LuckSystem
- **masagrator** — RawSize bug identification (CZ3 layers)
- **G2-Games** — CZ4 reference (lbee-utils)
- **Yoremi** — all patches, AIR French translation, GUI

---
---

# V3.20 — Auto-sélection du plugin script + durcissement LOG_BEGIN dans la GUI Dialogue

13/06/2026

## Ajout : décompilation/import plus sûrs quand le plugin est oublié

### Contexte

Un testeur CartagraHD a signalé que des lignes `LOG_BEGIN` traduites semblaient ne pas s'appliquer après repack. Le dossier fourni montrait pourtant que les `.txt` prêts à repacker contenaient déjà le texte anglais, dont la toute première ligne visible :

```text
LOG_BEGIN ("The roar of water fills my ears.")
```

Les tests de round-trip ont confirmé que `LOG_BEGIN` est bien importé quand le plugin CartagraHD est chargé. Ce n'était donc pas un bug moteur confirmé ; le cas le plus probable était une mauvaise manipulation : repack avec l'OPCODE, mais sans le plugin Python correspondant, ce qui faisait tomber l'outil sur le parser générique.

### Fix / garde-fou

**CLI**

- `script decompile` et `script import` auto-sélectionnent maintenant le plugin Python frère quand `-O` pointe vers une arborescence standard et que `-p` est vide :
  - `data/GAME.txt` -> `data/GAME.py`
  - `data/GAME/OPCODE.txt` -> `data/GAME.py`
- Les opcodes dépendants du plugin, comme `MESSAGE`, `LOG_BEGIN` et `SELECT`, restent donc sur le chemin d'import/export spécifique au jeu même si l'utilisateur oublie de sélectionner manuellement le plugin.
- La CLI affiche une ligne explicite, par exemple :

```text
[INFO] Auto-selected plugin: data\CartagraHD.py
```

**GUI**

- La détection des lignes Dialogue est plus stricte :
  - accepte `LOG_BEGIN` derrière les préfixes `labelN:` et `globalN:`;
  - ne confond plus `MESSAGE_CLEAR`, `MESSAGE_WAIT` et les opcodes similaires avec de vrais `MESSAGE`;
  - aligne l'aide visible de la GUI avec le backend : `MESSAGE`, `LOG_BEGIN` et `SELECT`.
- Ajout d'un test backend GUI ciblé couvrant extraction et import de `LOG_BEGIN` normal, labellisé et global-labellisé.
- Les scripts npm frontend lancent maintenant Vite via `node ./node_modules/vite/bin/vite.js` au lieu d'exécuter directement le shim `.bin/vite`. Cela évite `sh: 1: vite: Permission denied` lors des builds Linux si `node_modules/.bin/vite` a perdu son bit exécutable.
- Passage des libellés CLI et GUI en `v3.20`.

### Tests réalisés

- `go test ./... -count=1` depuis `SourcesGUI-wails` : OK.
- `go test ./cmd ./script` : OK.
- `go test ./... -run '^$'` : OK.
- `npm run build` depuis `SourcesGUI-wails/frontend` : OK.
- Repack du cas CartagraHD fourni sans passer `-p`; auto-sélection de `data\CartagraHD.py`.
- Redécompilation du PAK obtenu : `0000-op0_HD.txt:144` contient bien :

```text
LOG_BEGIN ("The roar of water fills my ears.")
```

---

# V3.1.9 — Correction ONGOTO CartagraHD + support multi-goto + fix dump chaîne longueur zéro

06/06/2026

## Bug corrigé : choix CartagraHD cassés après traduction (offsets ONGOTO non recalculés)

### Problème

La traduction des scripts CartagraHD produisait des choix en jeu incorrects : les offsets de branche étaient exportés comme nombres bruts et restaient figés quand une ligne de dialogue avant un choix grossissait en longueur. Deux bugs indépendants en étaient la cause :

1. `ONGOTO` n'était pas défini dans `base/cartagrahd.py` — l'opcode tombait en `UNDEFINED()`, et ses cibles de saut étaient dumpées comme valeurs `uint16` brutes au lieu de références `{goto label...}`. La machinerie d'import ne recalculant les offsets qu'exprimés sous forme de labels, tout changement de taille après un ONGOTO laissait toutes les branches en aval pointer vers de mauvaises positions.

2. Le parser de labels du moteur ne gérait qu'un seul token `{goto ...}` par ligne. ONGOTO porte N cibles de branche sur la même ligne ; toutes les cibles après la première étaient silencieusement ignorées.

Un troisième bug mineur était présent : `operator/util.go` émettait un caractère parasite dans le dump de chaîne lorsqu'il rencontrait une entrée de longueur zéro.

### Fix (4 fichiers — CLI)

**`data/base/cartagrahd.py`**
- Ajout du handler `ONGOTO` : lit le nombre de branches N, puis lit N offsets `uint16` et les émet comme références `{goto label_NNNN}`, sur le modèle de `IFN`/`IFY`.

**`script/model.go`**
- Extension du modèle `JumpParam` / label pour stocker une slice de cibles de saut par ligne au lieu d'une seule cible, permettant à une ligne de script de contenir plusieurs tokens `{goto ...}`.

**`script/script.go`**
- Mise à jour de `Export()` et `Import()` pour itérer sur toutes les cibles de saut d'une ligne (pas seulement la première) lors de la construction et de la consommation des références de labels.
- `Import()` : recalcule chaque offset de cible indépendamment, de sorte que les N branches d'un ONGOTO sont correctement repointed après un changement de taille.

**`game/operator/util.go`**
- Correction du cas limite chaîne longueur zéro dans le helper de dump de chaînes : une entrée de longueur zéro n'émet plus de caractère parasite avant le délimiteur fermant.

### Tests réalisés

- `go test ./script ./game/operator` : OK.
- Round-trip CartagraHD original (sans traduction) : les 277 entrées internes du PAK sont byte-identiques ; hash PAK identique à l'original.
- Test de régression — ligne 46 allongée : les cibles ONGOTO passent de `3730 / 8104` à `3758 / 8132`, correspondant exactement aux deux nouvelles positions absolues.
- `go test ./...` : les échecs restants sont des échecs pre-existants sur fixtures absentes (FONT.PAK, SP.py, LOOPERS SCRIPT.PAK) — sans rapport avec ce patch.

### Note

Les anciens dumps CartagraHD contenant des nombres ONGOTO bruts (ex. `65530, 12835, …`) doivent être ré-extraits avec le plugin corrigé avant réimport. Les anciens dumps ne contiennent pas les tokens `{goto label}` nécessaires au recalcul des offsets.

### Jeux concernés

CartagraHD — tout script contenant ONGOTO (branches de choix).

---

# V3.1.3 — Patch 3 : GUI — Détection automatique des presets de jeu depuis data/

28/02/2026

## Amélioration : scan automatique des presets de jeu remplace la sélection manuelle

### Problème
Le dropdown Game de la GUI (ajouté au Patch 2) ne listait que LB_EN et SP. Les utilisateurs travaillant sur d'autres jeux (AIR, KANON, HARMONIA, LOOPERS, LUNARiA, PlanetarianSG, CartagraHD) devaient parcourir manuellement les fichiers OPCODE et plugin à chaque fois.

### Fix (4 fichiers — GUI uniquement)

**`app.go`**
- Ajout du struct `GamePreset` et de la fonction `ScanGameData()`
- Scanne le dossier `data/` à côté de l'exécutable lucksystem
- Détecte tous les fichiers `.txt` (définitions OPCODE) récursivement, excluant `data/base/`
- Associe chaque jeu à son plugin `.py` si présent à `data/GAME.py`
- Retourne une liste triée de presets avec chemins absolus

**`frontend/src/App.svelte`**
- Dropdown dynamique "Game preset" dans les formulaires Decompile et Compile, peuplé depuis `ScanGameData()`
- La sélection d'un preset remplit automatiquement Opcode, Plugin et Game
- Boutons de parcours manuel toujours disponibles (réinitialisent le preset à "— Manual —")
- Presets re-scannés quand le chemin lucksystem est changé via "Locate"

**`frontend/wailsjs/go/main/App.js`** + **`App.d.ts`** — ajout du binding `ScanGameData()`

### Presets détectés (dossier data/ standard)
AIR (plugin), CartagraHD (plugin), HARMONIA (plugin), KANON (plugin), LB_EN, LOOPERS (plugin), LUNARiA (plugin), PlanetarianSG (plugin), SP (plugin)

### Jeux concernés
Tous les jeux — améliore le workflow pour chaque jeu supporté.

---

# V3.1.3 — Patch 2 : Flag `--game` / `-g` pour forcer le type de jeu (CLI + GUI)

28/02/2026

## Fix : override explicite du type de jeu pour fiabilité multiplateforme

### Problème
Sous Linux, l'auto-détection du Patch 1 pouvait échouer si le fichier OPCODE était placé dans un dossier arbitraire (ex: `~/Bureau/OPCODE.txt` au lieu de `data/LB_EN/OPCODE.txt`). Le nom du dossier parent ne correspondait pas à "LB_EN" → retombe en "Custom" → opérateur générique → MESSAGE en codepoints bruts.

### Fix — CLI (3 fichiers)

**`cmd/script.go`**
- Ajout de la variable `ScriptGameName` et du flag persistant `--game`/`-g` sur la commande `script`

**`cmd/scriptDecompile.go`**
- Nouvelle fonction `resolveGameName()` : chaîne de priorité flag `--game` → auto-détect depuis chemin OPCODE → fallback "Custom"
- `detectGameName()` amélioré avec 2 stratégies : (1) match dossier parent (original), (2) recherche dans tout le chemin (nouveau)

**`cmd/scriptImport.go`**
- Utilise `resolveGameName()` partagé, suppression logique dupliquée et import `fmt` inutilisé

### Fix — GUI (4 fichiers)

**`app.go`**
- Paramètre `gameName` ajouté aux signatures de `ScriptDecompile()` et `ScriptCompile()`
- Passe `-g gameName` à lucksystem quand non vide

**`frontend/src/App.svelte`**
- Variable `gameName` et dropdown Game (Auto-detect / LB_EN / SP) dans les formulaires Decompile et Compile

**`frontend/wailsjs/go/main/App.js`** + **`App.d.ts`** — bindings mis à jour pour le nouveau paramètre

### Chaîne de priorité
```
Flag --game (-g LB_EN)  →  priorité maximale (override explicite)
Auto-détect depuis -O   →  Stratégie 1 : dossier parent, Stratégie 2 : dans tout le chemin
Fallback                →  "Custom" → NewGeneric()
```

### Utilisation
```bash
# Méthode 1 : Flag explicite (plus fiable, surtout sous Linux)
lucksystem script decompile -s SCRIPT.PAK -c UTF-8 -o output -O OPCODE.txt -g LB_EN

# Méthode 2 : Auto-détect depuis le chemin (fonctionne si LB_EN apparaît dans le chemin)
lucksystem script decompile -s SCRIPT.PAK -c UTF-8 -o output -O /path/LB_EN/OPCODE.txt
```

### Jeux concernés
LB_EN et SP — tout jeu utilisant un opérateur Go sans plugin Python.

---

# V3.1.3 — Patch 1 : Correction auto-détection GameName pour décompilation scripts

27/02/2026

## Bug corrigé : opcodes MESSAGE/SELECT/BATTLE exportés en codepoints numériques au lieu de texte

### Problème
La décompilation des scripts LB_EN avec `lucksystem script decompile -s SCRIPT.PAK -O data/LB_EN/OPCODE.txt` produisait des lignes MESSAGE avec des codepoints Unicode bruts au lieu de texte lisible :

```
MESSAGE (0, 12502, 12523, 12523, 12523, 12523, 8230, 12288, ...)
```

Sortie attendue :
```
MESSAGE (0, "ブルルルル…　ブルルルル…", "Burururururu...  Burururururu...", 0x5)
```

Les 161 scripts étaient affectés — aucun texte de dialogue visible, uniquement des séquences numériques. Le même problème affectait `script import` (l'aller-retour échouait).

### Cause racine — `GameName: "Custom"` codé en dur dans scriptDecompile.go / scriptImport.go

Les deux fichiers passaient `GameName: "Custom"` à `game.NewGame()`, quel que soit le chemin OPCODE fourni. Chaîne de dispatch :

```
scriptDecompile.go:  GameName: "Custom"          ← toujours en dur
        ↓
vm.go NewVM():  switch "Custom" → pas de match (ni "LB_EN" ni "SP")
        ↓
vm.Operate = nil → fallback NewGeneric()         ← filet de sécurité patch 15
        ↓
Generic n'a pas de méthode MESSAGE() → dispatch vers UNDEFINED()
        ↓
UNDEFINED() appelle AllToUint16() → dump des codepoints en nombres
```

L'opérateur `LB_EN` (`operator/LB_EN.go`) implémente déjà complètement MESSAGE, SELECT, BATTLE, TASK, SAYAVOICETEXT et VARSTR_SET avec les appels `DecodeString()` — il n'était simplement jamais instancié car le GameName ne correspondait jamais à `"LB_EN"` dans le switch.

### Note sur le patch 15 (v3.1)
Le patch 15 documentait "auto-détection du GameName depuis le chemin OPCODE" mais l'implémentation était incomplète — la fonction `detectGameName()` n'était pas présente dans les fichiers livrés. Le fallback generic ajouté au patch 15 empêchait le crash nil pointer mais ne résolvait pas le décodage du texte. Ce patch complète l'auto-détection.

### Fix (2 fichiers)

**`cmd/scriptDecompile.go`**
- Ajout de `detectGameName(opcodePath string) string` : extrait le dossier parent du chemin OPCODE via `filepath.Dir()` + `filepath.Base()`, comparaison insensible à la casse avec les jeux connus (`LB_EN`, `SP`)
- L'auto-détection ne s'exécute que si aucun fichier plugin (`-p`) n'est fourni (les plugins ont la priorité)
- Affiche `[INFO] Auto-detected game: LB_EN (from OPCODE path)` quand un match est trouvé
- Remplacement de `GameName: "Custom"` par `GameName: gameName`

**`cmd/scriptImport.go`**
- Même logique d'auto-détection via `detectGameName()` (définie dans `scriptDecompile.go`, même package `cmd`)
- Remplacement de `GameName: "Custom"` par `GameName: gameName`

### Chaîne de priorité
```
Plugin (-p fichier.py)  →  priorité maximale (toujours utilisé si fourni)
Auto-détect depuis -O   →  "data/LB_EN/OPCODE.txt" → "LB_EN"
Fallback                →  "Custom" → NewGeneric()
```

### Jeux concernés
LB_EN et SP — tout jeu utilisant un opérateur Go (pas un plugin Python) avec son fichier OPCODE dans un sous-dossier correspondant au nom du jeu.

---

# V3.1.2 — Patch 1 : Correction séparateurs de chemins PAK Import/Export (Windows)

26/02/2026

## Bug corrigé : crash `pak replace` en mode dossier + corruption CZ via chemins mixtes

### Problème 1 — Crash en mode dossier
`lucksystem pak replace -s OTHCG.PAK -i <dossier> -o out.PAK` crashait avec `strconv.Atoi: parsing "C:\Users\...\msg_01k_en": invalid syntax`. Tous les fichiers du dossier étaient skippés avec "Skip File", puis la dernière erreur `strconv.Atoi` se propageait en crash fatal.

### Problème 2 — Corruption CZ en mode liste (anciens fichiers liste)
La réinjection de fichiers CZ dans un PAK via un fichier liste généré par une ancienne version produisait des images CZ corrompues (non-remplacées) en jeu. Les fichiers remplacés eux-mêmes étaient corrects — seuls les fichiers voisins non modifiés montraient de la corruption.

### Cause racine — `path.Base()` / `path.Join()` vs `filepath` (Windows)

Le fichier `pak.go` utilisait le package `path` (POSIX uniquement, séparateur `/`) au lieu de `path/filepath` (séparateur natif de l'OS) à trois endroits :

**Import() ligne 532** — `path.Base(file)` : Sous Windows, `filepath.Walk` retourne des chemins avec `\`. `path.Base("C:\...\msg_01k_en")` traite `\` comme des caractères normaux et retourne le **chemin complet** comme nom de fichier. Résultat : `CheckName` échoue toujours → `strconv.Atoi(chemin_complet)` → crash.

**Export() lignes 412/415** — `path.Join(dir, name)` : Générait des fichiers liste avec des séparateurs mixtes (`C:\Users\...\OTHCG_extracted/aug_01`). Quand ces fichiers liste étaient ensuite utilisés pour `pak replace -l`, les chemins mixtes causaient des échecs subtils de correspondance lors du rebuild du PAK, corrompant les calculs d'offset pour les fichiers non-remplacés.

### Bugs additionnels corrigés

**Fuite de variable erreur** — `id, err = strconv.Atoi(name)` écrivait dans le `err` du scope externe. Après un `continue` sur le dernier fichier, `return err` en fin de fonction retournait l'erreur Atoi au lieu de `nil`. Fix : variable locale `parseErr` via `:=`.

**Fuite de descripteur de fichier** — `fs, _ := os.Open(file)` ouvrait les fichiers avant validation. Quand un fichier était skippé via `continue`, le descripteur n'était jamais fermé. Fix : `fs.Close()` avant chaque `continue` + gestion d'erreur sur `os.Open`.

### Fix (1 fichier)

**`pak/pak.go`**
- `Import()` mode `"dir"` : `path.Base()` → `filepath.Base()` ; variables locales `parseErr`/`openErr` ; `fs.Close()` sur tous les chemins de skip
- `Export()` mode `"all"` : `path.Join()` → `filepath.Join()` (lignes 412, 415)
- Suppression de l'import `"path"` inutilisé

### Jeux concernés
Tous les jeux — toute opération PAK replace sous Windows était affectée.

---

# V3.1 — Patch 2 : Réduction de la verbosité des warnings d'opcodes indéfinis

25/02/2026

## Amélioration : accumulation silencieuse des warnings d'opcodes indéfinis

### Problème
Lors du `script decompile` sur Little Busters EN, les 1 461 warnings d'opcodes indéfinis (`Operation不存在 HAIKEI_SET`, etc.) étaient affichés un par un via `glog.V(5).Infoln()`. Sur des machines lentes ou en utilisant la GUI, cela créait une boucle apparemment infinie — la console scrollait des warnings pendant plus de 2 minutes, donnant l'impression que l'outil était bloqué. La décompilation elle-même ne prend que ~5 secondes.

### Fix (3 fichiers)

**`game/operator/undefined_operate.go`** — Remplacement de l'appel `glog.V(5).Infoln()` par opcode par un `opcodeTracker` thread-safe qui accumule silencieusement les compteurs dans une `map[string]int`. Expose `PrintUndefinedOpcodeSummary()` qui affiche un seul bloc résumé trié après traitement, puis réinitialise le tracker.

**`cmd/scriptDecompile.go`** — Ajout de l'appel `operator.PrintUndefinedOpcodeSummary()` après `g.RunScript()`.

**`cmd/scriptImport.go`** — Même appel résumé après `g.RunScript()`.

### Résultat
Au lieu de 1 461 lignes de warning individuelles :
```
[INFO] 1461 undefined opcodes skipped (15 unique types):
  HAIKEI_SET            x312
  WAIT                  x245
  DRAW                  x198
  ...
These are non-text opcodes (visual/audio/system) and can be safely ignored for translation work.
```

### Jeux concernés
Tous les jeux utilisant l'opérateur générique ou tout jeu avec des opcodes indéfinis (Little Busters EN, Kanon, Harmonia, LOOPERS, LUNARiA, Planetarian).

---

# V3.1 — Patch 1 : Correction décompilation scripts Little Busters EN

24/02/2026

## Bug corrigé : crash nil pointer + entrées de scripts invalides dans `script decompile`

### Problème
`lucksystem script decompile -s SCRIPT.PAK -c UTF-8 -o WORK -O data\LB_EN\OPCODE.txt` crashait immédiatement avec un nil pointer dereference. Deux bugs indépendants empêchaient Little Busters EN (et tout jeu sans flag `-p` plugin) d'être décompilé.

### Cause racine 1 — Nil pointer dans `NewVM()` sur GameName inconnu
`scriptDecompile.go` codait en dur `GameName: "Custom"` quand aucun `-p` plugin n'était spécifié. Dans `vm.go:NewVM()`, le switch ne matchait que `"LB_EN"` et `"SP"` — toute autre valeur laissait `vm.Operate = nil`. La ligne 48 appelait ensuite `vm.Operate.Init(vm.Runtime)` → panic nil pointer dereference.

### Cause racine 2 — SEEN8500/SEEN8501 (tables de données parsées comme scripts)
SCRIPT.PAK contient 169 entrées mais seulement 167 sont des scripts. SEEN8500 et SEEN8501 sont des tables de données du mini-jeu de baseball avec `firstLen=0` (premiers 2 octets des données). Quand `script.LoadScript()` appelait `restruct.Unpack` avec `Len=0`, il calculait `size = Len - 4` → underflow non signé → crash.

### Fix (5 fichiers)

**`game/operator/generic.go` (nouveau)** — Opérateur fallback générique qui gère les opcodes courants (IFN, IFY, GOTO, JUMP, FARCALL, GOSUB, EQU, EQUN, ADD, RANDOM) via `LucaOperateDefault` + `LucaOperateExpr` embarqués. Les opcodes inconnus sont dumpés en `UNDEFINED` avec les paramètres en uint16.

**`game/VM/vm.go`** — Nil guard après le switch GameName : si `vm.Operate` est toujours nil, instanciation de `operator.NewGeneric()` avec un warning. Cela empêche le crash nil pointer pour tout jeu non supporté.

**`cmd/scriptDecompile.go`** — Auto-détection du GameName depuis le chemin OPCODE : `data\LB_EN\OPCODE.txt` → dossier parent `LB_EN` → utilise l'opérateur spécifique LB_EN (handlers MESSAGE, SELECT, BATTLE, TASK). Retombe sur `"Custom"` si le nom de dossier ne correspond à aucun jeu connu.

**`cmd/scriptImport.go`** — Même logique d'auto-détection que scriptDecompile.go.

**`game/game.go`** — Ajout de `isValidScript()` (rejette les entrées avec `firstLen < 4`) et `safeLoadScript()` (recovery de panic). SEEN8500/SEEN8501 sont désormais skippés avec un warning au lieu de crasher.

### Résultat
- 161 scripts décompilés avec succès (SEEN8500/8501 skippés avec warning)
- 102 795 lignes MESSAGE extraites (format bilingue JP/EN)
- 1 461 warnings d'opcodes non gérés pour les opcodes visuels/audio (HAIKEI_SET, INIT, DRAW, WAIT, BGM, SE…) — attendu, n'affecte pas l'extraction du texte

### Jeux concernés
Tout jeu sans fichier plugin Python dédié, dont Little Busters EN, Kanon, Harmonia, LOOPERS, LUNARiA, Planetarian.

---

# V3 — Patch 3 : Correction import CZ2 (fonts) + améliorations GUI

22/02/2026

## Bug corrigé : crash à la réimport CZ2 (`font edit`)

### Problème
`lucksystem font edit` crashait avec `invalid argument` en modes **append** ou **insert**. Cause racine : `ReplaceChars()` augmente la hauteur de l'image quand on ajoute des caractères, mais `CzHeader.Width/Heigth` n'était jamais mis à jour. `Cz2Image.Import()` alors :
1. Rognait la nouvelle image (plus haute) aux anciennes dimensions via `FillImage()`
2. Atteignait le chemin de mismatch de taille et retournait `nil` — silencieusement, car `err` était déclaré mais jamais assigné
3. Laissait `Raw` et `OutputInfo` vides → `WriteStruct` appelait `restruct.Pack` sur une struct vide → `invalid argument`

### Fix
- `czimage/cz2.go` — `Import()` : quand les dimensions diffèrent (cas valide après append/insert), mise à jour de `CzHeader.Width/Heigth` au lieu du retour silencieux. Remplacement de l'assertion directe `*image.NRGBA` par une conversion sûre pour tout format PNG.
- `czimage/cz2.go` — ajout de la méthode `SetDimensions(w, h uint16)` sur `Cz2Image`
- `font/font.go` — `Write()` : appel de `SetDimensions()` avant `CzImage.Import()` pour synchroniser le header avec l'image produite par `ReplaceChars()`

### Améliorations GUI (v3 GUI)
- **Fenêtre CMD cachée** : les appels subprocess n'affichent plus de popup console sur Windows (`SysProcAttr{HideWindow: true}` — fichier build platform-specific)
- **Bouton Stop** : un bouton **■ Stop** apparaît dans le header de la console pendant toute opération en cours ; annule le subprocess via `context.WithCancel`
- **Champs de sortie Font Edit** : Output CZ et Output info sont maintenant en saisie libre (plus de validation d'extension Windows via SaveFileDialog)
- **PAK Font Replace** : ajout du mode fichier liste (identique à CG Replace) pour éviter le crash `strconv.Atoi` en mode dossier

---

# V3 — Patch 2 : LuckSystem GUI

21/02/2026

## Nouveau : Interface graphique (dépôt séparé)

Une GUI standalone pour LuckSystem construite avec **Wails** (Go + Svelte). La GUI appelle `lucksystem.exe` via subprocess, conformément à la [recommandation architecturale de wetor](https://github.com/wetor/LuckSystem) — aucun code source de LuckSystem n'est embarqué.

### Fonctionnalités
- **8 opérations** : Script Decompile/Compile, PAK Extract/Replace, Font Extract/Edit, Image Export/Import
- **Mode batch** pour Image Export/Import (dossiers entiers)
- Console temps réel avec logs colorés
- Détection automatique de `lucksystem.exe` (même dossier, CWD, ou PATH)
- Icône personnalisée
- Interface inspirée de SiglusTools

### Architecture
```
LuckSystemGUI.exe  ←→  lucksystem.exe (subprocess)
   (Wails/Go)              (outil CLI)
```

### Compilation
Requis : Go 1.23+, Node.js, Wails CLI
```bash
cd frontend && npm install && cd ..
wails build
```

---

# V3 — Patch 1 : CZ1 32-bit Import/Export + CZ0 logging

20/02/2026

## Fichiers modifiés
- `czimage/cz1.go` — réécriture Import/Export/Write
- `czimage/cz.go` — gestion gracieuse des fichiers non-CZ
- `czimage/cz0.go` — ajout log V(0) dans decompress()

## Bugs corrigés

### 1. Extended header manquant dans Write()
Le `Write()` original n'écrivait que les 15 bytes du `CzHeader` struct, ignorant les 13 bytes d'extended header (offsets, crop, bounds). Le fichier produit avait la block table à l'offset 15 au lieu de 28 → crash à la relecture.

**Fix** : Sauvegarde des bytes raw 15→HeaderLength dans `ExtendedHeader` au `Load()`, réécriture dans `Write()`.

### 2. Import() ne gérait que l'alpha
L'Import 32-bit ne compressait que le canal A (`data[i] = pic.A`), jetant RGB. Résultat : écran blanc/transparent en jeu.

**Fix** : Import multi-mode selon Colorbits (4, 8, 24, 32). Le mode 32-bit fait une copie directe de `pic.Pix` (RGBA).

### 3. Colorbits > 32 (palette 8-bit)
Les fichiers CZ1 palette utilisent Colorbits=248 (0xF8), un marqueur propriétaire Visual Art's. LuckSystem ne le reconnaissait pas → palette ignorée → `GetOutputInfo()` lisait la palette comme block table → crash (`slice bounds out of range`).

**Fix** : Normalisation `if Colorbits > 32 → Colorbits = 8` (même approche que lbee-utils).

### 4. Fichiers non-CZ dans les PAK
Les fichiers sans magic "CZ" (ex: トーンカーブ_夕/夜, des LUTs 768 bytes) causaient un `glog.Fatalln("Unknown Cz image type")`.

**Fix** : Vérification du magic avant unpacking, retour `nil` avec warning au lieu de crash.

### 5. Palette BGRA dans Write()
La palette est lue en BGRA (fichier) et stockée en NRGBA (Go). L'ancien Write via `restruct` sérialisait en RGBA → couleurs inversées R↔B.

**Fix** : Écriture manuelle de chaque entrée palette en [B,G,R,A].

### 6. CZ0 invisible dans les logs d'extraction
Les fichiers CZ0 n'avaient que du logging `V(6)` (debug profond), alors que CZ4 log en `V(0)` (toujours visible). Lors de l'extraction d'un PAK contenant un mix CZ0/CZ4, les dernières lignes visibles avant un CZ0 provenaient du CZ4 précédent, donnant l'impression que les CZ0 étaient traités comme CZ4.

**Fix** : Ajout d'un `glog.V(0).Infof("Decompress CZ0: %dx%d, Colorbits=%d")` dans `cz0.go:decompress()` (ligne 78).

## Format CZ1 confirmé
- 32-bit : pixels stockés en **RGBA** (pas BGRA comme CZ3)
- 8-bit palette : entrées stockées en **BGRA**, données = 1 byte/pixel (index)
- Extended header : 13 bytes obligatoires (même structure que Cz3Header)

## Statut
- ✅ CZ1 32-bit : round-trip OK, testé en jeu (systemmenu FR)
- ✅ CZ1 8-bit palette : ok (system_icon, NUM files)
- ✅ Fichiers non-CZ : warning au lieu de crash
- ✅ CZ0 : correctement identifié dans les logs d'extraction

### Mergé upstream (PR #35)
- **Fix crash décompresseur CZ2 (fonts)** — `czimage/lzw.go` (vérification de limites dans `decompressLZW2`)

---

# Version 2 (7 patches)

18/02/2026

## Patch 1 — Import de scripts à longueur variable
**Fichier :** `script/script.go`

L'import de scripts traduits échouait avec un panic quand la traduction avait une longueur différente de l'original. Le code vérifiait strictement `len(paramList) == len(code.Params)`, ce qui bloquait toute traduction plus longue ou plus courte.

- Suppression de la vérification stricte du nombre de paramètres
- Ajout de bounds checking dans la boucle de conversion et le merge des paramètres
- Les offsets de jump (GOTO, IFN, IFY…) sont recalculés automatiquement

## Patch 2 — Correction du pipeline CZ3 (export/import PNG)
**Fichiers :** `czimage/cz3.go`, `czimage/imagefix.go`

L'export et l'import de CZ3 corrompaient silencieusement les données pixels.

- **Magic byte** : `Write()` écrasait le magic "CZ3" → "CZ0", rendant le fichier illisible par le jeu
- **Format NRGBA** : Conversion automatique de tout format PNG en NRGBA 32-bit avant encodage
- **Buffer aliasing** : `DiffLine()` et `LineDiff()` partageaient des slices au lieu de copier, provoquant une corruption des données delta

## Patch 3 — Corruption mémoire dans le décompresseur LZW
**Fichier :** `czimage/lzw.go`

Le décompresseur LZW ajoutait des entrées dictionnaire qui référençaient directement le slice `w` au lieu d'en faire une copie. Les anciennes entrées du dictionnaire pointaient vers des données corrompues.

- Allocation explicite de `newEntry` avec copie de `w` avant ajout au dictionnaire

## Patch 4 — RawSize incorrect dans la table de blocs CZ
**Fichier :** `czimage/util.go`

Bug critique causant la corruption visuelle des CG en jeu (artefacts colorés). Les fonctions `Compress()` et `Compress2()` calculaient un `RawSize` erroné pour chaque bloc LZW.

1. **Carry-over non compensé** : Le dernier élément LZW reporté au bloc suivant n'était pas déduit du compteur de bytes.
2. **Encodage UTF-8 de Go** : `len(string(byte(200)))` retourne 2 au lieu de 1 pour les octets > 127, causant des erreurs ±1 sur les RawSize.

## Patch 5 — Support du format CZ4
**Fichiers :** `czimage/cz4.go` (nouveau), `czimage/imagefix.go`, `czimage/cz.go`

Ajout du décodage et de l'encodage du format CZ4, utilisé dans les jeux récents (Little Busters English, LOOPERS, Harmonia, Kanon 2024).

Le CZ4 diffère du CZ3 par le stockage séparé des canaux RGB (w×h×3) et Alpha (w×h), chacun avec son propre delta line encoding indépendant. Le LZW et le calcul de blockHeight sont identiques au CZ3.

## Patch 6 — Padding d'alignement dans pak.go
**Fichier :** `pak/pak.go`

Après l'écriture d'un PAK reconstruit (quand les fichiers remplacés sont plus grands que les originaux), le fichier n'était pas aligné sur la taille de bloc, ce qui pouvait causer des erreurs de lecture.

- Ajout de padding zéro en fin de fichier pour aligner sur `BlockSize`

## Patch 7 — Correction AIR.py (résolution du module base)
**Fichier :** `data/AIR.py`

Le script de définition AIR.py utilisait `from base.air import *` pour importer les fonctions de `data/base/air.py`. Cet import échouait systématiquement en mode `script import` car le working directory de LuckSystem n'est pas `data/`.

- Fusion des fonctions de `base/air.py` directement dans `AIR.py`
- Ajout de la fonction `ONGOTO` qui était absente
- Suppression de la dépendance `from base.air import *`

## Jeux testés
- AIR (Steam) — traduction française complète, SYSCG.pak 51/51 (CZ3+CZ4), SCRIPT.pak import/export
- Summer Pockets — fix RawSize confirmé (rapport masagrator)
- Kanon — fix CZ2 font confirmé

## Crédits
- **wetor** — LuckSystem original
- **masagrator** — identification du bug RawSize (CZ3 layers)
- **G2-Games** — référence CZ4 (lbee-utils)
- **Yoremi** — tous les patches, traduction française d'AIR, GUI
