#!/bin/bash
# Description: Deploy the app to the VPS


SSH_KEY_PATH=~/.ssh/id_rsa
TARGET_HOST=vps

# exit when any command fails
set -e

# echo "===> Updating remote server dependencies"
# ssh -i $SSH_KEY_PATH $TARGET_HOST 'sudo apt-get update && sudo apt-get upgrade -y && sudo apt-get autoremove -y'

echo "===> Zipping the binary"
zip gourmetapp.zip gourmet-app

echo "===> Copying the binary into server with temporary location so the downtime is minimal"
scp -i $SSH_KEY_PATH gourmetapp.zip $TARGET_HOST:/home/ubuntu/gourmet-tmp.zip & scp -i $SSH_KEY_PATH ./operations/gourmet.service $TARGET_HOST:/tmp/gourmet.service


echo "===> Unzipping the binary"
ssh -i $SSH_KEY_PATH $TARGET_HOST 'unzip -o /home/ubuntu/gourmet-tmp.zip -d /home/ubuntu/'

# echo "===> Move the migration files"
# scp -i $SSH_KEY_PATH -r db/ $TARGET_HOST:/home/ubuntu

echo "===> Moving the service file at the right place"
ssh -i $SSH_KEY_PATH $TARGET_HOST 'sudo mv /tmp/gourmet.service /etc/systemd/system/gourmet.service'

echo "===> Reloading the daemon & Stopping the service"
ssh -i $SSH_KEY_PATH $TARGET_HOST 'sudo systemctl daemon-reload && sudo systemctl stop gourmet'

echo "===> Moving the binary at the right place (overwriting the old one quickly)"
ssh -i $SSH_KEY_PATH $TARGET_HOST 'sudo mv /home/ubuntu/gourmet-app /home/ubuntu/gourmet/'

echo "===> Starting the service"
echo "If it fails, you can check the logs with: journalctl -u gourmet -f. Possible errors are:"
echo "- the binary is not executable, in that case, you can fix it with: sudo chmod +x /home/ubuntu/gourmet-back"
echo "- the .env file is not present"
ssh -i $SSH_KEY_PATH $TARGET_HOST 'sudo systemctl start gourmet'

echo "===> Cleaning up"
rm gourmet-app gourmetapp.zip
ssh -i $SSH_KEY_PATH $TARGET_HOST 'rm /home/ubuntu/gourmet-tmp.zip'

echo "===> Done"
