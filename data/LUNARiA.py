import core
from base.lunaria import *

def Init():
    core.set_config(
        expr_charset=core.Charset_UTF8,
        text_charset=core.Charset_Unicode,
        default_export=True
    )

def MESSAGE():
    core.read_uint16(True)                    # voice_id
    txt = core.read_len_str(core.text)        # msg_jp_len, msg_jp_str
    if len(txt) > 0:                          # for messages without text
        core.read_len_str(core.Charset_UTF8)  # msg_en_len, msg_en_str
        core.read_len_str(core.text)          # msg_cn_len, msg_cn_str
    if core.can_read():
        core.read(False)
    core.end()

def LOG_BEGIN():
    core.read_uint8(False)
    core.read_uint8(False)
    core.read_uint8(False)
    core.read_len_str(core.text)  # msg_jp_len, msg_jp_str
    core.read_len_str(core.text)  # msg_en_len, msg_en_str
    core.read_len_str(core.text)  # msg_cn_len, msg_cn_str
    core.read(False)
    core.end()

def MOVIE():
    core.read_len_str(core.Charset_UTF8)
    core.read(False)
    core.end()

def VARSTR_SET():
    core.read_len_str(core.Charset_UTF8)     # jp_len, jp_str
    core.read_uint16(True)                   # var_id
    core.read_len_str(core.Charset_Unicode)  # en_len, en_str
    core.read_len_str(core.Charset_Unicode)  # en_len2, en_str2
    core.read_len_str(core.Charset_Unicode)  # en_len3, en_str3
    core.end()

def DIALOG():
    core.read_uint16(False)
    core.read_uint16(False)
    core.read_len_str(core.text)
    core.read(False)
    core.end()