import core
from base.sp import *

# 在没有载入OPCODE时，进行手动映射以导出脚本
# 此方式导出的脚本修改后导入将会无法使用！！
opcode_dict = {
    '0x22': 'MESSAGE',
    '0x25': 'SELECT'
}

def Init():
    # core.Charset_UTF8
    # core.Charset_Unicode
    # core.Charset_SJIS
    core.set_config(expr_charset=core.Charset_SJIS, # 表达式编码, core.expr
                    text_charset=core.Charset_Unicode, # 文本编码, core.text
                    default_export=True) # 未定义指令的参数是否全部导出

def MESSAGE():
    core.read_uint16(True)
    core.read_len_str(core.text)
    core.read_uint8()
    core.end()

def SELECT():
    core.read_uint16()
    core.read_uint16()
    core.read_uint16(False)
    core.read_uint16(False)
    core.read_len_str(core.text)
    core.read(False)
    core.end()
