#!/bin/bash
set -euo pipefail

BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/afritechonline_$TIMESTAMP.sql"

mkdir -p "$BACKUP_DIR"

echo "Creating database backup..."
docker compose exec -T db mysqldump -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" > "$BACKUP_FILE"

echo "Backup created: $BACKUP_FILE"
echo "Compressing..."
gzip "$BACKUP_FILE"

echo "Backup complete: ${BACKUP_FILE}.gz"
