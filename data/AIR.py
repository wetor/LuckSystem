import core

# ============================================================================
# CONFIGURATION
# ============================================================================

def Init():
    # Configuration des encodages pour AIR Steam
    core.set_config(
        expr_charset=core.Charset_UTF8,      # Expression encoding
        text_charset=core.Charset_Unicode,   # Text encoding (UTF-16)
        default_export=True                  # Export undefined operations
    )

# ============================================================================
# JUMP OPERATIONS (anciennement dans base/air.py)
# ============================================================================

def IFN():
    """Conditional jump if false"""
    core.read_len_str(core.expr)
    core.read_jump()
    core.end()

def IFY():
    """Conditional jump if true"""
    core.read_len_str(core.expr)
    core.read_jump()
    core.end()

def FARCALL():
    """Call to another script file"""
    core.read_uint16(True)
    file = core.read_len_str(core.expr)
    core.read_jump(file)
    core.end()

def GOTO():
    """Unconditional jump"""
    core.read_jump()
    core.end()

def ONGOTO():
    """Multi-way jump (switch/case)"""
    core.read_len_str(core.expr)  # Expression to evaluate
    count = core.read_uint16(True)  # Number of jump targets
    for i in range(count):
        core.read_jump()
    core.end()

def GOSUB():
    """Call subroutine"""
    core.read_uint16(True)
    core.read_jump()
    core.end()

def JUMP():
    """Jump to file"""
    file = core.read_len_str(core.expr)
    if core.can_read():
        core.read_jump(file)
    core.end()

def JUMPPOINT():
    """Jump target marker"""
    core.end()

def RETURN():
    """Return from subroutine"""
    core.end()

def FARRETURN():
    """Return from far call"""
    core.end()

# ============================================================================
# TEXT OPERATIONS (anciennement dans data/AIR.py)
# ============================================================================

def MESSAGE():
    """Display message with Japanese, English, and Chinese text"""
    core.read_uint16(False)
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.expr)    # en_len, en_str
        core.read_len_str(core.text)    # zh_len, zh_str
    if core.can_read():
        core.read(False)
    core.end()

def MOVIE():
    """Play movie file"""
    core.read_uint16(False)
    core.read_str(core.expr)            # file_name
    core.read(False)
    core.end()

def VARSTR_SET():
    """Set variable string"""
    core.read_uint16(False)             # const value
    core.read_str(core.expr)            # filename_str
    core.read_uint16(True)              # varstr_id
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.text)    # en_len, en_str
        core.read_len_str(core.text)    # zh_len, zh_str
    if core.can_read():
        core.read(False)
    core.end()

def DIALOG():
    """Display dialog box"""
    core.read_uint16(False)
    core.read_uint16(False)
    core.read_len_str(core.text)        # jp_len, jp_str
    core.read(True)
    core.end()

def LOG_BEGIN():
    """Begin log entry"""
    core.read_uint8(False)
    core.read_uint8(False)
    core.read_uint8(False)
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.text)    # en_len, en_str
        core.read_len_str(core.text)    # zh_len, zh_str
    if core.can_read():
        core.read(False)
    core.end()

def SELECT():
    """Display selection menu"""
    core.read_uint16()
    core.read_uint16()
    core.read_uint16(False)
    core.read_uint16(False)
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.text)    # en_len, en_str
        core.read_len_str(core.text)    # zh_len, zh_str
    if core.can_read():
        core.read(False)
    core.end()
