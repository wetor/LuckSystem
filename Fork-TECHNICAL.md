# V3.1 — Patch 1 : Correction décompilation scripts Little Busters EN

## Fichiers modifiés
- `game/operator/generic.go` — **nouveau** : opérateur fallback générique
- `game/VM/vm.go` — nil guard + instanciation fallback
- `cmd/scriptDecompile.go` — auto-détection GameName depuis chemin OPCODE
- `cmd/scriptImport.go` — même auto-détection
- `game/game.go` — validation des entrées scripts + panic recovery

## Bug 1 : nil pointer dereference dans `NewVM()` (vm.go:48)

### Contexte
La commande `script decompile` sans flag `-p` (plugin Python) passe par `scriptDecompile.go` qui construit les options VM avec `GameName: "Custom"`. La chaîne d'appel :

```
scriptDecompile.go → game.LoadGame() → game.load() → VM.NewVM(opts)
                                                          ↓
                                                   switch opts.GameName {
                                                   case "LB_EN": vm.Operate = operator.NewLB_EN()
                                                   case "SP":    vm.Operate = operator.NewSP()
                                                   }
                                                   // "Custom" → aucun case → vm.Operate = nil
                                                   vm.Operate.Init(vm.Runtime)  ← PANIC
```

### Stack trace
```
goroutine 1 [running]:
lucksystem/game/VM.(*VM).NewVM(...)
    game/VM/vm.go:48
lucksystem/game.(*Game).load(...)
    game/game.go:53
lucksystem/cmd.scriptDecompile(...)
    cmd/scriptDecompile.go:26
```

### Fix

**`cmd/scriptDecompile.go` — auto-détection**
```go
func detectGameName(opcodePath string) string {
    if opcodePath == "" { return "Custom" }
    dir := filepath.Dir(opcodePath)
    dirName := strings.ToUpper(filepath.Base(dir))
    knownGames := []string{"LB_EN", "SP"}
    for _, g := range knownGames {
        if dirName == g { return g }
    }
    return "Custom"
}
```
Avec `data\LB_EN\OPCODE.txt` → `filepath.Dir` = `data\LB_EN` → `filepath.Base` = `LB_EN` → match → utilise opérateur LB_EN (MESSAGE, SELECT, BATTLE, TASK, SAYAVOICETEXT, VARSTR_SET, IMAGELOAD, MOVE).

**`game/VM/vm.go` — nil guard**
```go
if vm.Operate == nil {
    glog.Warningf("No game-specific operator for '%s', using generic fallback\n", opts.GameName)
    vm.Operate = operator.NewGeneric()
}
vm.Runtime = runtime.NewRuntime(opts.Mode)
vm.Operate.Init(vm.Runtime)
```

**`game/operator/generic.go` — opérateur générique**
```go
type Generic struct {
    LucaOperateUndefined  // UNDEFINED dump pour opcodes inconnus
    LucaOperateDefault    // IFN, IFY, GOTO, JUMP, FARCALL, GOSUB
    LucaOperateExpr       // EQU, EQUN, ADD, RANDOM
}

func (g *Generic) Init(ctx *runtime.Runtime) {
    ctx.Init(charset.ShiftJIS, charset.Unicode, true)
}
```
Implémente l'interface `api.Operator` (seule méthode requise : `Init`). Les opcodes courants sont gérés par les structs embarquées, les opcodes inconnus sont dumpés via `LucaOperateUndefined`.

## Bug 2 : SEEN8500/SEEN8501 — tables de données parsées comme scripts

### Contexte
SCRIPT.PAK de Little Busters contient 169 entrées dans sa table d'index (2916 slots, dont 169 non-nuls). Parmi elles, SEEN8500 et SEEN8501 ne sont **pas** des scripts mais des tables de données du mini-jeu de baseball.

