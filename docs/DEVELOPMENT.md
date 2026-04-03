# YogduOJ Development Guide

This guide covers how to set up, develop, and contribute to YogduOJ locally.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Quick Start with Docker Compose](#quick-start-with-docker-compose)
3. [Local Development Setup](#local-development-setup)
   - [Backend](#backend-setup)
   - [Frontend](#frontend-setup)
   - [Judge Service](#judge-service-setup)
4. [Project Structure](#project-structure)
5. [Code Style Guidelines](#code-style-guidelines)
6. [Testing Commands](#testing-commands)
7. [i18n: Adding New Translations](#i18n-adding-new-translations)
8. [Adding New API Endpoints](#adding-new-api-endpoints)
9. [Adding New Frontend Pages](#adding-new-frontend-pages)

---

## Prerequisites

Install the following tools before starting development:

| Tool | Version | Purpose |
|------|---------|---------|
| [Go](https://go.dev/dl/) | 1.22+ | Backend language |
| [Node.js](https://nodejs.org/) | 20+ (LTS) | Frontend build toolchain |
| [npm](https://www.npmjs.com/) | 10+ | Package manager (or pnpm 8+) |
| [Docker](https://docs.docker.com/get-docker/) | 24+ | Container runtime |
| [Docker Compose](https://docs.docker.com/compose/install/) | v2+ | Multi-container orchestration |
| [MySQL](https://dev.mysql.com/downloads/mysql/) | 8.0+ | Database (or use Docker) |
| [Git](https://git-scm.com/) | 2.40+ | Version control |
| [Make](https://www.gnu.org/software/make/) | any | Build automation (optional but recommended) |

### Optional but Recommended

- [golangci-lint](https://golangci-lint.run/usage/install/) -- Go linting aggregator
- [VS Code](https://code.visualstudio.com/) with the following extensions:
  - `golang.Go` -- Go language support
  - `vue.volar` -- Vue 3 + TypeScript support
  - `dbaeumer.vscode-eslint` -- ESLint integration

---

## Quick Start with Docker Compose

The fastest way to get the full stack running is via Docker Compose. This starts MySQL, the backend API, the frontend (served by Nginx), and the judge service.

```bash
# Clone the repository
git clone https://github.com/Yogdunana/yogduoj.git
cd yogduoj

# Start the development environment
make dev

# View logs from all services
make dev-logs

# Stop the environment
make dev-down
```

Or equivalently, using Docker Compose directly:

```bash
# Development environment
docker-compose -f deploy/docker-compose.dev.yml up --build -d

# Production environment
docker-compose -f deploy/docker-compose.yml up -d

# Stop
docker-compose -f deploy/docker-compose.dev.yml down
```

After starting, the services are available at:

| Service | URL |
|---------|-----|
| Frontend (Nginx) | http://localhost |
| Backend API | http://localhost:8080 |
| API Health Check | http://localhost:8080/api/v1/health |
| MySQL | localhost:3306 |

---

## Local Development Setup

For active development, it is recommended to run the backend and frontend individually (not in Docker) so you get hot-reload and easier debugging. Only MySQL needs to run in Docker.

### Backend Setup

1. **Start MySQL** (if not already running):

```bash
docker run -d \
  --name yogduoj-mysql \
  -e MYSQL_ROOT_PASSWORD=yogduoj_dev \
  -e MYSQL_DATABASE=yogduoj \
  -p 3306:3306 \
  mysql:8.0 \
  --character-set-server=utf8mb4 \
  --collation-server=utf8mb4_unicode_ci
```

2. **Configure the backend**:

Edit `backend/config.yaml` to match your local MySQL credentials:

```yaml
server:
  port: 8080
  mode: debug          # debug | release | test

database:
  host: localhost
  port: 3306
  user: root
  password: yogduoj_dev
  dbname: yogduoj
  charset: utf8mb4
  parse_time: true
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

jwt:
  secret: "yogduoj-default-secret-change-in-production"
  access_ttl: 2h
  refresh_ttl: 168h

rate_limit:
  requests_per_minute: 60
  burst: 100

log:
  level: debug
  file_path: logs/app.log
  max_size: 100
  max_backups: 10
  max_age: 30
  compress: true
```

Alternatively, override config values with environment variables using the `YOGDUOJ_` prefix:

```bash
export YOGDUOJ_DATABASE_PASSWORD=yogduoj_dev
export YOGDUOJ_SERVER_MODE=debug
```

3. **Install dependencies and run**:

```bash
cd backend
go mod download
go run ./cmd/server
```

The backend will auto-migrate all database tables on startup. You should see logs like:

```
INFO  Database connected successfully
INFO  Database migrations completed
INFO  Server starting  {"address": ":8080"}
```

4. **Verify the backend**:

```bash
curl http://localhost:8080/api/v1/health
# Expected: {"service":"yogduoj-backend","status":"ok"}
```

### Frontend Setup

1. **Install dependencies**:

```bash
cd frontend
npm install
# or: pnpm install
```

2. **Start the development server**:

```bash
npm run dev
# or: pnpm run dev
```

The dev server starts at `http://localhost:5173` by default with hot-module replacement (HMR).

3. **Build for production**:

```bash
npm run build
# or: pnpm run build
```

Output goes to `frontend/dist/`.

4. **Preview the production build**:

```bash
npm run preview
```

### Judge Service Setup

The judge service is a separate Go binary that communicates with the backend via gRPC.

```bash
cd judge
go mod download
go run ./cmd/judge
```

> **Note:** The judge service requires Linux namespaces, cgroups, and seccomp for sandboxing. It only runs on Linux hosts. On macOS/Windows, run it inside a Docker container.

---

## Project Structure

```
yogduoj/
├── backend/                    # Go backend service
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # Application entry point
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go       # Configuration loading (Viper)
│   │   ├── handler/            # HTTP handlers (controllers)
│   │   │   ├── admin_handler.go
│   │   │   ├── announcement_handler.go
│   │   │   ├── auth_handler.go
│   │   │   ├── contest_handler.go
│   │   │   ├── problem_handler.go
│   │   │   ├── submission_handler.go
│   │   │   ├── team_handler.go
│   │   │   └── user_handler.go
│   │   ├── middleware/          # Gin middleware
│   │   │   ├── auth.go         # JWT authentication
│   │   │   ├── cors.go         # CORS configuration
│   │   │   ├── i18n.go         # Internationalization
│   │   │   ├── logger.go       # Request logging (Zap)
│   │   │   └── ratelimit.go    # Rate limiting
│   │   ├── migration/
│   │   │   └── migrate.go      # GORM auto-migration
│   │   ├── model/              # GORM models (22 models)
│   │   │   ├── user.go
│   │   │   ├── team.go
│   │   │   ├── problem.go
│   │   │   ├── submission.go
│   │   │   ├── contest.go
│   │   │   ├── announcement.go
│   │   │   ├── cheat_record.go
│   │   │   ├── ai_record.go
│   │   │   ├── import_record.go
│   │   │   ├── ctf_resource.go
│   │   │   └── system_config.go
│   │   ├── repository/         # Data access layer
│   │   │   ├── user_repo.go
│   │   │   ├── team_repo.go
│   │   │   ├── problem_repo.go
│   │   │   ├── submission_repo.go
│   │   │   ├── contest_repo.go
│   │   │   └── announcement_repo.go
│   │   ├── router/
│   │   │   └── router.go       # Route registration
│   │   ├── service/            # Business logic layer
│   │   │   ├── auth_service.go
│   │   │   ├── user_service.go
│   │   │   ├── team_service.go
│   │   │   ├── problem_service.go
│   │   │   ├── submission_service.go
│   │   │   ├── contest_service.go
│   │   │   ├── announcement_service.go
│   │   │   ├── judge_service.go
│   │   │   ├── anti_cheat_service.go
│   │   │   ├── ai_service.go
│   │   │   ├── import_service.go
│   │   │   └── system_service.go
│   │   └── pkg/
│   │       └── jwt/            # JWT token management
│   ├── locales/                # Backend i18n files
│   │   ├── en.yaml             # English translations
│   │   └── zh.yaml             # Chinese translations
│   ├── config.yaml             # Default configuration
│   ├── Dockerfile              # Multi-stage Docker build
│   ├── go.mod
│   └── go.sum
│
├── frontend/                   # Vue 3 frontend
│   ├── src/
│   │   ├── api/                # API client modules
│   │   │   ├── index.ts        # Axios instance + interceptors
│   │   │   ├── auth.ts
│   │   │   ├── user.ts
│   │   │   ├── problem.ts
│   │   │   ├── submission.ts
│   │   │   ├── contest.ts
│   │   │   ├── team.ts
│   │   │   ├── announcement.ts
│   │   │   └── admin.ts
│   │   ├── assets/
│   │   │   └── styles/
│   │   │       ├── global.scss
│   │   │       └── variables.scss
│   │   ├── components/
│   │   │   ├── common/         # Shared components
│   │   │   │   ├── ConfirmDialog.vue
│   │   │   │   ├── LoadingSpinner.vue
│   │   │   │   └── Pagination.vue
│   │   │   └── layout/         # Layout components
│   │   │       ├── AdminLayout.vue
│   │   │       ├── AppFooter.vue
│   │   │       ├── AppHeader.vue
│   │   │       ├── AppLayout.vue
│   │   │       └── AppSidebar.vue
│   │   ├── i18n/               # Frontend i18n
│   │   │   ├── index.ts
│   │   │   ├── en.json
│   │   │   └── zh.json
│   │   ├── router/
│   │   │   └── index.ts        # Vue Router + navigation guards
│   │   ├── stores/
│   │   │   ├── auth.ts         # Authentication state
│   │   │   └── theme.ts        # Theme (dark/light) state
│   │   ├── types/
│   │   │   └── index.ts        # TypeScript type definitions
│   │   ├── views/              # Page components
│   │   │   ├── Home.vue
│   │   │   ├── Login.vue
│   │   │   ├── Register.vue
│   │   │   ├── ProblemList.vue
│   │   │   ├── ProblemDetail.vue
│   │   │   ├── SubmissionList.vue
│   │   │   ├── SubmissionDetail.vue
│   │   │   ├── ContestList.vue
│   │   │   ├── ContestDetail.vue
│   │   │   ├── ContestRanking.vue
│   │   │   ├── TeamList.vue
│   │   │   ├── TeamCreate.vue
│   │   │   ├── TeamDetail.vue
│   │   │   ├── UserProfile.vue
│   │   │   ├── UserPublicProfile.vue
│   │   │   ├── CTFPractice.vue
│   │   │   ├── CTFCategoryView.vue
│   │   │   ├── AnnouncementList.vue
│   │   │   ├── AnnouncementDetail.vue
│   │   │   ├── HelpCenter.vue
│   │   │   └── admin/          # Admin panel views
│   │   │       ├── Dashboard.vue
│   │   │       ├── UserManage.vue
│   │   │       ├── ProblemManage.vue
│   │   │       ├── ProblemCreate.vue
│   │   │       ├── ProblemEdit.vue
│   │   │       ├── ContestManage.vue
│   │   │       ├── ContestCreate.vue
│   │   │       ├── ContestEdit.vue
│   │   │       ├── SubmissionManage.vue
│   │   │       ├── JudgeMonitor.vue
│   │   │       ├── AnnouncementManage.vue
│   │   │       ├── AIProblemManage.vue
│   │   │       ├── ImportManage.vue
│   │   │       ├── CheatManage.vue
│   │   │       ├── SystemConfig.vue
│   │   │       └── DIYTemplates.vue
│   │   ├── App.vue
│   │   └── main.ts
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── package.json
│   └── index.html
│
├── judge/                      # Judge service (gRPC)
│   ├── cmd/
│   │   └── judge/
│   ├── proto/                  # gRPC protobuf definitions
│   └── sandbox/                # Sandbox implementation
│
├── deploy/                     # Deployment configurations
│   ├── docker-compose.yml          # Production compose
│   └── docker-compose.dev.yml      # Development compose
│
├── docs/                       # Documentation
│   ├── DEV_LOG.md
│   ├── DEVELOPMENT.md
│   └── DEPLOYMENT.md
│
├── .gitignore
├── Makefile                    # Build automation
├── README.md                   # English README
└── README_CN.md                # Chinese README
```

### Architecture Overview

YogduOJ follows a **three-layer architecture** on the backend:

```
HTTP Request
    |
    v
[Router] --> [Middleware Chain: CORS -> Logger -> Rate Limit -> Auth]
    |
    v
[Handler]  -- validates request, calls service, returns response
    |
    v
[Service]  -- business logic, orchestrates repositories
    |
    v
[Repository] -- database CRUD via GORM
    |
    v
[MySQL]
```

The frontend uses a **component-based architecture** with:

```
[Vue Router] --> [Navigation Guards (auth check)]
    |
    v
[Layout Component] (AppLayout / AdminLayout)
    |
    v
[View Component] --> [Pinia Store] --> [API Module] --> [Backend]
```

---

## Code Style Guidelines

### Go Backend

- Follow [Effective Go](https://go.dev/doc/effective_go) and the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
- Use `gofmt` for formatting (run `go fmt ./...`).
- Use `golangci-lint` for linting (`golangci-lint run ./...`).
- Package naming: lowercase, single word, no underscores.
- Exported functions must have doc comments.
- Error handling: always check errors; use `fmt.Errorf("context: %w", err)` for wrapping.
- Use structured logging via `zap` (never use `fmt.Println` for logging).

**Example handler:**

```go
// GetUser retrieves a user by ID.
// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    user, err := h.userService.GetUserByID(uint(id))
    if err != nil {
        zap.L().Error("failed to get user", zap.Error(err))
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": user})
}
```

### TypeScript / Vue Frontend

- Use TypeScript strict mode; avoid `any` when possible.
- Use `<script setup lang="ts">` for Vue Single File Components.
- Follow Vue 3 Composition API style (no Options API).
- Use Naive UI components; do not import Element Plus.
- Use `vue-i18n` for all user-facing strings (no hardcoded text).
- Use Pinia stores for shared state; avoid prop drilling for deep state.
- Use SCSS for component styles; leverage CSS variables from `variables.scss`.
- File naming: PascalCase for components (`ProblemList.vue`), camelCase for utilities.

**Example component:**

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NCard, NButton } from 'naive-ui'

const { t } = useI18n()
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  // fetch data...
  loading.value = false
})
</script>

<template>
  <NCard :title="t('problems.title')">
    <NButton :loading="loading">
      {{ t('common.refresh') }}
    </NButton>
  </NCard>
</template>
```

---

## Testing Commands

### Run All Tests

```bash
make test
```

### Backend Tests

```bash
# Run all backend tests with race detection and coverage
cd backend
go test ./... -race -coverprofile=coverage.out

# Run a specific package's tests
go test ./internal/service -v

# View coverage report
go tool cover -func=coverage.out
```

### Frontend Tests

```bash
cd frontend
npm run test:unit
# or: pnpm run test:unit
```

### Linting

```bash
# Lint everything
make lint

# Backend only
cd backend && golangci-lint run ./...

# Frontend only
cd frontend && npm run lint
# or: pnpm run lint
```

### Building

```bash
# Build all
make build

# Backend binary
make build-backend    # outputs to bin/yogduoj-server

# Frontend assets
make build-frontend   # outputs to frontend/dist/

# Judge binary
make build-judge      # outputs to bin/yogduoj-judge
```

### Cleaning

```bash
make clean
```

---

## i18n: Adding New Translations

YogduOJ supports English (EN) and Chinese (ZH) on both the backend and frontend.

### Backend i18n

Backend translations are stored in `backend/locales/` as YAML files.

1. Open `backend/locales/en.yaml` and add your new message key:

```yaml
# existing keys...
newMessageKey: "This is a new message"
```

2. Open `backend/locales/zh.yaml` and add the Chinese translation:

```yaml
newMessageKey: "这是一条新消息"
```

3. Use the key in your handler or service via the `go-i18n` localizer:

```go
msg := localizer.MustLocalize(&i18n.LocalizeConfig{
    MessageID: "newMessageKey",
})
```

### Frontend i18n

Frontend translations are stored in `frontend/src/i18n/` as JSON files.

1. Open `frontend/src/i18n/en.json` and add your key under the appropriate section:

```json
{
  "problems": {
    "newField": "New Field Label"
  }
}
```

2. Open `frontend/src/i18n/zh.json` and add the Chinese translation:

```json
{
  "problems": {
    "newField": "新字段标签"
  }
}
```

3. Use the key in your Vue component:

```vue
<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<template>
  <span>{{ t('problems.newField') }}</span>
</template>
```

### Adding a New Language

To add a new language (e.g., Japanese):

1. **Backend**: Create `backend/locales/ja.yaml` with all keys translated.
2. **Frontend**: Create `frontend/src/i18n/ja.json` with all keys translated.
3. **Frontend registration**: Update `frontend/src/i18n/index.ts` to import and register the new locale.
4. **Backend registration**: Update the i18n middleware initialization to support the new locale code.

---

## Adding New API Endpoints

Follow these steps to add a new API endpoint to the backend.

### Step 1: Define the Model (if needed)

If the endpoint involves a new database entity, create a model file in `backend/internal/model/`:

```go
// backend/internal/model/my_entity.go
package model

import "gorm.io/gorm"

type MyEntity struct {
    gorm.Model
    Name   string `json:"name" gorm:"size:255;not null"`
    UserID uint   `json:"user_id" gorm:"index"`
}
```

Register it in `backend/internal/migration/migrate.go`:

```go
db.AutoMigrate(&model.MyEntity{})
```

### Step 2: Create the Repository

```go
// backend/internal/repository/my_entity_repo.go
package repository

import (
    "github.com/Yogdunana/yogduoj/backend/internal/model"
    "gorm.io/gorm"
)

type MyEntityRepository struct {
    db *gorm.DB
}

func NewMyEntityRepository(db *gorm.DB) *MyEntityRepository {
    return &MyEntityRepository{db: db}
}

func (r *MyEntityRepository) Create(entity *model.MyEntity) error {
    return r.db.Create(entity).Error
}

func (r *MyEntityRepository) GetByID(id uint) (*model.MyEntity, error) {
    var entity model.MyEntity
    err := r.db.First(&entity, id).Error
    return &entity, err
}

func (r *MyEntityRepository) List(page, pageSize int) ([]model.MyEntity, int64, error) {
    var entities []model.MyEntity
    var total int64
    r.db.Model(&model.MyEntity{}).Count(&total)
    err := r.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities).Error
    return entities, total, err
}
```

### Step 3: Create the Service

```go
// backend/internal/service/my_entity_service.go
package service

import (
    "github.com/Yogdunana/yogduoj/backend/internal/model"
    "github.com/Yogdunana/yogduoj/backend/internal/repository"
)

type MyEntityService struct {
    repo *repository.MyEntityRepository
}

func NewMyEntityService(repo *repository.MyEntityRepository) *MyEntityService {
    return &MyEntityService{repo: repo}
}

func (s *MyEntityService) GetByID(id uint) (*model.MyEntity, error) {
    return s.repo.GetByID(id)
}
```

### Step 4: Create the Handler

```go
// backend/internal/handler/my_entity_handler.go
package handler

import (
    "net/http"
    "strconv"

    "github.com/Yogdunana/yogduoj/backend/internal/service"
    "github.com/gin-gonic/gin"
)

type MyEntityHandler struct {
    service *service.MyEntityService
}

func NewMyEntityHandler(service *service.MyEntityService) *MyEntityHandler {
    return &MyEntityHandler{service: service}
}

func (h *MyEntityHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    entity, err := h.service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": entity})
}
```

### Step 5: Register the Route

Update `backend/internal/router/router.go`:

1. Add the handler to the `Router` struct.
2. Accept it in `NewRouter()`.
3. Register the route in `Setup()`.

```go
// In the Router struct
myEntityHandler *handler.MyEntityHandler

// In Setup()
myEntities := api.Group("/my-entities")
{
    myEntities.GET("/:id", r.myEntityHandler.GetByID)
}
```

### Step 6: Wire Up Dependencies

In `backend/cmd/server/main.go`, initialize the repository, service, and handler, then pass them to the router.

### Step 7: Add i18n Keys

Add error/success message keys to both `backend/locales/en.yaml` and `backend/locales/zh.yaml`.

---

## Adding New Frontend Pages

### Step 1: Create the View Component

Create a new `.vue` file in `frontend/src/views/`:

```vue
<!-- frontend/src/views/MyNewPage.vue -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NCard, NDataTable } from 'naive-ui'
import { myApiFunction } from '@/api/myModule'

const { t } = useI18n()
const loading = ref(false)
const data = ref([])

onMounted(async () => {
  loading.value = true
  try {
    const res = await myApiFunction()
    data.value = res.data.data
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="my-new-page">
    <NCard :title="t('myPage.title')">
      <NDataTable :data="data" :loading="loading" />
    </NCard>
  </div>
</template>

<style scoped lang="scss">
.my-new-page {
  padding: 24px;
}
</style>
```

### Step 2: Create the API Module

Create a new API file in `frontend/src/api/`:

```typescript
// frontend/src/api/myModule.ts
import api from './index'
import type { ApiResponse } from '@/types'

export function myApiFunction() {
  return api.get<ApiResponse<MyDataType[]>>('/my-endpoint')
}
```

### Step 3: Register the Route

Add the route to `frontend/src/router/index.ts`:

```typescript
// Under the AppLayout children array:
{
  path: 'my-page',
  name: 'MyNewPage',
  component: () => import('@/views/MyNewPage.vue'),
},
```

If the page requires authentication, add the meta field:

```typescript
{
  path: 'my-page',
  name: 'MyNewPage',
  component: () => import('@/views/MyNewPage.vue'),
  meta: { requiresAuth: true },
},
```

For admin-only pages, add it under the `/admin` route group:

```typescript
// Under the AdminLayout children array:
{
  path: 'my-admin-page',
  name: 'AdminMyPage',
  component: () => import('@/views/admin/MyAdminPage.vue'),
},
```

### Step 4: Add i18n Translations

Add translation keys to `frontend/src/i18n/en.json` and `frontend/src/i18n/zh.json`:

```json
// en.json
{
  "myPage": {
    "title": "My New Page"
  }
}

// zh.json
{
  "myPage": {
    "title": "我的新页面"
  }
}
```

### Step 5: Add Navigation Entry (optional)

If the page should appear in the sidebar or header navigation, update the corresponding layout component (`AppSidebar.vue` or `AppHeader.vue`) with a new router-link entry.

---

## Useful Make Targets

Run `make help` to see all available targets. The most commonly used ones:

| Command | Description |
|---------|-------------|
| `make dev` | Start development environment with Docker Compose |
| `make dev-down` | Stop development environment |
| `make dev-logs` | Tail development logs |
| `make backend` | Run backend locally (with hot-reload via `go run`) |
| `make frontend` | Run frontend dev server locally |
| `make build` | Build all services |
| `make test` | Run all tests |
| `make lint` | Run all linters |
| `make clean` | Clean all build artifacts |
| `make deploy` | Deploy to production |
| `make migrate` | Run database migrations |
