# LuckSystem GUI (Windows)

Graphical interface for [LuckSystem](https://github.com/wetor/LuckSystem), the Visual Art's/Key visual novel translation toolkit.

Interface graphique pour [LuckSystem](https://github.com/wetor/LuckSystem), l'outil de traduction de visual novels Visual Art's/Key.

![LuckSystem GUI](screenshot.png)

## Architecture

The GUI is a **standalone wrapper** — it does NOT embed LuckSystem source code. It calls `lucksystem.exe` via subprocess, exactly like you would from a terminal.

```
LuckSystemGUI.exe  ←→  lucksystem.exe (subprocess)
   (Wails/Go)              (CLI tool)
```

This design follows [wetor's recommendation](https://github.com/wetor/LuckSystem) to keep the GUI separated from the core tool for cross-platform compatibility and maintainability.

## Setup

1. Download `lucksystem.exe` from [LuckSystem releases](https://github.com/wetor/LuckSystem/releases) (or build from the [Yoremi fork](https://github.com/yoremi-trad-fr/LuckSystem-2.3.2-Yoremi-Update))
2. Place `lucksystem.exe` next to `LuckSystemGUI.exe`
3. Run `LuckSystemGUI.exe`

The GUI auto-detects `lucksystem.exe` in the same directory, current working directory, or system PATH. You can also manually locate it by clicking the path indicator in the title bar.

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

## Supported games

All games using the ProtoDB / LUCA System engine:
AIR, CLANNAD, Kanon, Little Busters, Summer Pockets, Harmonia, LOOPERS, LUNARiA, Planetarian, etc.

## Build from source

Requires: [Go 1.23+](https://go.dev/), [Node.js](https://nodejs.org/), [Wails CLI](https://wails.io/)

```bash
cd frontend && npm install && cd ..
go mod tidy
wails dev          # Development with hot-reload
wails build        # Build to build/bin/LuckSystemGUI.exe
```

> ⚠️ **Do NOT run `npm audit fix --force`** — it upgrades Svelte/Vite to incompatible major versions.

## Credits

- **[wetor](https://github.com/wetor)** — [LuckSystem](https://github.com/wetor/LuckSystem) core CLI tool
- **Yoremi** — GUI development, [Yoremi fork](https://github.com/yoremi-trad-fr/LuckSystem-2.3.2-Yoremi-Update) patches

## License

MIT
