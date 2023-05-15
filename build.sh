# /bin/bash

GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o hqdragondownloader-win64.exe main.go
# $env:GOOS='windows'; $env:GOARCH=386; go build -ldflags -H=windowsgui -o hqdragondownloader-win32.exe main.go

upx -9 -k hqdragondownloader-win64.exe
# upx -9 -k hqdragondownloader-win32.exe

md5sum ./hqdragondownloader-win64.exe
sha256sum -hashfile ./hqdragondownloader-win64.exe

# certutil -hashfile .\hqdragondownloader-win32.exe MD5
# certutil -hashfile .\hqdragondownloader-win32.exe SHA256