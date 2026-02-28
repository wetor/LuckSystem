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
Any game without a dedicated Python plugin file, including Little Busters EN, Kanon, Harmonia, LOOPERS, LUNARiA, Planetarian.

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
