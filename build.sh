#!/bin/bash
DT=`date +%F`
GC=`git rev-parse --short HEAD`
VERSION_FILE=github.com/HewlettPackard/hpecli/pkg/version

go build -o hpecli -ldflags "-X '$VERSION_FILE.version=0.0.1' -X '$VERSION_FILE.buildDate=$DT' -X '$VERSION_FILE.gitCommitId=$GC'" ./cmd/hpecli
