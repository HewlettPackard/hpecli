@echo off
for /f "usebackq tokens=*" %%a in (`git rev-parse --short HEAD`) do set GIT_COMMIT=%%a

SET ISO_DATE=%DATE:~10,4%-%DATE:~4,2%-%DATE:~7,2%

@echo on
go build -o hpe.exe -ldflags "-X main.semanticVer=0.0.0 -X main.buildDate=%ISO_DATE% -X main.gitCommitID=%GIT_COMMIT%" .\cmd\hpecli


