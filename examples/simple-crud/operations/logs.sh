#!/bin/bash
# Description: Deploy the app to the VPS


SSH_KEY_PATH=~/.ssh/id_rsa
TARGET_HOST=vps

# exit when any command fails
set -e

echo "===> Checking the status of the service"
ssh -i $SSH_KEY_PATH $TARGET_HOST 'sudo systemctl status gourmet'
ssh -i $SSH_KEY_PATH $TARGET_HOST 'journalctl -u gourmet -f'





