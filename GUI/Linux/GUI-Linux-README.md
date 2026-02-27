LuckSystem GUI (Linux)

Graphical interface for LuckSystem
, the Visual Art's/Key visual novel translation toolkit.

Interface graphique pour LuckSystem
, l'outil de traduction de visual novels Visual Art's/Key.

Architecture

The GUI is a standalone wrapper — it does NOT embed LuckSystem source code.
It calls lucksystem via subprocess, exactly like from a terminal.

LuckSystemGUI      ←→  lucksystem (subprocess)
   (Wails/Go)              (CLI tool)

This design follows wetor's recommendation to keep the GUI separated from the core tool for cross-platform compatibility and maintainability.

Setup

Download the Linux binary lucksystem from LuckSystem releases
(or build from the Yoremi fork)

Place lucksystem in the same directory as LuckSystemGUI

Make sure both files are executable:

chmod +x lucksystem
chmod +x LuckSystemGUI

Run:

./LuckSystemGUI

The GUI auto-detects lucksystem in:

the same directory

the current working directory

the system PATH

You can also manually locate it by clicking the path indicator in the title bar.

Features
Operation	Description
Script Decompile	Extract scripts from SCRIPT.PAK to text files
Script Compile	Repack translated scripts into a new SCRIPT.PAK
PAK Extract	Extract all files from any .PAK archive
PAK Replace	Replace files inside a .PAK archive
Font Extract	Export CZ font atlas to PNG + charset list
Font Edit	Redraw/append characters using a TTF font
Image Export	Convert CZ images to PNG (single or batch)
Image Import	Convert PNG back to CZ format (single or batch)
Supported games

All games using the ProtoDB / LUCA System engine:

AIR, CLANNAD, Kanon, Little Busters, Summer Pockets, Harmonia, LOOPERS, LUNARiA, Planetarian, etc.

Build from source (Linux)

Requires:

Go 1.23+

Node.js

Wails CLI

cd frontend && npm install && cd ..
go mod tidy
wails dev          # Development with hot-reload
wails build        # Build to build/bin/LuckSystemGUI

If you are cross-compiling from Windows:

set GOOS=linux
set GOARCH=amd64
wails build
Notes

No .exe extension on Linux

If the app does not launch, check execution permissions

Wayland users may need XWayland depending on desktop environment

Credits

wetor — LuckSystem core CLI tool
Yoremi — GUI development, Yoremi fork patches

License

MIT
