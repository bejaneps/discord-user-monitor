SET GOPATH=%cd%
SET GO111MODULE=on

mkdir bin

go build -o bin\scrapper.exe cmd\scrapper\main.go
