#!/bin/sh

GOARCH=arm GOOS=linux go build main.go
rsync -avz -e "ssh -o StrictHostKeyChecking=no" --progress main medusa.local:~/medusa/medusa
rsync -avz -e "ssh -o StrictHostKeyChecking=no" --progress ../core/core.cfg.json medusa.local:~/medusa/
