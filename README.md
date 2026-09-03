# IEat 校园食堂评价

IEat 是面向校内成员的食堂窗口评价社区和管理平台。首版采用 Go + Chi + pgx + PostgreSQL 后端，以及 React + TypeScript + Vite + Ant Design 响应式前端，整体以模块化单体开发。

## 开发环境

请先安装：

- Git
- Go 1.26 或更高版本
- Node.js 22.12 或更高版本；推荐 Node.js 24 LTS
- npm（随 Node.js 安装）
- PostgreSQL 16～18

Go 的当前稳定版本和支持策略见 [Go Release History](https://go.dev/doc/devel/release)，Node.js 应选择仍在维护的 LTS 版本，见 [Node.js Releases](https://nodejs.org/en/about/previous-releases)。本项目使用的 Vite 要求 Node.js 20.19+ 或 22.12+，见 [Vite 8](https://vite.dev/blog/announcing-vite8)。PostgreSQL 各主版本支持周期见 [PostgreSQL Versioning Policy](https://www.postgresql.org/support/versioning/)。

Linux 还需要 Bash；Windows 启动脚本使用系统自带的 PowerShell。开发脚本只管理项目依赖和进程，不安装上述系统软件。

## 准备 PostgreSQL

使用 PostgreSQL 管理员账号执行：

```sql
CREATE ROLE ieat WITH LOGIN PASSWORD 'ieat';
CREATE DATABASE ieat_dev OWNER ieat;
CREATE DATABASE ieat_test OWNER ieat;
```

Arch Linux/Omarchy 可以通过下面的命令进入 PostgreSQL：

```bash
sudo -iu postgres psql
```

Windows 通常使用安装 PostgreSQL 时设置的管理员密码：

```powershell
psql -U postgres
```

这些账号和密码只用于个人开发环境，请勿用于试运行或生产部署。

## 首次启动

复制环境配置：

```bash
cp .env.example .env
```

Windows PowerShell：

```powershell
Copy-Item .env.example .env
```

按需修改 `.env` 中的 `DATABASE_URL`。随后运行对应脚本：

```bash
./scripts/dev.sh
```

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\dev.ps1
```

脚本会下载项目依赖、执行数据库迁移、初始化开发数据，并同时启动：

- 前端：http://localhost:5173
- 后端：http://localhost:8080
- 健康检查：http://localhost:8080/healthz

按 `Ctrl+C` 可以停止前后端。

## 开发账号

开发环境关闭注册和邮件发送。运行启动脚本后可使用以下账号登录，密码默认为 `.env` 中的 `DEV_SEED_PASSWORD`：

| 角色 | 邮箱 |
| --- | --- |
| 普通用户 | `student@mail.ustc.edu.cn` |
| 管理部门人员 | `manager@ustc.edu.cn` |
| 平台管理员 | `admin@ustc.edu.cn` |

开发数据初始化只允许在 `APP_ENV=development` 时执行。

## 常用命令

```bash
go run ./cmd/migrate
go run ./cmd/seed
go run ./cmd/api
go test ./...
go vet ./...
npm --prefix web run dev
npm --prefix web run build
```

数据库集成测试应使用独立的 `ieat_test` 数据库，不要清空 `ieat_dev`。

## 目录

```text
cmd/                    Go 可执行程序入口
internal/auth/          登录、会话与访问身份
internal/catalog/       食堂、楼层和窗口目录
internal/database/      PostgreSQL 连接与迁移
internal/devseed/       仅开发环境可用的固定数据
internal/httpapi/       HTTP 路由和通用中间件
web/                    React 响应式前端
scripts/                Linux 与 Windows 开发脚本
docs/adr/               已确认的架构决定
CONTEXT.md              领域词汇和业务规则
TODO.md                 模块任务及负责人
```

开发时先阅读 [CONTEXT.md](CONTEXT.md)、[架构决定](docs/adr/) 和 [TODO.md](TODO.md)。
