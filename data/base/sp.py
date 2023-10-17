import core

def IFN():
    # IFN (expr_str, {jump})
    core.read_str(core.expr)
    core.read_jump()
    core.end()

def IFY():
    # IFY (expr_str, {jump})
    core.read_str(core.expr)
    core.read_jump()
    core.end()

def FARCALL():
    # FARCALL (index, file_str, {jump})
    core.read_uint16(True)
    file = core.read_str(core.expr)
    core.read_jump(file)
    core.end()

def GOTO():
    # GOTO ({jump})
    core.read_jump()
    core.end()

def GOSUB():
    # GOTO (int, {jump})
    core.read_uint16(True)
    core.read_jump()
    core.end()

def JUMP():
    # JUMP (file_str, {jump})
    file = core.read_str(core.expr)
    if core.can_read():
        core.read_jump(file)
    core.end()
