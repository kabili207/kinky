@ECHO OFF

echo Compiling...
echo Windows:
set GOOS=windows
echo   x86
set GOARCH=386
go build -i -tags release -o build/kinky_win32.exe ./cmd/kinky
echo   x64
set GOARCH=amd64
go build -i -tags release -o build/kinky_win64.exe ./cmd/kinky

echo Linux:
set GOOS=linux
echo   x86
set GOARCH=386
go build -i -tags release -o build/kinky_l32 ./cmd/kinky
echo   x64
set GOARCH=amd64
go build -i -tags release -o build/kinky_l64 ./cmd/kinky
echo   arm
set GOARCH=arm
go build -i -tags release -o build/kinky_arm ./cmd/kinky
echo   arm64
set GOARCH=arm64
go build -i -tags release -o build/kinky_arm64 ./cmd/kinky
