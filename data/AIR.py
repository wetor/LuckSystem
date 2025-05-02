import core
from base.air import *

def Init():
    # core.Charset_UTF8
    # core.Charset_Unicode
    # core.Charset_SJIS
    core.set_config(expr_charset=core.Charset_UTF8, # 表达式编码, core.expr
                    text_charset=core.Charset_Unicode, # 文本编码, core.text
                    default_export=True) # 未定义指令的参数是否全部导出

def MESSAGE():
    core.read_uint16(False)
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.expr)    # en_len, en_str
        core.read_len_str(core.text)    # zh_len, zh_len
    if core.can_read():
        core.read(False)
    core.end()

def MOVIE():
    core.read_uint16(False)
    core.read_str(core.expr)            # file_name
    core.read(False)
    core.end()

def VARSTR_SET():
    core.read_uint16(False)             # const value
    core.read_str(core.expr)            # filename_str
    core.read_uint16(True)              # varstr_id
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.text)    # en_len, en_str
        core.read_len_str(core.text)    # zh_len, zh_len
    if core.can_read():
        core.read(False)
    core.end()

def DIALOG():
    core.read_uint16(False)
    core.read_uint16(False)
    core.read_len_str(core.text)        # jp_len, jp_str
    core.read(True)
    core.end()

def LOG_BEGIN():
    core.read_uint8(False)
    core.read_uint8(False)
    core.read_uint8(False)
    txt = core.read_len_str(core.text)  # jp_len, jp_str
    if len(txt) > 0:
        core.read_len_str(core.text)    # en_len, en_str
        core.read_len_str(core.text)    # zh_len, zh_len
    if core.can_read():
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
        core.read_len_str(core.text)    # zh_len, zh_len
    if core.can_read():
        core.read(False)
    core.end()