### Analyse du format PAK
```
SCRIPT.PAK header: 2916 entries, block_size=4, flags=0x200

Script normal (SEEN0513):
  offset=0, length=49052, data: [25 00 5b 01 ...]
  → firstLen = 0x0025 = 37 bytes (longueur de la première CodeLine)

Table de données (SEEN8500):
  offset=17210032, length=7962, data: [00 00 5b 01 ...]
  → firstLen = 0x0000 = 0 (pas un script!)

Table de données (SEEN8501):
  offset=17217996, length=52724, data: [00 00 5b 01 ...]
  → firstLen = 0x0000 = 0 (pas un script!)
```

### Mécanisme du crash
Dans `script.LoadScript()`, le parsing d'une CodeLine lit `Len` (uint16) en premiers 2 octets, puis appelle `restruct.Unpack` avec une taille calculée comme `size = Len - 4`. Avec `Len = 0` :
```
size = uint16(0) - 4 = 65532 (underflow uint16)
→ restruct.Unpack tente de lire 65532 bytes → panic
```

### Fix

**`game/game.go` — validation + recovery**
```go
func isValidScript(data []byte) bool {
    if len(data) < 4 { return false }
    firstLen := binary.LittleEndian.Uint16(data[0:2])
    if firstLen < 4 { return false }  // Minimum CodeLine = 4 bytes (len + opcode + flag)
    return true
}

func safeLoadScript(opts *script.LoadOptions) (scr *script.Script, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("failed to parse script '%s': %v", opts.Entry.Name, r)
            scr = nil
        }
    }()
    scr = script.LoadScript(opts)
    return scr, nil
}
```

Dans `load()` :
```go
if !isValidScript(entry.Data) {
    glog.Warningf("Skipping invalid script entry '%s' (firstLen < 4, likely data table)\n", entry.Name)
    continue
}
scr, loadErr := safeLoadScript(&script.LoadOptions{Entry: entry})
if loadErr != nil {
    glog.Warningf("Skipping script '%s': %v\n", entry.Name, loadErr)
    continue
}
```
Double protection : validation structurelle + panic recovery pour les cas imprévus.

## Résultat extraction Little Busters EN

| Métrique | Valeur |
|----------|--------|
| Entrées PAK | 169 (sur 2916 slots) |
| Scripts valides | 167 |
| Scripts skippés | 2 (SEEN8500, SEEN8501) |
| Scripts exportés | 161 (+ 6 scripts utilitaires : _ARFLAG, _VARSTR, _QUAKE, _SAYAVOICE, _KEYWORD, et scripts sans MESSAGE) |
| Lignes MESSAGE | 102 795 |
| Warnings opcodes | 1 461 (HAIKEI_SET, INIT, DRAW, WAIT, BGM, SE, FARRETURN, END, BTFUNC, KOEP, etc.) |

Les warnings sont des opcodes visuels/audio non implémentés dans l'opérateur LB_EN. Ils n'affectent pas l'extraction du texte — les scripts sont exportés correctement avec toutes les lignes MESSAGE au format bilingue `MESSAGE (id, "日本語", "English")`.

