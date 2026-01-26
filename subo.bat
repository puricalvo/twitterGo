@echo off
set GOOS=linux
set GOARCH=amd64

go build -o bootstrap main.go

del main.zip 2>nul
tar -a -cf main.zip bootstrap

git add .
git commit -m "Deploy Lambda"
git push