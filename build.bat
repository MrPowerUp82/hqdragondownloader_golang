go build -ldflags -H=windowsgui -o hqdragondownloader.exe main.go

upx -9 -k hqdragondownloader.exe

certutil -hashfile .\hqdragondownloader.exe MD5
certutil -hashfile .\hqdragondownloader.exe SHA256