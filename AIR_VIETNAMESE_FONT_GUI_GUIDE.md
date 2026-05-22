# AIR Vietnamese Font Procedure - GUI Guide

This guide explains how to test a Vietnamese-capable TTF with AIR Steam fonts using the LuckSystem Yoremi GUI.

It is written for GUI use only. You only need `LuckSystemGUI.exe`, the matching `lucksystem.exe`, the AIR font PAK files, a charset text file, and a TTF.

The tested target is the English font slot. For quick in-game checks, generating only `FONT_GOTHIC1.PAK` plus the matching `FONT__INFO.PAK` is enough.

## Recommended Settings

Use these settings for AIR Vietnamese tests:

| GUI field | Value |
|----------|-------|
| Font family | `GOTHIC1` |
| Slot | English / Japanese font slot, not Chinese `ZC` |
| Mode | `Insert at index` |
| Index | `7091` |
| Charset file | Missing-only Vietnamese charset, 102 characters |
| TTF | The Vietnamese-capable TTF you want to test |
| Files to copy back to AIR | `FONT_GOTHIC1.PAK` and `FONT__INFO.PAK` |

Do not use `Append to end` for AIR Steam Vietnamese tests.

## Why These Settings

AIR already contains some Vietnamese-compatible accented Latin characters. If they are replaced, the result can look worse than the original. The safer method is:

1. Keep the characters already present in AIR unchanged.
2. Inject only the missing Vietnamese characters.
3. Replace unused cells at the end of the original font table.

For the tested English/Japanese AIR font table:

- Total glyph cells: `7193`
- Characters already usable: `32`
- Missing Vietnamese characters to inject: `102`
- Insert start index: `7193 - 102 = 7091`

The current GUI does not expose vertical offset adjustment. GUI-only tests may still show a small baseline difference depending on the TTF. If the font is readable but sits too high or too low, test another TTF or ask the maintainer for a prepared build.

## Charset Files

### Use This Charset For GUI Insert Mode

Create a UTF-8 text file containing exactly these 102 missing Vietnamese characters:

```text
ăđơưĂĐƠƯảạắằẳẵặấầẩẫậẻẽẹếềểễệỉĩịỏọốồổỗộớờởỡợủũụứừửữựỳỷỹỵẢẠẮẰẲẴẶẤẦẨẪẬẺẼẸẾỀỂỄỆỈĨỊỎỌỐỒỔỖỘỚỜỞỠỢỦŨỤỨỪỬỮỰỲỶỸỴ
```

The same charset is also included in the repository as:

```text
examples/AIR_vietnamese_missing_102.txt
```

This is the charset to use in the GUI with:

- Mode: `Insert at index`
- Index: `7091`

### Do Not Re-inject These Existing Characters

These 32 characters are already present in AIR and should stay at their original indexes:

```text
âêôÂÊÔáàãéèíìóòõúùýÁÀÃÉÈÍÌÓÒÕÚÙÝ
```

## Step 1 - Back Up The Original Files

Before editing, keep a safe copy of the original AIR files:

```text
AIR/files/font_win32_1280/FONT__INFO.PAK
AIR/files/font_win32_1280/FONT_GOTHIC1.PAK
```

Work in a separate folder and copy the finished files back only when you are ready to test.

## Step 2 - Extract The Font PAKs

Open the GUI and use:

```text
PAK (Font) -> Font Extract
```

Extract both files:

```text
FONT__INFO.PAK
FONT_GOTHIC1.PAK
```

Use `UTF-8` as the charset value for extraction and replacement.

Recommended output layout:

```text
work/
  INFO/
  GOTHIC1/
  FONT__INFO_list.txt
  FONT_GOTHIC1_list.txt
```

The extracted font files usually have no extension. Keep their names exactly as extracted.

After extraction, make a backup copy of the whole `work` folder. The easiest GUI workflow is to edit the extracted working files in place, because the generated `_list.txt` files already point to that folder.

## Step 3 - Preview One Font Size

Use:

```text
FONT -> Font Extract
```

Select one matching pair:

```text
Source CZ file:   work/GOTHIC1/<GOTHIC1 size file>
Source info file: work/INFO/<matching info size file>
Output PNG:       work/preview.png
Output charset:   work/original_charset.txt
```

The `CZ` file and the `info` file must have the same size number. For example, if you choose the Gothic size `28`, use the matching `info28`.

This step is only for checking the original font and confirming that the pair is correct.

