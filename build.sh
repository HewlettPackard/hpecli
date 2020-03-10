#!/bin/bash
ISO_DATE=`date +%F`
GIT_COMMIT=`git rev-parse --short HEAD`

go build -o hpecli -ldflags "-X main.sematicVer=0.0.1 -X main.buildDate=$ISO_DATE -X main.gitCommitID=$GIT_COMMIT" ./cmd/hpecli
