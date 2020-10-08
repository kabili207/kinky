@ECHO OFF
go generate ./...
go build -i -tags dev -o build/kinky.exe ./cmd/kinky
build\kinky.exe %*
