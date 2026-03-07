#!/bin/bash

# Load config file
if [ -f .deploy.conf ]; then
    source .deploy.conf
else
    echo "Error: .deploy.conf not found"
    echo "Create it from .deploy.conf.example"
    exit 1
fi


echo -e "DEPLOYING Tiumo to ${PI_IP}"

echo -e "COPYING the binary to Raspberry [${PI_USER}@${PI_IP}:${TMP_DIR}]"
scp -q ./build/web ${PI_USER}@${PI_IP}:${TMP_DIR}

ssh -t ${PI_USER}@${PI_IP} << EOF
echo -e "BACKING UP the current version to [${TARGET_DIR}/backup]"
sudo cp ${TARGET_DIR}/web ${TARGET_DIR}/backup/web-backup-$(date +%Y%m%d-%H%M%S)

echo -e "INSTALLING the new version to [${TARGET_DIR}]"
sudo mv ${TMP_DIR}/web ${TARGET_DIR}

echo -e "RESTARTING the service"
sudo systemctl restart tiimit
sudo systemctl status tiimit
EOF