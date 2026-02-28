# LuckSystem GUI (Linux)

Graphical interface for [LuckSystem](https://github.com/wetor/LuckSystem), the Visual Art's/Key visual novel translation toolkit.

Interface graphique pour [LuckSystem](https://github.com/wetor/LuckSystem), l'outil de traduction de visual novels Visual Art's/Key.

## Architecture

The GUI is a **standalone wrapper** — it does NOT embed LuckSystem source code. It calls `lucksystem` via subprocess, exactly like you would from a terminal.

```
LuckSystemGUI      ←→  lucksystem (subprocess)
   (Wails/Go)              (CLI tool)
```

This design follows [wetor's recommendation](https://github.com/wetor/LuckSystem) to keep the GUI separated from the core tool for cross-platform compatibility and maintainability.

## Setup

1. Download the Linux binary `lucksystem` from [LuckSystem releases](https://github.com/wetor/LuckSystem/releases) (or build from the [Yoremi fork](https://github.com/yoremi-trad-fr/LuckSystem-2.3.2-Yoremi-Update))
2. Place `lucksystem` in the same directory as `LuckSystemGUI`
3. Make sure both files are executable:

```bash
chmod +x lucksystem
chmod +x LuckSystemGUI
```

4. Run:

```bash
./LuckSystemGUI
```

The GUI auto-detects `lucksystem` in the same directory, current working directory, or system PATH. You can also manually locate it by clicking the path indicator in the title bar.

## Features

| Operation | Description |
|-----------|-------------|
| **Script Decompile** | Extract scripts from SCRIPT.PAK to text files |
| **Script Compile** | Repack translated scripts into a new SCRIPT.PAK |
| **PAK Extract** | Extract all files from any .PAK archive |
| **PAK Replace** | Replace files inside a .PAK archive |
| **Font Extract** | Export CZ font atlas to PNG + charset list |
| **Font Edit** | Redraw/append characters using a TTF font |
| **Image Export** | Convert CZ images to PNG (single or batch) |
| **Image Import** | Convert PNG back to CZ format (single or batch) |
| **Dialogue Extract** | Extract translatable dialogue from decompiled scripts to TSV (single file or batch) |
| **Dialogue Import** | Reimport translated dialogue from TSV back into scripts (single file or batch) |

### Dialogue Extract / Import

The Dialogue Extract and Import functions provide a streamlined translation workflow based on TSV files, replacing manual script editing.

**Extract** scans decompiled script files (`.txt`) for translatable lines (`MESSAGE` and `LOG_BEGIN` entries) and exports them to tab-separated `.tsv` files. The language columns are numbered (Lang 1, Lang 2, Lang 3, Lang 4) rather than named, since the order of languages varies between games. You select which columns to extract via checkboxes.

**Import** reads a translated `.tsv` file and reinjects the text back into the corresponding decompiled script. You select which column number contains the target language. Matching is done by sequential ID for robustness.

Both operations support single-file and batch modes. The format auto-detection scans the script to determine the number of available language columns.

**TSV format example:**
```
ID	TAG	Lang 2
1	MESSAGE	`Rin@❝Stop bullying the weak!❞
2	MESSAGE	`Riki@❝Masato... where are you going?❞
3	LOG_BEGIN	Chapter 1
```

## Supported games

All games using the ProtoDB / LUCA System engine:
AIR, CLANNAD, Kanon, Little Busters, Summer Pockets, Harmonia, LOOPERS, LUNARiA, Planetarian, etc.

## Build from source (Linux)

Requires: [Go 1.23+](https://go.dev/), [Node.js](https://nodejs.org/), [Wails CLI](https://wails.io/)

```bash
cd frontend && npm install && cd ..
go mod tidy
wails dev          # Development with hot-reload
wails build        # Build to build/bin/LuckSystemGUI
```

> ⚠️ **Do NOT run `npm audit fix --force`** — it upgrades Svelte/Vite to incompatible major versions.

## Notes

- No `.exe` extension on Linux
- If the app does not launch, check execution permissions
- Wayland users may need XWayland depending on desktop environment

## Credits

- **[wetor](https://github.com/wetor)** — [LuckSystem](https://github.com/wetor/LuckSystem) core CLI tool
- **Yoremi** — GUI development, [Yoremi fork](https://github.com/yoremi-trad-fr/LuckSystem-2.3.2-Yoremi-Update) patches

## License

MIT
