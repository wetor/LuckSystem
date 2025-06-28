import core
from base.butterflyspoisonbc import *

def Init():
    # core.Charset_UTF8
    # core.Charset_Unicode
    # core.Charset_SJIS
    core.set_config(
        expr_charset=core.Charset_UTF8,    # core.expr
        text_charset=core.Charset_Unicode, # core.text
        default_export=True
    )

def MESSAGE():
    core.read_uint16(False)
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.expr)    # en_len, en_str
    if core.can_read():
        core.read(False)
    core.end()

def MOVIE():
    core.read_uint16(False)
    core.read_str(core.expr)            # filename
    core.read(False)
    core.end()

def VARSTR_SET():
    core.read_uint16(False)
    core.read_len_str(core.text)
    if core.can_read():
        core.read(False)
    core.end()

def DIALOG():
    core.read_uint16(False)
    core.read_uint16(False)
    core.read_len_str(core.text)        # jp_len, jp_str
    core.read(False)
    core.end()

def SELECT():
    core.read_uint16()
    core.read_uint16()
    core.read_uint16(False)
    core.read_uint16(False)
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.text)    # en_len, en_str
    if core.can_read():
        core.read(False)
    core.end()