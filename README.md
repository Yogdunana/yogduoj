# YogduOJ - Online Judge System

<p align="center">
  <strong>YogduOJ</strong> - A self-developed Online Judge system for programming contests, algorithm competitions, and CTF cybersecurity practice.
</p>

<p align="center">
  <a href="./README_CN.md">简体中文</a> | English
</p>

---

## Overview

**YogduOJ** is a full-featured, self-developed Online Judge (OJ) system designed for programming contests, algorithm competitions, and CTF cybersecurity training. Built with a modern tech stack, it supports a wide range of contest formats and provides a comprehensive platform for both individual and team-based competitions.

## Features

### Contest Formats

YogduOJ supports **8 contest formats** out of the box:

| Format | Description |
|--------|-------------|
| **ACM** | Classic ACM-ICPC scoring with real-time ranking |
| **OI** | Olympiad in Informatics with partial scoring |
| **IOI** | International Olympiad in Informatics format |
| **CF** | Codeforces-style rating-based contests |
| **CTF** | Capture The Flag cybersecurity challenges |
| **AWD** | Attack With Defense - real-time security competition |
| **ISW** | Information Security Warfare format |
| **DIY** | Fully customizable contest format |

### Core Features

- **Individual & Team Modes** - Support for both solo participants and team competitions
- **Multi-language Support** - C, C++, Java, Python, Go, Rust, and more
- **AI-Assisted Problem Generation** - Interface reserved for AI-powered problem creation
- **Cross-Platform Problem Import** - Import problems from other OJ platforms
- **Anti-Cheat Detection** - Built-in code similarity detection and plagiarism prevention
- **Real-time Ranking** - Live scoreboard with customizable display options
- **Sandbox Judge** - Custom-built sandboxed judge service for secure code execution
- **Docker Deployment** - Easy deployment with Docker and Docker Compose

### System Features

- User registration, login, and role-based access control (Admin / Contestant / Judge)
- Problem management with rich text editors, test case upload, and sample I/O
- Contest management with flexible scheduling and permission controls
- Submission history, code review, and statistics
- Announcement and discussion system
- File and attachment management

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Frontend** | Vue 3 + TypeScript + Vite + Pinia + Element Plus |
| **Backend** | Go + Gin + GORM |
| **Database** | MySQL 8.0 |
| **Judge** | Custom sandbox (Go) |
| **Cache** | Redis |
| **Deployment** | Docker + Docker Compose + Nginx |

## Project Structure

```
yogduoj/
├── backend/          # Go backend service (Gin + GORM)
│   ├── cmd/          # Entry points (server, migrate, etc.)
│   ├── internal/     # Internal packages (handlers, models, services)
│   ├── config/       # Configuration files
│   └── docs/         # Swagger documentation
├── frontend/         # Vue 3 frontend
│   ├── src/          # Source code (components, views, stores, router)
│   ├── public/       # Static assets
│   └── vite.config.ts
├── judge/            # Judge service (Go)
│   ├── cmd/          # Judge entry point
│   ├── sandbox/      # Sandbox implementation
│   └── runners/      # Language runners
├── deploy/           # Deployment configurations
│   ├── docker-compose.yml
│   ├── docker-compose.dev.yml
│   ├── Dockerfile.backend
│   ├── Dockerfile.frontend
│   └── Dockerfile.judge
├── Makefile          # Build automation
├── .gitignore
├── README.md
└── README_CN.md
```

## Quick Start

### Prerequisites

- **Go** >= 1.21
- **Node.js** >= 18
- **pnpm** >= 8
- **MySQL** >= 8.0
- **Redis** >= 7.0
- **Docker** & **Docker Compose** (for containerized deployment)

### Development (Docker)

The easiest way to get started is with Docker Compose:

```bash
# Clone the repository
git clone https://github.com/Yogdunana/YogduOJ.git
cd YogduOJ

# Start development environment
make dev
```

### Development (Local)

```bash
# Install dependencies
make build

# Start backend (terminal 1)
make backend

# Start frontend (terminal 2)
make frontend

# Start judge service (terminal 3)
make judge
```

### Build & Deploy

```bash
# Build all services
make build

# Run tests
make test

# Run linters
make lint

# Deploy to production
make deploy
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make dev` | Start development environment (Docker Compose) |
| `make build` | Build all services |
| `make backend` | Build and run backend locally |
| `make frontend` | Build and run frontend dev server |
| `make judge` | Build and run judge service |
| `make test` | Run all tests |
| `make lint` | Run all linters |
| `make clean` | Clean build artifacts |
| `make deploy` | Deploy to production |
| `make help` | Show all available commands |

## Configuration

### Backend Configuration

Edit `backend/config/config.yaml` to configure:

- Database connection (MySQL)
- Redis connection
- JWT secret and expiration
- Judge service connection
- File upload settings

### Frontend Configuration

Create a `.env.local` file in the `frontend/` directory:

```env
VITE_API_BASE_URL=http://localhost:8080/api
VITE_WS_URL=ws://localhost:8080/ws
```

## License

This project is developed and maintained by [Yogdunana-悠渡](https://github.com/Yogdunana). All rights reserved.

---

<p align="center">
  Powered by YogduOJ, Copyright Yogdunana-悠渡
</p>
