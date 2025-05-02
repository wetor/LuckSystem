import core
from base.air import *

def Init():
    # core.Charset_UTF8
    # core.Charset_Unicode
    # core.Charset_SJIS
    core.set_config(
        expr_charset=core.Charset_UTF8,
        text_charset=core.Charset_Unicode,
        default_export=True
    )

def MOVIE():
    core.read_len_str(core.Charset_UTF8)
    core.read(False)
    core.end()

def MESSAGE():
    core.read_uint16(True)              # voice_id
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.expr)    # en_len, en_str
        core.read_len_str(core.text)    # fr_len, fr_str
        core.read_len_str(core.text)    # cn1_len, cn1_str
        core.read_len_str(core.text)    # cn2_len, cn2_str
    if core.can_read():
        core.read(False)
    core.end()

def VARSTR_SET():
    core.read_len_str(core.expr)  # filename_len, filename_str
    core.read_uint16(False)        # var_id
    core.read_len_str(core.text)  # jp_len, jp_str
    core.read_len_str(core.text)  # en_len, en_str
    core.read_len_str(core.text)  # fr_len, fr_str
    core.read_len_str(core.text)  # cn1_len, cn1_str
    core.read_len_str(core.text)  # cn2_len, cn2_str
    core.end()

def LOG_BEGIN():
    core.read_uint8(False)
    core.read_uint8(False)
    core.read_uint8(False)
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.text)    # en_len, en_str
        core.read_len_str(core.text)    # fr_len, fr_str
        core.read_len_str(core.text)    # cn1_len, cn1_str
        core.read_len_str(core.text)    # cn2_len, cn2_str
    if core.can_read():
        core.read(False)
    core.end()


def DIALOG():
    core.read_uint16(False)
    core.read_uint16(False)
    core.read_len_str(core.text)        # jp_len, jp_str
    if core.can_read():
        core.read(True)
    core.end()