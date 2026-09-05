# 本地开发默认使用内存 Mock 数据

Windows 和 Linux 的日常开发默认设置 `DATA_BACKEND=mock`，Go API 在进程内提供固定角色账号、会话和示例食堂，不连接、安装或启动 PostgreSQL；进程重启后这些状态恢复初始值。既有 PostgreSQL 迁移与实现继续保留，只有显式设置 `DATA_BACKEND=postgres` 时才连接数据库。Mock 与 PostgreSQL 分别实现身份和食堂目录模块所需的小接口，不引入覆盖整个系统的通用数据库适配层；后续业务优先完成 Mock 行为，数据库集成集中在独立环境验收。
