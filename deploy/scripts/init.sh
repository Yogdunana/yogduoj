#!/bin/bash
set -euo pipefail

DEPLOY_DIR="/opt/yogduoj"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "========================================="
echo "  YogduOJ - First-time Deployment Setup"
echo "========================================="
echo ""

# Create deployment directory
echo "[1/5] Creating deployment directory: ${DEPLOY_DIR}"
sudo mkdir -p "${DEPLOY_DIR}"
sudo chown "$(id -u):$(id -g)" "${DEPLOY_DIR}"

# Copy .env.example to .env if it does not exist
echo "[2/5] Setting up environment configuration..."
if [ ! -f "${DEPLOY_DIR}/.env" ]; then
    cp "${PROJECT_DIR}/deploy/.env.example" "${DEPLOY_DIR}/.env"
    echo "  -> Created .env from .env.example"
    echo "  -> IMPORTANT: Edit ${DEPLOY_DIR}/.env and update passwords/secrets before continuing!"
else
    echo "  -> .env already exists, skipping."
fi

# Copy docker-compose files
echo "[3/5] Copying Docker Compose configuration..."
cp "${PROJECT_DIR}/deploy/docker-compose.yml" "${DEPLOY_DIR}/docker-compose.yml"
cp -r "${PROJECT_DIR}/deploy/nginx" "${DEPLOY_DIR}/nginx"
cp -r "${PROJECT_DIR}/deploy/mysql" "${DEPLOY_DIR}/mysql"
echo "  -> Copied docker-compose.yml, nginx config, and MySQL init script."

# Create Docker volumes
echo "[4/5] Ensuring Docker volumes exist..."
docker volume create yogduoj_mysql_data 2>/dev/null || true
docker volume create yogduoj_submission_data 2>/dev/null || true
docker volume create yogduoj_problem_data 2>/dev/null || true
docker volume create yogduoj_cert_data 2>/dev/null || true
echo "  -> Volumes ready."

# Start services
echo "[5/5] Starting YogduOJ services..."
cd "${DEPLOY_DIR}"
docker compose up -d

# Wait for MySQL to become healthy
echo ""
echo "Waiting for MySQL to become healthy..."
timeout=120
elapsed=0
until docker compose exec -T mysql mysqladmin ping -h localhost -u root --silent 2>/dev/null; do
    if [ "$elapsed" -ge "$timeout" ]; then
        echo "ERROR: MySQL did not become healthy within ${timeout}s."
        exit 1
    fi
    sleep 5
    elapsed=$((elapsed + 5))
    echo "  ... waiting (${elapsed}s/${timeout}s)"
done
echo "  -> MySQL is healthy!"

# Show status
echo ""
echo "========================================="
echo "  YogduOJ Services Status"
echo "========================================="
docker compose ps
echo ""
echo "Deployment complete!"
echo "  - Frontend:  http://localhost:23196"
echo "  - Backend:   http://localhost:8080 (internal)"
echo "  - Judge:     gRPC on localhost:50051 (internal)"
echo "  - MySQL:     localhost:3306 (internal)"
echo ""
echo "Next steps:"
echo "  1. Edit ${DEPLOY_DIR}/.env with your secrets"
echo "  2. Place SSL certificates in the cert_data volume"
echo "  3. Restart: cd ${DEPLOY_DIR} && docker compose restart"
