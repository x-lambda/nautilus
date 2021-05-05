# logging

`golang`常用的开源日志库
* [`logrus`](https://github.com/sirupsen/logrus)
* [`zerolog`](https://github.com/rs/zerolog)
* [`zap`](https://github.com/uber-go/zap)

3个库各有千秋，`logrus`用的人多，`zerolog`主打高性能，低消耗，`zap`Uber出品，自带大厂光环。

设计目标
* 支持现有`ELK`方案
* 兼容`ctx`，可以通过`ctx`获取`trace id`
* 不输出和代码相关，所在行数等信息
* 使用简单，只接受输出信息`msg`一个参数

配置选择
* `LOG_LEVEL`: 日志等级
* `LOG_AGENT`: ELK配置

默认配置
* 输出到标准输出
* `json`格式
* `time`格式：`2006-01-02 15:04:05.xxxxxx`