import core
from base.cartagrahd import *

def Init():
    # core.Charset_UTF8
    # core.Charset_Unicode
    # core.Charset_SJIS
    core.set_config(expr_charset=core.Charset_UTF8, # 表达式编码, core.expr
                    text_charset=core.Charset_Unicode, # 文本编码, core.text
                    default_export=True) # 未定义指令的参数是否全部导出

def MESSAGE():
    core.read_uint16(True)
    core.read_len_str(core.text)
    core.read_uint8()
    core.read(False)
    core.end()

def SELECT():
    core.read_uint16()
    core.read_uint16()
    core.read_uint16()
    core.read_uint16()
    core.read_len_str(core.text)
    core.read(True)
    core.end()

def DIALOG():
    core.read_uint16(False)
    core.read_uint16(False)
    core.read_len_str(core.text)
    core.read(False)
    core.end()

def LOG_BEGIN():
    core.read_uint8(False)
    core.read_uint8(False)
    core.read_uint8(False)
    core.read_len_str(core.text)
    core.read(False)
    core.end()