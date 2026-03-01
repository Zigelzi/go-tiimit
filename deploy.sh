#!/bin/bash

# Load config file
if [ -f .deploy.conf ]; then
    source .deploy.conf
else
    echo "Error: .deploy.conf not found"
    echo "Create it from .deploy.conf.example"
    exit 1
fi

echo -e "DEPLOYING Tiumo to ${PI_IP}\n"

echo -e "COPYING the binary to Raspberry [${PI_USER}@${PI_IP}:${TMP_DIR}]"
scp ./build/web ${PI_USER}@${PI_IP}:${TMP_DIR}

ssh ${PI_USER}@${PI_IP} << EOF
echo -e "BACKING UP the current version to [${TARGET_DIR}/backup/]\n"
sudo cp ${TMP_DIR}/web ${TARGET_DIR}/backup/web-backup-$(date +%Y%m%d-%H%M%S)

echo -e "INSTALLING the new version to [${TARGET_DIR}]\n"
sudo mv ${TMP_DIR}/web ${TARGET_DIR}

echo -e "RESTARTING the service\n"
sudo systemctl restart tiimit.service
sudo systemctl status tiimit
EOF