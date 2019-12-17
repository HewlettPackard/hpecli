@echo off
for /f "usebackq tokens=*" %%a in (`git rev-parse --short HEAD`) do set GIT_COMMIT=%%a

SET ISO_DATE=%DATE:~10,4%-%DATE:~4,2%-%DATE:~7,2%
SET RELEASE_PKG=github.com/HewlettPackard/hpecli/pkg/version

@echo on
go build -o hpecli.exe -ldflags "-X %RELEASE_PKG%.version=v0.0.1 -X %RELEASE_PKG%.builtAt=%ISO_DATE% -X %RELEASE_PKG%.gitCommit=%GIT_COMMIT%" .\cmd\hpecli


