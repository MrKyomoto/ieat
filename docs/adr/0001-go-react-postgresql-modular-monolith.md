# 采用 Go、React 和 PostgreSQL 构建模块化单体

项目采用单套响应式 React + TypeScript + Vite + Ant Design 前端，以及 Go + Chi + pgx + PostgreSQL 后端，并以模块化单体形式部署。该组合满足学生社区、窗口级授权、流水导入和统计需求，同时保持四人团队及 Agent 辅助开发的反馈链路简短；暂不采用 Rust、微服务或通用数据库适配层，只为尚未确定的身份来源和可能变化的流水格式保留明确边界。
