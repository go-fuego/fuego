#!/bin/bash
# Description: Upload the SQLite database to the server

SSH_KEY_PATH=~/.ssh/id_rsa
TARGET_HOST=vps

# exit when any command fails
set -e

echo "===> Copy the SQLite database"
scp -i $SSH_KEY_PATH recipe.db $TARGET_HOST:/home/ubuntu/gourmet.db
