import core

def IFN():
    # IFN (int, expr_str, {jump})
    core.read_uint16(True)
    core.read_str(core.expr)
    core.read_jump()
    core.end()

def IFY():
    # IFY (int, expr_str, {jump})
    core.read_uint16(True)
    core.read_str(core.expr)
    core.read_jump()
    core.end()

def GOTO():
    # GOTO ({jump})
    core.read_jump()
    core.end()

def JUMP():
    # JUMP (int, file_str, {jump})
    core.read_uint16(True)
    file = core.read_str(core.expr)
    if core.can_read():
        core.read_jump(file)
    core.end()
