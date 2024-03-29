#!/bin/sh

GOARCH=amd64 GOOS=linux go build main.go
rsync -avz -e "ssh -o StrictHostKeyChecking=no" --progress main dkg@core.local:~/medusa/medusa
rsync -avz -e "ssh -o StrictHostKeyChecking=no" --progress ../core/core.cfg.json dkg@core.local:~/medusa/