### Portée
Affecte tout jeu LucaSystem sans fichier plugin Python dédié. La protection est double :
1. **Auto-détection** : évite 99% des cas (l'utilisateur spécifie le bon chemin OPCODE)
2. **Generic fallback** : empêche le crash pour tout jeu inconnu (opérateur minimal)
3. **Validation scripts** : protège contre les entrées PAK non-script dans tout jeu

---

# V3 — Patch 3 : Correction import CZ2 (fonts)

## Fichiers modifiés
- `czimage/cz2.go` — `Import()` : gestion correcte du redimensionnement + assertion de type sûre
- `czimage/cz2.go` — ajout de `SetDimensions(w, h uint16)`
- `font/font.go` — `Write()` : synchronisation du `CzHeader` avant `Import()`

## Bug : crash `invalid argument` dans `font edit` (modes append/insert)

### Contexte
`lucksystem font edit` avec les modes `-a` (append) ou `-i N` (insert) modifie une police CZ2 en ajoutant de nouveaux glyphes. La chaîne d'appel complète est :

```
fontEdit.go → LucaFont.Import() → LucaFont.ReplaceChars() → LucaFont.Write()
                                                                     ↓
                                                         CzImage.Import(img, fillSize=true)
                                                                     ↓
                                                              CzImage.Write()
```

### Cause racine (3 bugs imbriqués)

**Bug 1 — Dimensions non synchronisées**
`ReplaceChars()` crée une nouvelle image de dimensions :
```go
imageW := size*100 + 4
imageH := size * int(math.Ceil(float64(f.Info.CharNum) / 100.0))  // augmente si CharNum > old
```
Mais `CzHeader.Width` et `CzHeader.Heigth` conservent les valeurs du fichier original. Lors de l'appel `f.CzImage.Import(img, true)`, ces vieilles dimensions sont utilisées comme cible.

**Bug 2 — Retour nil silencieux dans Cz2Image.Import()**
```go
func (cz *Cz2Image) Import(r io.Reader, fillSize bool) error {
    var err error          // err = nil, jamais assigné sur ce chemin
    // ...
    pic = FillImage(pic, oldWidth, oldHeight)  // tronque la nouvelle image !
    if width != pic.Rect.Size().X || height != pic.Rect.Size().Y {
        return err         // retourne nil — silencieusement
    }
    // Le reste (compression) n'est jamais atteint
    // Raw et OutputInfo restent nil/vides
```
La fonction retourne `nil` (succès apparent) mais `cz.Raw` et `cz.OutputInfo` restent dans leur état initial non-compressé.

**Bug 3 — WriteStruct sur OutputInfo vide**
```go
func (cz *Cz2Image) Write(w io.Writer) error {
    err = WriteStruct(w, &cz.CzHeader, cz.Cz2Header, cz.ColorPanel, cz.OutputInfo)
    // cz.OutputInfo = nil → restruct.Pack → "invalid argument"
```

### Fix

**`czimage/cz2.go` — Import()**
```go
// AVANT : retour silencieux si dimensions différentes
if width != pic.Rect.Size().X || height != pic.Rect.Size().Y {
    glog.V(2).Infof("图片大小不匹配...")
    return err  // err = nil, Raw reste vide → crash dans Write()
}

// APRÈS : mise à jour du header pour accepter les nouvelles dimensions
if width != pic.Rect.Size().X || height != pic.Rect.Size().Y {
    // Cas légitime : ReplaceChars a agrandi l'image (append/insert)
    cz.Width  = uint16(pic.Rect.Size().X)
    cz.Heigth = uint16(pic.Rect.Size().Y)
    width  = pic.Rect.Size().X
    height = pic.Rect.Size().Y
    glog.V(2).Infof("CZ2 dimensions updated: %dx%d -> %dx%d", ...)
}
// → suite du code : compression normale, Raw et OutputInfo correctement remplis
```

**`czimage/cz2.go` — SetDimensions()**
```go
func (cz *Cz2Image) SetDimensions(w, h uint16) {
    cz.CzHeader.Width  = w
    cz.CzHeader.Heigth = h
}
```

**`font/font.go` — Write()**
```go
// Sync CzHeader avec les dimensions de l'image AVANT Import()
if setter, ok := f.CzImage.(interface {
    SetDimensions(w, h uint16)
}); ok {
    setter.SetDimensions(
        uint16(f.Image.Bounds().Size().X),
        uint16(f.Image.Bounds().Size().Y),
    )
}
err = f.CzImage.Import(img, true)
```

**`czimage/cz2.go` — safe type assertion**
```go
// AVANT : assertion directe, panic si PNG n'est pas NRGBA
pic := cz.PngImage.(*image.NRGBA)

// APRÈS : conversion sûre pour tout format PNG
var pic *image.NRGBA
switch src := cz.PngImage.(type) {
case *image.NRGBA:
    pic = src
default:
    // Conversion explicite pixel par pixel
    dst := image.NewNRGBA(src.Bounds())
    // ...
    pic = dst
}
```

### Portée
Affecte uniquement `font edit` en modes append (`-a`) et insert (`-i N`). Le mode redraw (`-r`) ne change pas les dimensions et n'est pas impacté. Le format CZ2 est utilisé exclusivement pour les polices de caractères (FONT.PAK, FONT__INFO.PAK).

### Jeux concernés
Tous les jeux Visual Art's/Key utilisant des polices CZ2 (AIR, Kanon, CLANNAD, Little Busters, etc.).

---



# V3 — Patch 1 : CZ1 32-bit Import/Export + CZ0 logging

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

**Fix** : Ajout d'un `glog.V(0).Infof("Decompress CZ0: %dx%d, Colorbits=%d")` dans `cz0.go:decompress()` (ligne 78) pour identifier clairement le format dans les logs.

## Format CZ1 confirmé
- 32-bit : pixels stockés en **RGBA** (pas BGRA comme CZ3)
- 8-bit palette : entrées stockées en **BGRA**, données = 1 byte/pixel (index)
- Extended header : 13 bytes obligatoires (même structure que Cz3Header)

## Statut
- ✅ CZ1 32-bit : round-trip OK, testé en jeu (systemmenu FR)
- ⏳ CZ1 8-bit palette : code prêt, à tester (system_icon, NUM files)
- ✅ Fichiers non-CZ : warning au lieu de crash
- ✅ CZ0 : correctement identifié dans les logs d'extraction


# Technical Analysis — LuckSystem-Yoremi-version 2

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
| `game/operator/generic.go` | 15 | New: generic fallback operator |
| `game/VM/vm.go` | 15 | Nil guard + generic fallback instantiation |
| `cmd/scriptDecompile.go` | 15 | Auto-detection of GameName from OPCODE path |
| `cmd/scriptImport.go` | 15 | Same auto-detection as scriptDecompile |
| `game/game.go` | 15 | isValidScript() + safeLoadScript() |
| `czimage/cz2.go` | 13 | CZ2 Import() resize fix + SetDimensions() |
| `font/font.go` | 13 | Write() header sync before Import() |
| `czimage/cz1.go` | 8-10 | CZ1 32-bit + 8-bit Import/Export rewrite |
| `czimage/cz.go` | 10, 5 | Non-CZ graceful handling + CZ4 dispatcher |
| `czimage/cz0.go` | 11 | CZ0 logging visibility |
| `script/script.go` | 1 | Variable-length script import |
| `czimage/cz3.go` | 2 | Magic byte, NRGBA conversion, logging |
| `czimage/imagefix.go` | 2, 5 | Buffer aliasing fix + CZ4 delta encode/decode |
| `czimage/lzw.go` | 3 | LZW dictionary memory corruption |
| `czimage/util.go` | 4 | RawSize carry-over + UTF-8 length |
| `czimage/cz4.go` | 5 | New: CZ4 format support |
| `pak/pak.go` | 6 | PAK block alignment padding |
| `data/AIR.py` | 7 | Self-contained module, no base/ dependency |

## Known remaining issues

- **CZ1 8-bit palette**: `GetOutputInfo()` does not account for the 1024-byte color palette between the header and block table for CZ1 images with `Colorbits=8` or `Colorbits=4`. This causes a panic on palette-indexed CZ1 files (e.g., `system_icon_*` in PARTS.pak).
- **CZ1 32-bit import**: `cz1.go Import()` only exports the alpha channel (`data[i] = A`), discarding RGB. Re-importing a CZ1 32-bit RGBA image produces a blank/white result (e.g., `systemmenu` in PARTS.pak).
- **Non-CZ files in PAK**: Some files extracted from PAK archives are not CZ images (e.g., トーンカーブ_夕/夜 are 768-byte RGB tone curve LUTs). These should be detected and skipped gracefully instead of crashing with "Unknown CZ image type".
