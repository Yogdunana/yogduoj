# YogduOJ - 在线评测系统

<p align="center">
  <strong>YogduOJ</strong> - 自主研发的在线评测系统，适用于程序设计竞赛、算法竞赛及 CTF 网络安全训练。
</p>

<p align="center">
  简体中文 | <a href="./README.md">English</a>
</p>

---

## 项目简介

**YogduOJ** 是一个功能完善的自主研发在线评测系统（Online Judge），专为程序设计竞赛、算法竞赛和 CTF 网络安全训练而设计。采用现代化技术栈构建，支持多种竞赛模式，为个人和团队竞赛提供全方位的平台支持。

## 功能特性

### 竞赛模式

YogduOJ 内置支持 **8 种竞赛模式**：

| 模式 | 说明 |
|------|------|
| **ACM** | 经典 ACM-ICPC 赛制，实时排名 |
| **OI** | 信息学奥林匹克，支持部分分 |
| **IOI** | 国际信息学奥林匹克赛制 |
| **CF** | Codeforces 风格的积分赛制 |
| **CTF** | 夺旗赛网络安全挑战 |
| **AWD** | 攻防对抗 - 实时安全竞赛 |
| **ISW** | 信息安全攻防赛制 |
| **DIY** | 完全自定义竞赛模式 |

### 核心功能

- **个人赛与团队赛** - 同时支持个人参赛和团队协作竞赛
- **多语言支持** - C、C++、Java、Python、Go、Rust 等主流编程语言
- **AI 辅助出题** - 预留 AI 辅助题目生成接口
- **跨平台题目导入** - 支持从其他 OJ 平台导入题目
- **防作弊检测** - 内置代码相似度检测与抄袭预防机制
- **实时排名** - 实时排行榜，支持多种展示方式
- **沙箱评测** - 自研沙箱评测服务，确保代码安全执行
- **Docker 部署** - 基于 Docker 和 Docker Compose 的一键部署

### 系统功能

- 用户注册、登录及基于角色的权限控制（管理员 / 参赛者 / 评委）
- 题目管理，支持富文本编辑器、测试用例上传、样例输入输出
- 竞赛管理，支持灵活的赛程安排和权限控制
- 提交记录、代码查看与统计分析
- 公告与讨论系统
- 文件与附件管理

## 技术栈

| 层级 | 技术 |
|------|------|
| **前端** | Vue 3 + TypeScript + Vite + Pinia + Element Plus |
| **后端** | Go + Gin + GORM |
| **数据库** | MySQL 8.0 |
| **评测** | 自研沙箱（Go） |
| **缓存** | Redis |
| **部署** | Docker + Docker Compose + Nginx |

## 项目结构

```
yogduoj/
├── backend/          # Go 后端服务（Gin + GORM）
│   ├── cmd/          # 入口文件（server、migrate 等）
│   ├── internal/     # 内部包（handlers、models、services）
│   ├── config/       # 配置文件
│   └── docs/         # Swagger 文档
├── frontend/         # Vue 3 前端
│   ├── src/          # 源代码（components、views、stores、router）
│   ├── public/       # 静态资源
│   └── vite.config.ts
├── judge/            # 评测服务（Go）
│   ├── cmd/          # 评测入口
│   ├── sandbox/      # 沙箱实现
│   └── runners/      # 语言运行器
├── deploy/           # 部署配置
│   ├── docker-compose.yml
│   ├── docker-compose.dev.yml
│   ├── Dockerfile.backend
│   ├── Dockerfile.frontend
│   └── Dockerfile.judge
├── Makefile          # 构建自动化
├── .gitignore
├── README.md
└── README_CN.md
```

## 快速开始

### 环境要求

- **Go** >= 1.21
- **Node.js** >= 18
- **pnpm** >= 8
- **MySQL** >= 8.0
- **Redis** >= 7.0
- **Docker** 和 **Docker Compose**（容器化部署）

### 开发模式（Docker）

使用 Docker Compose 是最简单的启动方式：

```bash
# 克隆仓库
git clone https://github.com/Yogdunana/YogduOJ.git
cd YogduOJ

# 启动开发环境
make dev
```

### 开发模式（本地）

```bash
# 安装依赖并构建
make build

# 启动后端（终端 1）
make backend

# 启动前端（终端 2）
make frontend

# 启动评测服务（终端 3）
make judge
```

### 构建与部署

```bash
# 构建所有服务
make build

# 运行测试
make test

# 运行代码检查
make lint

# 部署到生产环境
make deploy
```

## Makefile 命令

| 命令 | 说明 |
|------|------|
| `make dev` | 启动开发环境（Docker Compose） |
| `make build` | 构建所有服务 |
| `make backend` | 本地构建并运行后端 |
| `make frontend` | 本地构建并运行前端开发服务器 |
| `make judge` | 本地构建并运行评测服务 |
| `make test` | 运行所有测试 |
| `make lint` | 运行所有代码检查 |
| `make clean` | 清理构建产物 |
| `make deploy` | 部署到生产环境 |
| `make help` | 显示所有可用命令 |

## 配置说明

### 后端配置

编辑 `backend/config/config.yaml` 进行配置：

- 数据库连接（MySQL）
- Redis 连接
- JWT 密钥与过期时间
- 评测服务连接
- 文件上传设置

### 前端配置

在 `frontend/` 目录下创建 `.env.local` 文件：

```env
VITE_API_BASE_URL=http://localhost:8080/api
VITE_WS_URL=ws://localhost:8080/ws
```

## 许可证

本项目由 [Yogdunana-悠渡](https://github.com/Yogdunana) 开发和维护，保留所有权利。

---

<p align="center">
  Powered by YogduOJ, Copyright Yogdunana-悠渡
</p>
