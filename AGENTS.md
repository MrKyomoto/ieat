# 项目协作约定

开始任务前先阅读 `CONTEXT.md`、相关 `docs/adr/` 和 `TODO.md`。业务用语以 `CONTEXT.md` 为准；发现新概念或歧义时，先更新上下文再实现。

## 代码边界

- 后端保持 Go + Chi + pgx + PostgreSQL 模块化单体，不增加微服务、ORM、通用数据库适配层或 JWT。
- 前端保持 React + TypeScript + Vite + Ant Design 单套响应式 Web，不维护独立移动端项目。
- 身份来源和消费流水格式是明确的可变边界；其他位置不要为假设中的未来需求添加接口或工厂。
- 金额使用整数分，时间入库使用 UTC、展示和经营统计使用北京时间。
- 上传文件写入 `UPLOAD_DIR`，数据库不保存图片二进制。
- 日常开发默认使用 `DATA_BACKEND=mock`，不得要求开发者在 Windows 安装 PostgreSQL；需要持久化或数据库集成验证时才显式切换到 `postgres`。
- Mock 和 PostgreSQL 只实现各业务模块所需的小接口，不创建通用数据库抽象层。

## 开发要求

- 新功能和缺陷修复优先按 `.agents/skills/tdd/SKILL.md` 做小步测试驱动开发；测试描述使用 `CONTEXT.md` 中的领域词汇。
- 后端提交前运行 `go test ./...` 和 `go vet ./...`。
- 前端提交前运行 `npm --prefix web run build`。
- 只修改当前任务涉及的模块；跨模块契约通过 HTTP API 或数据库迁移明确表达。
- 代码标识符和提交信息使用英文，面向用户的文案及项目文档使用中文。
- `TODO.md` 由项目负责人分配；完成任务时勾选对应项，不自行重排其他成员任务。
