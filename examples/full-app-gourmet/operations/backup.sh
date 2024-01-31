#!/bin/bash
# Description: Deploy the app to the VPS


SSH_KEY_PATH=~/.ssh/id_rsa
TARGET_HOST=vps

# exit when any command fails
set -e

# echo "===> Updating remote server dependencies"
# ssh -i $SSH_KEY_PATH $TARGET_HOST 'sudo apt-get update && sudo apt-get upgrade -y && sudo apt-get autoremove -y'

echo "===> Copy the SQLite database"
scp -i $SSH_KEY_PATH $TARGET_HOST:/home/ubuntu/gourmet.db gourmet.bak.db
