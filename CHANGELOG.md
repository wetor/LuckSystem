# LuckSystem — Patches Yoremi

Fork de [LuckSystem 2.3.2](https://github.com/wetor/LuckSystem) avec corrections et ajouts pour le support de la traduction de visual novels utilisant le moteur ProtoDB/LUCA System (AIR, CLANNAD, Kanon, Summer Pockets, Harmonia, etc.).

## Corrections apportées

### Patch 1 — Import de scripts à longueur variable
**Fichier :** `script/script.go`

L'import de scripts traduits échouait avec un panic quand la traduction avait une longueur différente de l'original. Le code vérifiait strictement `len(paramList) == len(code.Params)`, ce qui bloquait toute traduction plus longue ou plus courte.

- Suppression de la vérification stricte du nombre de paramètres
- Ajout de bounds checking dans la boucle de conversion et le merge des paramètres
- Les offsets de jump (GOTO, IFN, IFY…) sont recalculés automatiquement

### Patch 2 — Correction du pipeline CZ3 (export/import PNG)
**Fichiers :** `czimage/cz3.go`, `czimage/imagefix.go`

L'export et l'import de CZ3 corrompaient silencieusement les données pixels.

- **Magic byte** : `Write()` écrasait le magic "CZ3" → "CZ0", rendant le fichier illisible par le jeu
- **Format NRGBA** : Conversion automatique de tout format PNG en NRGBA 32-bit avant encodage
- **Buffer aliasing** : `DiffLine()` et `LineDiff()` partageaient des slices au lieu de copier, provoquant une corruption des données delta

### Patch 3 — Corruption mémoire dans le décompresseur LZW
**Fichier :** `czimage/lzw.go`

Le décompresseur LZW ajoutait des entrées dictionnaire qui référençaient directement le slice `w` au lieu d'en faire une copie. Les anciennes entrées du dictionnaire pointaient vers des données corrompues.

- Allocation explicite de `newEntry` avec copie de `w` avant ajout au dictionnaire

### Patch 4 — RawSize incorrect dans la table de blocs CZ
**Fichier :** `czimage/util.go`

Bug critique causant la corruption visuelle des CG en jeu (artefacts colorés). Les fonctions `Compress()` et `Compress2()` calculaient un `RawSize` erroné pour chaque bloc LZW.

1. **Carry-over non compensé** : Le dernier élément LZW reporté au bloc suivant n'était pas déduit du compteur de bytes.
2. **Encodage UTF-8 de Go** : `len(string(byte(200)))` retourne 2 au lieu de 1 pour les octets > 127, causant des erreurs ±1 sur les RawSize.

### Patch 5 — Support du format CZ4
**Fichiers :** `czimage/cz4.go` (nouveau), `czimage/imagefix.go`, `czimage/cz.go`

Ajout du décodage et de l'encodage du format CZ4, utilisé dans les jeux récents (Little Busters English, LOOPERS, Harmonia, Kanon 2024).

Le CZ4 diffère du CZ3 par le stockage séparé des canaux RGB (w×h×3) et Alpha (w×h), chacun avec son propre delta line encoding indépendant. Le LZW et le calcul de blockHeight sont identiques au CZ3.

### Patch 6 — Padding d'alignement dans pak.go
**Fichier :** `pak/pak.go`

Après l'écriture d'un PAK reconstruit (quand les fichiers remplacés sont plus grands que les originaux), le fichier n'était pas aligné sur la taille de bloc, ce qui pouvait causer des erreurs de lecture.

- Ajout de padding zéro en fin de fichier pour aligner sur `BlockSize`

### Patch 7 — Correction AIR.py (résolution du module base)
**Fichier :** `data/AIR.py`

Le script de définition AIR.py utilisait `from base.air import *` pour importer les fonctions de `data/base/air.py`. Cet import échouait systématiquement en mode `script import` car le working directory de LuckSystem n'est pas `data/`, empêchant Python de résoudre le chemin `base/air`.

Erreur reproduite avec la commande documentée dans usage.md :
```
FileNotFoundError: 'Failed to resolve "base/air"'
panic: runtime error: invalid memory address or nil pointer dereference
```

- Fusion des fonctions de `base/air.py` directement dans `AIR.py` (IFN, IFY, FARCALL, GOTO, GOSUB, JUMP, etc.)
- Ajout de la fonction `ONGOTO` qui était absente
- Suppression de la dépendance `from base.air import *`

## Jeux testés

- AIR (Steam) — traduction française complète, SYSCG.pak 51/51 (CZ3+CZ4), SCRIPT.pak import/export
- Summer Pockets — fix RawSize confirmé (rapport masagrator)

## Crédits

- **wetor** — LuckSystem original
- **masagrator** — identification du bug RawSize (CZ3 layers)
- **Yoremi** — patches 1-7, traduction française d'AIR
- **G2-Games** — référence CZ4 (lbee-utils)
