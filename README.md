# LUCA System Tools

## LUCA Script



## 脚本结构

| 类型  |  长度 | 名称 |  说明 | 
|------|-------|-----|-----|
| uint16 |  2  | len | 代码长度|
| uint8 |  1  | opcode | 指令索引|
| uint8 |  1  | flag | 一个标志，值0~3|
| []uint16 |  2 * n  | data0 | 未知参数，其中n=flag(flag<3),n=2(flag==3)|
| params |  len -4 -2*n  | params | 参数|
| uint8 |  k  | align | 补齐位，其中k=len%2|


uint16 2 byte 整条语句的长度  
uint8 1 byte

## 更新日志

### 6.30
- 支持多游戏
- 设计参数、函数等结构

### 6.28
- 框架设计与编写
- 第三方包的选择与测试
- 支持LB_EN基本解析