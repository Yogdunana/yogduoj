# YogduOJ Deployment Guide

This guide covers deploying YogduOJ to a production server using Docker Compose and Nginx.

---

## Table of Contents

1. [Server Requirements](#server-requirements)
2. [First-Time Deployment](#first-time-deployment)
3. [Environment Variables Reference](#environment-variables-reference)
4. [Docker Commands Reference](#docker-commands-reference)
5. [SSL Certificate Setup (Let's Encrypt)](#ssl-certificate-setup-lets-encrypt)
6. [Backup and Restore](#backup-and-restore)
7. [Monitoring and Logs](#monitoring-and-logs)
8. [Troubleshooting Common Issues](#troubleshooting-common-issues)
9. [Updating the Deployment](#updating-the-deployment)

---

## Server Requirements

### Hardware

| Resource | Minimum | Recommended |
|----------|---------|-------------|
| CPU | 2 cores | 4 cores |
| RAM | 2 GB | 4 GB+ |
| Disk | 20 GB SSD | 50 GB+ SSD |

> The judge service requires additional resources depending on concurrent judging load. Each judge container uses approximately 256 MB RAM. With the default pool of 4 containers, allocate at least 1 GB additional RAM for judging.

### Software

| Software | Version | Notes |
|----------|---------|-------|
| Operating System | Ubuntu 22.04 LTS / Debian 12 | Linux required for judge sandbox (namespaces, cgroups, seccomp) |
| Docker | 24.0+ | Include Docker Compose v2 plugin |
| Git | 2.40+ | For cloning the repository |
| Nginx | 1.24+ | Reverse proxy (install via `apt`) |
| Certbot | 2.0+ | For Let's Encrypt SSL certificates |

### Network

| Item | Value |
|------|-------|
| Domain | `oj.yogdunana.com` |
| DNS A Record | `101.237.129.33` |
| Nginx Listen Port | `23196` |
| Backend Port (internal) | `8080` |
| Frontend Port (internal) | `80` |
| MySQL Port (internal) | `3306` |
| gRPC Port (internal) | `50051` |

The host Nginx listens on port `23196` and proxies requests to the Docker frontend/backend services. A separate reverse proxy (e.g., Cloudflare or host-level Nginx) forwards traffic from `oj.yogdunana.com` to `101.237.129.33:23196`.

---

## First-Time Deployment

### Step 1: Prepare the Server

```bash
# Update system packages
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER
# Log out and back in for group change to take effect

# Install Docker Compose v2 (included with Docker Engine above)
docker compose version

# Install Nginx and Certbot
sudo apt install -y nginx certbot python3-certbot-nginx

# Install Git
sudo apt install -y git
```

### Step 2: Clone the Repository

```bash
git clone https://github.com/Yogdunana/yogduoj.git
cd yogduoj
```

### Step 3: Configure Environment

Copy and edit the production configuration:

```bash
# The backend reads config.yaml by default.
# Override sensitive values with environment variables:
export YOGDUOJ_DATABASE_PASSWORD="your-strong-mysql-password"
export YOGDUOJ_JWT_SECRET="your-strong-random-jwt-secret-at-least-32-chars"
```

Or create a production config file:

```bash
cp backend/config.yaml backend/config.prod.yaml
# Edit config.prod.yaml with production values
```

Then set the config path:

```bash
export YOGDUOJ_CONFIG_PATH="config.prod.yaml"
```

### Step 4: Configure Docker Compose

Edit `deploy/docker-compose.yml` to set environment variables for each service. The production compose file should include:

- **MySQL**: Set `MYSQL_ROOT_PASSWORD` and `MYSQL_DATABASE`.
- **Backend**: Pass `YOGDUOJ_DATABASE_HOST`, `YOGDUOJ_DATABASE_PASSWORD`, `YOGDUOJ_JWT_SECRET`, and `YOGDUOJ_SERVER_MODE=release`.
- **Frontend**: No special configuration needed (static files served by Nginx).
- **Judge**: Configure gRPC connection to the backend.

Example environment section for the backend service:

```yaml
services:
  backend:
    build:
      context: ../backend
      dockerfile: Dockerfile
    environment:
      - YOGDUOJ_DATABASE_HOST=mysql
      - YOGDUOJ_DATABASE_PORT=3306
      - YOGDUOJ_DATABASE_USER=root
      - YOGDUOJ_DATABASE_PASSWORD=${MYSQL_PASSWORD}
      - YOGDUOJ_DATABASE_DBNAME=yogduoj
      - YOGDUOJ_JWT_SECRET=${JWT_SECRET}
      - YOGDUOJ_SERVER_MODE=release
      - YOGDUOJ_LOG_LEVEL=info
    depends_on:
      - mysql
    ports:
      - "8080:8080"
    volumes:
      - backend-logs:/app/logs
      - judge-workspace:/app/workspace
```

### Step 5: Configure Host Nginx

Create an Nginx configuration for the reverse proxy:

```bash
sudo nano /etc/nginx/sites-available/yogduoj
```

```nginx
server {
    listen 23196;
    server_name oj.yogdunana.com;

    # Frontend (served by Docker Nginx)
    location / {
        proxy_pass http://127.0.0.1:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Backend API
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support (for judge status streaming)
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_read_timeout 86400;
    }

    # gRPC (judge service, internal only)
    # location /grpc/ {
    #     grpc_pass grpc://127.0.0.1:50051;
    # }
}
```

Enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/yogduoj /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Step 6: Start the Services

```bash
cd yogduoj

# Build and start all services
docker compose -f deploy/docker-compose.yml up -d --build

# Verify all containers are running
docker compose -f deploy/docker-compose.yml ps

# Check logs
docker compose -f deploy/docker-compose.yml logs -f
```

### Step 7: Verify the Deployment

```bash
# Health check
curl http://localhost:8080/api/v1/health
# Expected: {"service":"yogduoj-backend","status":"ok"}

# Frontend via Nginx
curl -I http://localhost:23196
# Expected: HTTP/1.1 200 OK

# External access (from your local machine)
curl https://oj.yogdunana.com/api/v1/health
```

---

## Environment Variables Reference

All environment variables use the `YOGDUOJ_` prefix. They override the corresponding values in `config.yaml`.

### Server

| Variable | Default | Description |
|----------|---------|-------------|
| `YOGDUOJ_SERVER_PORT` | `8080` | Backend HTTP port |
| `YOGDUOJ_SERVER_MODE` | `debug` | Gin mode: `debug`, `release`, or `test` |
| `YOGDUOJ_CONFIG_PATH` | `config.yaml` | Path to the config file |

### Database

| Variable | Default | Description |
|----------|---------|-------------|
| `YOGDUOJ_DATABASE_HOST` | `localhost` | MySQL host |
| `YOGDUOJ_DATABASE_PORT` | `3306` | MySQL port |
| `YOGDUOJ_DATABASE_USER` | `root` | MySQL username |
| `YOGDUOJ_DATABASE_PASSWORD` | (empty) | MySQL password |
| `YOGDUOJ_DATABASE_DBNAME` | `yogduoj` | Database name |
| `YOGDUOJ_DATABASE_CHARSET` | `utf8mb4` | Character set |
| `YOGDUOJ_DATABASE_MAX_IDLE_CONNS` | `10` | Max idle connections |
| `YOGDUOJ_DATABASE_MAX_OPEN_CONNS` | `100` | Max open connections |
| `YOGDUOJ_DATABASE_CONN_MAX_LIFETIME` | `3600` | Connection max lifetime (seconds) |

### JWT

| Variable | Default | Description |
|----------|---------|-------------|
| `YOGDUOJ_JWT_SECRET` | `yogduoj-default-secret-change-in-production` | **MUST change in production** |
| `YOGDUOJ_JWT_ACCESS_TTL` | `2h` | Access token lifetime |
| `YOGDUOJ_JWT_REFRESH_TTL` | `168h` | Refresh token lifetime (7 days) |

### Rate Limiting

| Variable | Default | Description |
|----------|---------|-------------|
| `YOGDUOJ_RATE_LIMIT_REQUESTS_PER_MINUTE` | `60` | Max requests per minute per IP |
| `YOGDUOJ_RATE_LIMIT_BURST` | `100` | Burst allowance |

### Logging

| Variable | Default | Description |
|----------|---------|-------------|
| `YOGDUOJ_LOG_LEVEL` | `debug` | Log level: `debug`, `info`, `warn`, `error` |
| `YOGDUOJ_LOG_FILE_PATH` | `logs/app.log` | Log file path |
| `YOGDUOJ_LOG_MAX_SIZE` | `100` | Max log file size (MB) |
| `YOGDUOJ_LOG_MAX_BACKUPS` | `10` | Max number of old log files |
| `YOGDUOJ_LOG_MAX_AGE` | `30` | Max days to retain logs |
| `YOGDUOJ_LOG_COMPRESS` | `true` | Compress rotated logs |

### Docker Compose Variables

| Variable | Description |
|----------|-------------|
| `MYSQL_ROOT_PASSWORD` | Root password for MySQL container |
| `MYSQL_DATABASE` | Database name to create |
| `JWT_SECRET` | JWT signing secret (passed to backend) |

---

## Docker Commands Reference

### Service Management

```bash
# Start all services (production)
docker compose -f deploy/docker-compose.yml up -d

# Start with rebuild
docker compose -f deploy/docker-compose.yml up -d --build

# Stop all services
docker compose -f deploy/docker-compose.yml down

# Stop and remove volumes (WARNING: deletes database data)
docker compose -f deploy/docker-compose.yml down -v

# Restart a single service
docker compose -f deploy/docker-compose.yml restart backend

# View status of all services
docker compose -f deploy/docker-compose.yml ps
```

### Logs

```bash
# Tail logs from all services
docker compose -f deploy/docker-compose.yml logs -f

# Tail logs from a specific service
docker compose -f deploy/docker-compose.yml logs -f backend
docker compose -f deploy/docker-compose.yml logs -f mysql
docker compose -f deploy/docker-compose.yml logs -f frontend

# Last 100 lines
docker compose -f deploy/docker-compose.yml logs --tail=100 backend
```

### Images

```bash
# Build all images
docker compose -f deploy/docker-compose.yml build

# Build a specific image
docker compose -f deploy/docker-compose.yml build backend

# Rebuild without cache
docker compose -f deploy/docker-compose.yml build --no-cache backend
```

### Volumes

```bash
# List Docker volumes
docker volume ls

# Inspect a specific volume
docker volume inspect yogduoj_mysql-data

# Remove unused volumes
docker volume prune
```

### Resource Usage

```bash
# View resource usage of running containers
docker stats

# Disk usage by Docker
docker system df
```

### Cleanup

```bash
# Remove stopped containers, unused networks, dangling images, and build cache
docker system prune

# Remove everything (including unused volumes and images)
docker system prune -a --volumes
```

---

## SSL Certificate Setup (Let's Encrypt)

### Option A: Direct SSL on the Server

If the server is directly accessible on port 443:

```bash
# Obtain certificate
sudo certbot --nginx -d oj.yogdunana.com

# Certbot will automatically modify the Nginx config to:
# 1. Listen on port 443 with SSL
# 2. Redirect HTTP (port 80) to HTTPS
# 3. Set up auto-renewal via cron

# Test auto-renewal
sudo certbot renew --dry-run
```

### Option B: SSL via External Reverse Proxy (Cloudflare, etc.)

If SSL termination happens at an external reverse proxy (e.g., Cloudflare) before reaching your server:

1. Configure the external proxy to forward HTTPS traffic to `http://101.237.129.33:23196`.
2. The host Nginx on port 23196 serves HTTP only; SSL is handled upstream.
3. No Let's Encrypt certificate is needed on the server.

### Option C: SSL on Host Nginx with Custom Port

If you need SSL on port 23196 (non-standard):

```bash
# Obtain certificate using standalone mode (port 80 must be free temporarily)
sudo certbot certonly --standalone -d oj.yogdunana.com

# Then update the Nginx config:
sudo nano /etc/nginx/sites-available/yogduoj
```

```nginx
server {
    listen 23196 ssl;
    server_name oj.yogdunana.com;

    ssl_certificate /etc/letsencrypt/live/oj.yogdunana.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/oj.yogdunana.com/privkey.pem;

    # SSL hardening
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # ... rest of proxy config ...
}
```

### Auto-Renewal

Certbot installs a systemd timer or cron job for auto-renewal. Verify it:

```bash
sudo systemctl status certbot.timer
# or
sudo crontab -l | grep certbot
```

---

## Backup and Restore

### Database Backup

```bash
# Create a backup
docker compose -f deploy/docker-compose.yml exec mysql \
  mysqldump -u root -p${MYSQL_PASSWORD} yogduoj \
  > backup_$(date +%Y%m%d_%H%M%S).sql

# Compress the backup
gzip backup_$(date +%Y%m%d_%H%M%S).sql
```

### Database Restore

```bash
# Restore from a backup file
gunzip backup_20260403_120000.sql.gz
docker compose -f deploy/docker-compose.yml exec -T mysql \
  mysql -u root -p${MYSQL_PASSWORD} yogduoj \
  < backup_20260403_120000.sql
```

### Full Backup (Database + Volumes)

```bash
# Stop services to ensure consistency
docker compose -f deploy/docker-compose.yml down

# Backup Docker volumes
docker run --rm -v yogduoj_mysql-data:/data -v $(pwd):/backup \
  alpine tar czf /backup/mysql-data-backup.tar.gz -C /data .

docker run --rm -v yogduoj_judge-workspace:/data -v $(pwd):/backup \
  alpine tar czf /backup/judge-workspace-backup.tar.gz -C /data .

# Restart services
docker compose -f deploy/docker-compose.yml up -d
```

### Automated Backups with Cron

Create a backup script:

```bash
sudo nano /opt/yogduoj-backup.sh
```

```bash
#!/bin/bash
set -euo pipefail

BACKUP_DIR="/opt/yogduoj-backups"
mkdir -p "$BACKUP_DIR"
DATE=$(date +%Y%m%d_%H%M%S)
COMPOSE_FILE="/path/to/yogduoj/deploy/docker-compose.yml"

# Database backup
docker compose -f "$COMPOSE_FILE" exec -T mysql \
  mysqldump -u root -p"${MYSQL_PASSWORD}" yogduoj \
  | gzip > "$BACKUP_DIR/db_${DATE}.sql.gz"

# Keep only the last 30 days of backups
find "$BACKUP_DIR" -name "db_*.sql.gz" -mtime +30 -delete

echo "[$(date)] Backup completed: db_${DATE}.sql.gz"
```

```bash
chmod +x /opt/yogduoj-backup.sh

# Add to cron (daily at 3:00 AM)
echo "0 3 * * * /opt/yogduoj-backup.sh >> /var/log/yogduoj-backup.log 2>&1" | sudo crontab -
```

---

## Monitoring and Logs

### Viewing Logs

```bash
# All services
docker compose -f deploy/docker-compose.yml logs -f

# Specific service
docker compose -f deploy/docker-compose.yml logs -f backend

# Backend logs (inside container)
docker compose -f deploy/docker-compose.yml exec backend cat /app/logs/app.log
```

### Health Checks

```bash
# Backend API health
curl -s http://localhost:8080/api/v1/health | jq .

# MySQL connectivity
docker compose -f deploy/docker-compose.yml exec mysql \
  mysqladmin -u root -p${MYSQL_PASSWORD} ping

# Docker container health
docker inspect --format='{{.State.Health.Status}}' yogduoj-backend-1
```

### Resource Monitoring

```bash
# Real-time container stats
docker stats

# Disk usage
docker system df

# MySQL disk usage
docker compose -f deploy/docker-compose.yml exec mysql \
  mysql -u root -p${MYSQL_PASSWORD} -e "
    SELECT table_schema AS 'Database',
           ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS 'Size (MB)'
    FROM information_schema.tables
    WHERE table_schema = 'yogduoj'
    GROUP BY table_schema;
  "
```

### Log Rotation

Backend logs are automatically rotated by Zap with the settings in `config.yaml`:

- `max_size`: 100 MB per file
- `max_backups`: 10 rotated files
- `max_age`: 30 days retention
- `compress`: gzip enabled

Docker container logs are managed by the Docker logging driver. To limit Docker log size, add to `/etc/docker/daemon.json`:

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "50m",
    "max-file": "5"
  }
}
```

Then restart Docker:

```bash
sudo systemctl restart docker
```

---

## Troubleshooting Common Issues

### Backend fails to start

**Symptom:** Backend container exits immediately or keeps restarting.

```bash
# Check logs
docker compose -f deploy/docker-compose.yml logs backend

# Common causes:
# 1. Cannot connect to MySQL
#    -> Ensure MySQL container is healthy first
#    -> Check YOGDUOJ_DATABASE_HOST matches the MySQL service name in compose
# 2. Config file not found
#    -> Ensure config.yaml is copied in the Dockerfile
#    -> Or set YOGDUOJ_CONFIG_PATH
# 3. Port already in use
#    -> Check: sudo lsof -i :8080
```

### MySQL connection refused

**Symptom:** Backend logs show "failed to connect to database" or "connection refused".

```bash
# Check if MySQL is running
docker compose -f deploy/docker-compose.yml ps mysql

# Check MySQL logs
docker compose -f deploy/docker-compose.yml logs mysql

# Verify credentials
docker compose -f deploy/docker-compose.yml exec mysql \
  mysql -u root -p${MYSQL_PASSWORD} -e "SELECT 1"

# Common causes:
# 1. MySQL not fully started yet (wait for health check)
# 2. Wrong password in environment variables
# 3. Network issue between containers (check docker network)
```

### Frontend shows blank page or 502

**Symptom:** Browser shows blank page or Nginx 502 Bad Gateway.

```bash
# Check if frontend container is running
docker compose -f deploy/docker-compose.yml ps frontend

# Check frontend Nginx logs
docker compose -f deploy/docker-compose.yml logs frontend

# Check host Nginx
sudo nginx -t
sudo systemctl status nginx
sudo tail -f /var/log/nginx/error.log

# Common causes:
# 1. Frontend build failed (check build logs)
# 2. Host Nginx misconfigured (proxy_pass pointing to wrong port)
# 3. Docker port mapping conflict
```

### Judge service not working

**Symptom:** Submissions stuck in "Pending" or "Judging" status.

```bash
# Check judge service logs
docker compose -f deploy/docker-compose.yml logs judge

# Check if judge can reach backend via gRPC
docker compose -f deploy/docker-compose.yml exec judge \
  nc -zv backend 50051

# Common causes:
# 1. gRPC port not exposed or blocked
# 2. Judge sandbox not supported (requires Linux kernel features)
# 3. Insufficient permissions for namespace/cgroup operations
# 4. Judge workspace volume not mounted
```

### Permission denied errors

```bash
# Check container user
docker compose -f deploy/docker-compose.yml exec backend whoami

# Fix volume permissions
docker compose -f deploy/docker-compose.yml exec backend chmod -R 755 /app/logs
docker compose -f deploy/docker-compose.yml exec backend chmod -R 755 /app/workspace
```

### Out of disk space

```bash
# Check Docker disk usage
docker system df

# Clean up unused resources
docker system prune -a

# Check server disk usage
df -h
du -sh /var/lib/docker/*
```

### Database migration errors

```bash
# Check current migration state
docker compose -f deploy/docker-compose.yml exec backend \
  ls -la /app/

# If auto-migration fails, check the migration code in:
# backend/internal/migration/migrate.go

# For manual migration, connect to MySQL:
docker compose -f deploy/docker-compose.yml exec mysql \
  mysql -u root -p${MYSQL_PASSWORD} yogduoj
```

---

## Updating the Deployment

### Standard Update Procedure

```bash
cd yogduoj

# 1. Pull latest code
git pull origin main

# 2. Backup database before updating
docker compose -f deploy/docker-compose.yml exec mysql \
  mysqldump -u root -p${MYSQL_PASSWORD} yogduoj \
  | gzip > backup_pre_update_$(date +%Y%m%d_%H%M%S).sql.gz

# 3. Rebuild and restart with zero downtime (rolling update)
docker compose -f deploy/docker-compose.yml up -d --build

# 4. Verify all services are healthy
docker compose -f deploy/docker-compose.yml ps
curl http://localhost:8080/api/v1/health

# 5. Check for migration errors in logs
docker compose -f deploy/docker-compose.yml logs --tail=50 backend
```

### Rollback Procedure

If the update causes issues:

```bash
# 1. Rollback code
git checkout <previous-commit-hash>

# 2. Restore database
gunzip backup_pre_update_YYYYMMDD_HHMMSS.sql.gz
docker compose -f deploy/docker-compose.yml exec -T mysql \
  mysql -u root -p${MYSQL_PASSWORD} yogduoj \
  < backup_pre_update_YYYYMMDD_HHMMSS.sql

# 3. Rebuild and restart
docker compose -f deploy/docker-compose.yml up -d --build
```

### Forced Recreation

When you need to completely recreate containers (e.g., config changes that require it):

```bash
docker compose -f deploy/docker-compose.yml up -d --build --force-recreate
```

### Using the Makefile

```bash
# Deploy to production
make deploy

# Stop production
make deploy-stop
```
