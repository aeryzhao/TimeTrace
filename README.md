# 时间记录工具

## 运行

通过 `run.sh` 脚本可以启动、停止、重启时间记录工具。

启动

``` shell
./run.sh start
```
- 检查并清理已被占用的端口（8080 和 5173）。
- 后台启动 Go 后端服务，日志输出到 server.log 。
- 后台启动 Vue 前端服务，日志输出到 web.log 。

停止，根据端口号查找并终止前后端进程

``` shell
./run.sh stop
```


重启 
``` shell
./run.sh restart
```