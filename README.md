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
- Script Decompile / Compile
- PAK Extract / Replace (CG and Font workflows separated)
- Font Extract / Edit (append, insert, redraw modes)
- Image Export / Import (single file + batch folder mode)
- Real-time console output
- **Stop button** to cancel any running operation
- No CMD popup window during batch operations
- Auto-detection of `lucksystem.exe`

> Place `LuckSystemGUI.exe` in the same folder as `lucksystem.exe` to use.

---

## Patches

### Version 3 — Patch 3 *(latest)*

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
- **Little Busters English** — CZ4 confirmed

---

## Credits

- **[wetor](https://github.com/wetor)** — LuckSystem original
- **masagrator** — RawSize bug identification (CZ3 layers)
- **[G2-Games](https://github.com/G2-Games)** — CZ4 reference ([lbee-utils](https://github.com/G2-Games/lbee-utils))
- **Yoremi** — patches 1-14, AIR French translation, GUI
