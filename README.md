# IEat 校园食堂评价

IEat 是校内食堂窗口评价社区，同时提供食堂管理后台。目前默认使用内存 Mock 数据，**手动测试不需要安装 PostgreSQL**。

## 1. 准备环境

请先安装：

- Go 1.26 或更高版本
- Node.js 22.12 或更高版本（自带 npm）
- Git

可用下面的命令确认：

```powershell
go version
node --version
npm --version
```

首次克隆后手动准备项目配置和前端依赖。Windows PowerShell：

```powershell
Copy-Item .env.example .env
npm.cmd --prefix web install
```

Linux：

```bash
cp .env.example .env
npm --prefix web install
```

如果 `.env` 已存在，不要覆盖。启动脚本只启动对应进程，不检查环境、不安装依赖。

## 2. 分别启动前后端

打开两个终端。在 Windows PowerShell 中分别运行：

```powershell
# 终端一：后端
powershell -ExecutionPolicy Bypass -File .\scripts\backend.ps1

# 终端二：前端
powershell -ExecutionPolicy Bypass -File .\scripts\frontend.ps1
```

Linux：

```bash
# 终端一：后端
bash ./scripts/backend.sh

# 终端二：前端
bash ./scripts/frontend.sh
```

启动后访问：

- 前端：http://localhost:5173
- 后端健康检查：http://localhost:8080/healthz

修改前端代码后，Vite 会自动刷新页面。修改后端代码后，只需在后端终端按 `Ctrl+C`，再重新运行后端脚本。两个进程分别使用 `Ctrl+C` 停止。

## 3. 登录账号

三个账号的密码都是 `ieat-dev-only`。

| 角色 | 邮箱 |
| --- | --- |
| 普通用户 | `student@mail.ustc.edu.cn` |
| 管理部门人员 | `manager@ustc.edu.cn` |
| 平台管理员 | `admin@ustc.edu.cn` |

## 4. 手动测试清单

- 使用普通用户登录，确认能看到“示例食堂”、两层楼和两个示例窗口。
- 退出后使用管理部门人员登录，确认页面显示对应角色和经营看板入口。
- 使用平台管理员登录，确认页面显示平台管理入口。
- 在浏览器开发者工具中切换到手机尺寸（例如宽度 390px），检查页面和导航是否正常适配。
- 访问健康检查地址，确认返回 `{\"status\":\"ok\"}`。

Mock 数据只保存在当前后端进程中，重启后会恢复初始状态。

## 常见问题

- 提示缺少 `go`、`node` 或 `npm`：先安装对应开发环境，并重新打开终端。
- 前端提示缺少依赖：先运行 `npm.cmd --prefix web install`；Linux 使用 `npm --prefix web install`。
- 端口占用：关闭占用 `5173` 或 `8080` 端口的程序后重试。
- 登录返回 `403`：确认使用的是 `http://localhost:5173`，并检查 `.env` 中 `WEB_ORIGIN=http://localhost:5173`。
- 意外连接数据库：检查 `.env` 中是否为 `DATA_BACKEND=mock`。

## 开发检查

```bash
go test ./...
go vet ./...
npm --prefix web run build
```

以后需要测试 PostgreSQL 时，再将 `.env` 中的 `DATA_BACKEND` 改为 `postgres` 并配置 `DATABASE_URL`；当前阶段无需处理。项目规则和已确认设计见 [CONTEXT.md](CONTEXT.md)、[架构决策](docs/adr/) 和 [TODO.md](TODO.md)。
