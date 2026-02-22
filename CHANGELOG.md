# LuckSystem — Yoremi Fork — CHANGELOG

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
