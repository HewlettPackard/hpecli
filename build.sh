#!/bin/bash
DT=`date +%F`
GC=`git rev-parse --short HEAD`
RELEASE_PKG=github.com/HewlettPackard/hpecli/pkg/version

go build -o hpecli -ldflags "-X '$RELEASE_PKG.version=0.0.1' -X '$RELEASE_PKG.builtAt=$DT' -X '$RELEASE_PKG.gitCommit=$GC'" ./cmd/hpecli