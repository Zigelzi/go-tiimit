#!/bin/bash
set -euo pipefail

# Load config file
if [ -f .deploy.conf ]; then
    source .deploy.conf
else
    echo "Error: .deploy.conf not found"
    echo "Create it from .deploy.conf.example"
    exit 1
fi

echo "RUNNING tests"
go test ./...

echo "BUILDING Tiumo"
make prod/build-arm64
echo "DEPLOYING Tiumo to ${PI_IP}"

echo "COPYING the binary to Raspberry [${PI_USER}@${PI_IP}:${TMP_FILE}]"
scp -q ./build/web ${PI_USER}@${PI_IP}:${TMP_FILE}

ssh -t ${PI_USER}@${PI_IP} "
set -e
echo 'BACKING UP the current version to [${TARGET_DIR}/backup]'
sudo cp ${TARGET_DIR}/web ${TARGET_DIR}/backup/web-backup-$(date +%Y%m%d-%H%M%S)

echo 'INSTALLING the new version to [${TARGET_DIR}]'
sudo mv ${TMP_FILE} ${TARGET_DIR}/web

echo 'RESTARTING the service'
sudo systemctl restart tiimit
systemctl status tiimit
"