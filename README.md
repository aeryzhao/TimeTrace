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

## Docker Compose 部署

项目已提供容器化部署支持，适合快速在本地或服务器上拉起完整服务。

启动服务：

```shell
docker compose up -d --build
```

停止服务：

```shell
docker compose down
```

访问地址：

- 前端：`http://localhost:5173`
- 后端：通过前端同源代理访问 `/api/v1`

说明：

- `web` 服务会构建前端静态资源，并通过 `nginx` 提供访问。
- `server` 服务运行 Go 后端，并将 SQLite 数据库存放在 Docker volume `timetrace-data` 中。
- 如需清空容器内数据库，可执行 `docker compose down -v` 删除持久化卷。

## 开发说明

- 本地 `web` 开发环境已通过 Vite 代理 `/api` 到 `http://localhost:8080`。
- 前端默认请求同源 `/api/v1`，也可通过 `VITE_API_BASE_URL` 环境变量覆盖。
