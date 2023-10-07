## 使用Plugin支持任意游戏的脚本反编译与导入

插件为非标准的Python，语法类似Python3.4，缺少大量的内置库和一些特性，基本使用没有问题

### 编写规则
可参考 [LOOPERS.py](data/LOOPERS.py) 和 [SP.py](data/SP.py)  

#### 解析操作函数
在脚本种定义与OPCODE同名的函数（大小写一致），即可在反编译、导入时使用脚本中的函数去解析对应的操作参数  

#### Init函数
在脚本中定义的`Init`函数会在最开始执行一次，一般用来调用`set_config`，也可以进行其他初始化操作  

#### opcode_dict全局变量
在没有载入OPCODE或者存在未知OPCODE时，反编译会将操作名输出为`0x22` `0x3A`等十六进制标记，此时可以通过`opcode_dict`来手动映射十六进制标记和操作函数名，如下：  
```python
opcode_dict = {
    '0x22': 'MESSAGE',
    '0x25': 'SELECT'
}
```
此时对于未知的`0x22`操作，将会使用插件中的`MESSAGE`函数进行解析

### core包内置常量
```python
core.Charset_UTF8 = "UTF-8"
core.Charset_Unicode = "UTF-16LE"
core.Charset_SJIS = "Shift_JIS"
```
三种字符集的特点如下：
#### Charset_UTF8
1~3byte一字，固定的结尾0x00  
通常英文字符串使用此编码，目前如LOOPERS等新游戏也作为表达式(expr)和PAK包文件名编码也是UTF8

#### Charset_Unicode
2byte(一个uint16)一字，固定的结尾0x0000  
目前绝大多数的中文、日文文本均是此编码

#### Charset_SJIS
1~2byte一字，固定的结尾0x00
老游戏的表达式和PAK包编码，一般为单日语游戏使用  

### core包内置函数

#### set_config
`set_config(expr_charset, text_charset, default_export=True)`  
设置默认值  
设置`expr_charset`后可通过`core.expr`获取对应设置的值，通常与PAK的文件名编码一致    
设置`text_charset`后可通过`core.text`获取对应设置的值，也会成功`read_str`和`read_len_str`的默认字符集值  
`default_export`即为解析插件未定义OPCODE时，是否将自动解析的uint16参数导出

#### read
`read(export=False) -> list(int)`  
以uint16的形式读取**所有**参数，若剩余一位则将以uint8形式读取

#### read_uint8
`read_uint8(export=False) -> int`  
以uint8的形式读取一个参数

#### read_uint16
`read_uint16(export=False) -> int`  
以uint16的形式读取一个参数

#### read_uint32
`read_uint32(export=False) -> int`  
以uint32的形式读取一个参数

#### read_str
`read_str(charset=textCharset, export=True) -> str`  
按指定编码读取一个字符串  
注意：如确定字符串前面存在一个`uint16`为字符串长度，则需要使用`read_len_str`进行读取，否则长度错乱会导致导入出错

#### read_len_str
`read_len_str(charset=textCharset, export=True) -> str`  
按指定编码读取一个包含长度的字符串，导入时会自动计算新的长度 

#### read_jump (重要)
`read_jump(file='', export=True) -> int`  
以uint32的形式读取一个跳转位置参数，导入时将会自动重构跳转位置  
跨文件的跳转需要传入`file`参数，一般为前一个参数。跨文件跳转标记为`global233`   
跨文件的`file`参数，需要使用read_str或read_len_str读取，字符集为expr，通常与PAK文件名编码一致  
文件内跳转的标记为`label66`

#### end (重要)
`end()`  
将已经进行的读取提交，必须调用，否则将无法正常导出导入  

#### can_read
`can_read() -> bool`  
判断接下来是否可以进行读取

