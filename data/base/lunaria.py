import core

def GOTO():
    # GOTO ({jump})
    core.read_jump()
    core.end()

def FARCALL():
    # FARCALL (index, file_str, {jump})
    core.read_uint16(True)
    file = core.read_len_str(core.expr)
    core.read_jump(file)
    core.end()

def IFN():
    # IFN (expr_len, expr_str, {jump})
    core.read_len_str(core.expr)
    core.read_jump()
    core.end()