#!/bin/sh

GOARCH=amd64 GOOS=linux go build main.go
rsync -avz -e "ssh -o StrictHostKeyChecking=no" --progress main core.local:~/medusa/medusa
rsync -avz -e "ssh -o StrictHostKeyChecking=no" --progress ../core/core.cfg.json core.local:~/medusa/
