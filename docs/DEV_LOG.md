# YogduOJ Development Log

This file tracks development progress across sessions to prevent context loss.

---

## [2026-04-03] Session 2: Full Feature Implementation (Phases 2-7)

### Completed

#### Phase 2: User System
- User profile page with avatar upload, bio editing, and statistics display
- Team management: create/join/leave teams, team member roles (leader/member)
- Team profile page with member list and team submission statistics
- User settings page (language preference, theme toggle, password change)

#### Phase 3: Problem System
- Problem list page with search, filter by difficulty/category/tags, pagination
- Problem detail page with description, constraints, examples, and submission history
- Code editor integration with syntax highlighting and language selection
- Submission page with real-time status updates (pending/AC/WA/TLE/MLE/RE/CE)
- Problem tag management and category system

#### Phase 4: Judge System
- Sandbox implementation using Linux Namespace + Cgroup + Seccomp
- Container pool: pre-created isolated containers for judging
- gRPC communication between backend and judge service
- Language support: C/C++ (GCC), Java (OpenJDK), Python 3
- Standard problem checker and CTF flag checker
- Anti-cheat service: code similarity detection, submission pattern analysis
- Resource limits enforcement (time, memory, output size)

#### Phase 5: Contest System
- ICPC/IOI contest format support with configurable scoring
- Contest list, detail, and registration pages
- Real-time leaderboard with freeze board support
- Contest announcement system
- Contest problem management (bind existing problems)
- Contest participation tracking and submission filtering

#### Phase 6: Admin Panel
- Admin dashboard with system overview statistics
- User management: list, search, ban/unban, role assignment
- Problem management: create, edit, delete, import/export
- Contest management: create, configure, publish, archive
- Announcement management: create, edit, pin, schedule
- System configuration panel
- Import records and cross-platform import (JSON format)
- Cheat record review and management

#### Phase 7: CTF Module & Polish
- CTF resource model and flag-based checker
- CTF problem submission flow
- Homepage optimization: contest highlights, recent announcements, problem recommendations
- Help center page with FAQ and usage guides
- UI polish: responsive design improvements, loading states, error handling
- i18n completeness for EN/ZH across all pages

#### CI/CD & Deployment Configuration
- GitHub Actions CI pipeline (`.github/workflows/ci.yml`):
  - Backend: go vet, go test with race detection and coverage, go build
  - Frontend: npm ci, npm run build
  - Judge: go vet, go build
  - Triggers on push to main/dev and PRs to main
- GitHub Actions CD pipeline (`.github/workflows/deploy.yml`):
  - Docker Buildx setup and GHCR login
  - Build and push backend, frontend, judge images to ghcr.io/yogdunana/
  - SSH-based deployment to server (101.237.129.33:23196)
  - Post-deploy health check via curl
  - Secrets reminder (SERVER_PASSWORD required)
- Updated `deploy/docker-compose.yml`:
  - Changed build directives to image directives using ghcr.io/yogdunana/yogduoj-*:latest
  - Original build directives kept as comments for local development reference
- Created `deploy/docker-compose.prod.yml`:
  - Production override with resource limits (CPU/memory) for all services
  - JSON file logging with rotation (max-size, max-file) and service tags
  - Restart policies enforced
  - Usage: `docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d`

### In Progress
- Production deployment to server (101.237.129.33)
- GitHub repository creation and initial push
- Docker image build testing

### Next Steps
- Phase 8: Production deployment
  - Push code to GitHub repository
  - Configure GitHub Secrets (SERVER_PASSWORD)
  - Set up server environment (/opt/yogduoj)
  - Deploy via GitHub Actions or manual docker compose
  - DNS and SSL certificate configuration for oj.yogdunana.com
  - End-to-end testing on production

### Tech Stack
- Frontend: Vue3 + TypeScript + Vite + Pinia + Naive UI + vue-i18n + Axios
- Backend: Go 1.22 + Gin + GORM + Viper + Zap + go-i18n + JWT
- Database: MySQL 8.0 (utf8mb4)
- Judge: Custom Go service with gRPC, Linux Namespace + Cgroup + Seccomp sandbox
- Deploy: Docker + docker-compose + Nginx + GitHub Actions CI/CD
- Container Registry: GitHub Container Registry (ghcr.io/yogdunana/)
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
- Production images hosted on GHCR (ghcr.io/yogdunana/yogduoj-*)
- CD uses SSH action as jump host for deployment to server

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
