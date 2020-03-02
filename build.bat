@echo off
for /f "usebackq tokens=*" %%a in (`git rev-parse --short HEAD`) do set GIT_COMMIT=%%a

SET ISO_DATE=%DATE:~10,4%-%DATE:~4,2%-%DATE:~7,2%
SET VERSION_FILE=github.com/HewlettPackard/hpecli/pkg/version

@echo on
go build -o hpecli.exe -ldflags "-X %VERSION_FILE%.version=0.0.0 -X %VERSION_FILE%.buildDate=%ISO_DATE% -X %VERSION_FILE%.gitCommitId=%GIT_COMMIT%" .\cmd\hpecli