## Step 4 - Edit One Size With Another TTF

Use:

```text
FONT -> Font Edit
```

Fill the fields like this:

```text
Source CZ file:    work/GOTHIC1/<original size file>
Source info file:  work/INFO/<matching info size file>
TTF font file:     your_test_font.ttf
Mode:              Insert at index
Index:             7091
Charset file:      missing_vietnamese_102.txt
Output CZ:         work/GOTHIC1/<same size file name>
Output info:       work/INFO/<same info file name>
```

Important notes:

- The output paths should normally use the same filenames as the original extracted entries.
- In this simple workflow, the output overwrites the extracted working copy, not the original game PAK.
- The GUI output fields expect full paths without adding an extra extension.
- For a quick TTF check, patch one size first.
- For a real in-game test using only the GUI, repeat the same edit for every `GOTHIC1` size used by AIR.
- If you prefer a separate `work_patched` folder, copy the entire extracted folder first and update the copied `_list.txt` files so their paths point to `work_patched`.

Known AIR font sizes in this family include:

```text
12, 16, 22, 24, 25, 27, 28, 29, 30, 32, 33, 35, 38
```

## Step 5 - Rebuild The PAK Files

After editing the font size files, rebuild the two PAKs with:

```text
PAK (Font) -> Font Replace
```

Rebuild:

```text
FONT_GOTHIC1.PAK
FONT__INFO.PAK
```

Use the list files created during extraction when possible:

```text
FONT_GOTHIC1_list.txt
FONT__INFO_list.txt
```

List mode is recommended because it preserves the original entry order and names.

Do not rebuild from a folder that contains only the edited size files. The rebuild input must still point to a complete extracted font set, with the edited entries replacing the originals.

## Step 6 - Test In AIR

Copy only these rebuilt files into the English font folder used by the game:

```text
AIR/files/font_win32_1280/FONT_GOTHIC1.PAK
AIR/files/font_win32_1280/FONT__INFO.PAK
```

Start the game and check Vietnamese dialogue in the English slot.

If the characters appear but sit slightly too high or too low, the TTF works but needs a vertical offset adjustment. The GUI cannot currently change this offset.

## Which Font Edit Mode Should Be Used?

### Insert At Index - Recommended For AIR

Use this mode for AIR Vietnamese font tests.

Recommended values:

```text
Mode:  Insert at index
Index: 7091
Charset: missing Vietnamese charset, 102 characters
```

This replaces cells near the end of the font table that are not needed for the English Vietnamese script. It does not grow the font table, which is important because AIR is sensitive to the original font structure.

### Append To End - Not Recommended For AIR

Do not use this mode for AIR Vietnamese tests.

`Append to end` adds new characters after the existing table. This changes the font count and layout. AIR Steam may fail to start, ignore the new glyphs, or produce visual bugs depending on the rebuilt PAK and font table.

### Redraw All - Not Recommended For This Test

Do not use this mode for normal Vietnamese injection.

`Redraw all` redraws the whole font with the selected TTF. It is not an insertion mode: if a charset is provided, it starts remapping from the beginning of the table instead of using index `7091`; if no charset is provided, it redraws the existing font table. This is only useful if you intentionally want to replace the entire font style, and it can change many glyphs that already looked correct.

## Testing Another TTF

When trying a different TTF:

1. Confirm that the TTF supports Vietnamese combining and precomposed characters.
2. Use the same missing-only 102-character charset.
3. Use `Insert at index`.
4. Use index `7091`.
5. Patch `GOTHIC1` first.
6. Test in game before patching more families.

Good signs:

- No crash at game startup.
- No visual corruption in menus.
- Missing Vietnamese characters now appear.
- Existing accented characters such as `ó`, `à`, `ê`, `ô` still look stable.

Bad signs:

- Game does not start.
- Menu fonts become corrupted.
- Existing accented characters look different or broken.
- Vietnamese marks are clipped or too high/low.

If only the vertical position is wrong, the TTF may still be usable, but this version of the GUI cannot adjust that value directly.

## Quick Reference

For GUI testing, use:

```text
FONT -> Font Edit
Mode: Insert at index
Index: 7091
Charset: missing_vietnamese_102.txt
Family: GOTHIC1
Output PAKs: FONT_GOTHIC1.PAK + FONT__INFO.PAK
```

Avoid:

```text
Append to end
Redraw all
Full 134-character charset in GUI insert mode
Chinese ZC slot for the current English-slot test
```
