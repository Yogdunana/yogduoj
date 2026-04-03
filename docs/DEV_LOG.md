# YogduOJ Development Log

This file tracks development progress across sessions to prevent context loss.

---

## [2026-04-03] Session 1: Project Initialization (Phase 1)

### Completed
- Initialized Monorepo structure with .gitignore, Makefile, README (EN/CN)
- Backend Go project initialized with Gin + GORM + Viper + Zap + go-i18n
- All 22 GORM models defined (users, teams, problems, submissions, contests, announcements, cheat_records, ai_records, import_records, ctf_resources, system_configs, etc.)
- Handler -> Service -> Repository three-layer architecture implemented
- All routes registered: /api/v1/auth, /api/v1/users, /api/v1/teams, /api/v1/problems, /api/v1/submissions, /api/v1/contests, /api/v1/announcements, /api/v1/ctf, /api/v1/admin/*
- JWT auth middleware (access + refresh tokens), CORS, rate limiting, request logging, i18n middleware
- Auth fully implemented: register (bcrypt password hashing, username/email uniqueness), login (credential check, login attempt tracking, JWT generation), logout, refresh
- Frontend Vue3 + TypeScript + Vite project initialized
- Naive UI component library integrated
- Pinia stores (auth, theme), Vue Router with all routes and navigation guards
- vue-i18n with EN/ZH translations
- Layout components: AppHeader, AppSidebar, AppFooter (with "Powered by YogduOJ" copyright), AppLayout, AdminLayout
- Common components: Pagination, LoadingSpinner, ConfirmDialog
- All view pages created (Home with hero section, Login, Register, + placeholder views for all other pages)
- Dark space theme: #1a1a2e background, #00d4ff primary accent
- Docker configuration: docker-compose.yml (production), docker-compose.dev.yml (development)
- Nginx reverse proxy config for oj.yogdunana.com
- Judge service stub with gRPC proto definition
- Backend builds successfully (go build, go vet pass)
- Frontend builds successfully (npm run build, 0 errors)

### In Progress
- Docker deployment testing
- GitHub repository creation and initial push

### Next Steps
- Phase 2: User system (profile page, team management)
- Phase 3: Problem system (problem list, detail, code editor, submission)
- Phase 4: Judge system (sandbox, container pool, anti-cheat)
- Phase 5: Contest system (multi-format, leaderboard, freeze board)
- Phase 6: Admin panel (all management pages)
- Phase 7: CTF module, homepage optimization, help center
- Phase 8: Production deployment, CI/CD, server setup

### Tech Stack
- Frontend: Vue3 + TypeScript + Vite + Pinia + Naive UI + vue-i18n + Axios
- Backend: Go 1.22 + Gin + GORM + Viper + Zap + go-i18n + JWT
- Database: MySQL 8.0 (utf8mb4)
- Judge: Custom Go service with gRPC, Linux Namespace + Cgroup + Seccomp sandbox
- Deploy: Docker + docker-compose + Nginx + GitHub Actions CI/CD
- Domain: oj.yogdunana.com (DNS resolved to 101.237.129.33)
- Server: 101.237.129.33:23196

### Key Decisions
- No Element Plus - using Naive UI for unique UI design
- gRPC for backend-judge communication (not HTTP)
- Code stored on filesystem (Docker Volume), DB only stores paths
- Container pool: pre-create 4 isolated containers
- bcrypt cost=12 for password hashing
- JWT: Access Token 2h + Refresh Token 7d
- Nginx listens on port 23196, forwarded via host reverse proxy to oj.yogdunana.com
- AI features: interface reserved only, no implementation
- Cross-platform import: JSON format first
