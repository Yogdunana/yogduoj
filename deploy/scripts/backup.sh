#!/bin/bash
set -euo pipefail

# YogduOJ Database Backup Script
# Keeps the last 7 backups in /opt/yogduoj/backups/

DEPLOY_DIR="/opt/yogduoj"
BACKUP_DIR="${DEPLOY_DIR}/backups"
CONTAINER_NAME="yogduoj-mysql"
DB_NAME="yogduoj"
KEEP_COUNT=7

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/yogduoj_${TIMESTAMP}.sql.gz"

echo "========================================="
echo "  YogduOJ - Database Backup"
echo "========================================="
echo ""

# Ensure backup directory exists
mkdir -p "${BACKUP_DIR}"

# Run mysqldump inside the MySQL container and compress
echo "[1/3] Dumping database '${DB_NAME}' from container '${CONTAINER_NAME}'..."
docker exec "${CONTAINER_NAME}" \
    mysqldump \
    -u root \
    -p"${MYSQL_ROOT_PASSWORD:-}" \
    --single-transaction \
    --routines \
    --triggers \
    --add-drop-database \
    "${DB_NAME}" 2>/dev/null | gzip > "${BACKUP_FILE}"

if [ ! -s "${BACKUP_FILE}" ]; then
    echo "ERROR: Backup file is empty or mysqldump failed."
    rm -f "${BACKUP_FILE}"
    exit 1
fi

BACKUP_SIZE=$(du -h "${BACKUP_FILE}" | cut -f1)
echo "  -> Backup saved: ${BACKUP_FILE} (${BACKUP_SIZE})"

# Rotate old backups - keep last KEEP_COUNT
echo "[2/3] Rotating old backups (keeping last ${KEEP_COUNT})..."
BACKUP_COUNT=$(ls -1 "${BACKUP_DIR}"/yogduoj_*.sql.gz 2>/dev/null | wc -l)
if [ "${BACKUP_COUNT}" -gt "${KEEP_COUNT}" ]; then
    ls -1t "${BACKUP_DIR}"/yogduoj_*.sql.gz | tail -n +$((KEEP_COUNT + 1)) | xargs rm -f
    echo "  -> Removed $((BACKUP_COUNT - KEEP_COUNT)) old backup(s)."
else
    echo "  -> No rotation needed (${BACKUP_COUNT}/${KEEP_COUNT})."
fi

# Summary
echo "[3/3] Backup summary:"
echo "  -> Current backups:"
ls -lh "${BACKUP_DIR}"/yogduoj_*.sql.gz 2>/dev/null | awk '{print "     " $NF " (" $5 ")"}'
echo ""
echo "Backup complete!"
