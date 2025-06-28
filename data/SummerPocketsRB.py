import core
from base.summerpocketsrb import *

def Init():
    # core.Charset_UTF8
    # core.Charset_Unicode
    # core.Charset_SJIS
    core.set_config(
        expr_charset=core.Charset_UTF8,
        text_charset=core.Charset_Unicode,
        default_export=True
    )

def MESSAGE():
    core.read_uint16(True)              # voice_id
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.expr)    # en_len, en_str
        core.read_len_str(core.text)    # cn_len, cn_str
    if core.can_read():
        core.read(True)
    core.end()

def SELECT():
    core.read_uint16()
    core.read_uint16()
    core.read_uint16(False)
    core.read_uint16(False)
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.text)    # en_len, en_str
        core.read_len_str(core.text)    # zh_len, zh_len
    if core.can_read():
        core.read(False)
    core.end()

def VARSTR_SET():
    core.read_len_str(core.expr)       # filename_len, filename_str
    core.read_uint16(False)            # var_id
    if core.remaining_length() > 6:    # some instructions have 6 zero bytes instead of text
        core.read_len_str(core.text)   # jp_len, jp_str
        core.read_len_str(core.text)   # en_len, en_str
        core.read_len_str(core.text)   # cn_len, cn_str
    if core.can_read():
        core.read(False)
    core.end()

def LOG_BEGIN():
    core.read_uint16(False)
    core.read_uint8(False)
    core.read_len_str(core.text)      # jp_len, jp_str
    core.read_len_str(core.text)      # en_len, en_str

    # special check for the last LOG_BEGIN in script 30_蒼0826
    # for unknown reason, only two lines with en and cn text are passed into it
    if core.remaining_length() > 1:
        core.read_len_str(core.text)  # cn_len, cn_str

    if core.can_read():
        core.read(True)
    core.end()

def DIALOG():
    core.read_uint16(False)
    core.read_uint16(False)
    core.read_len_str(core.text)        # jp_len, jp_str
    if core.can_read():
        core.read(True)
    core.end()

def MOVIE():
    core.read_len_str(core.expr)
    if core.can_read():
        core.read(False)
    core.end()
