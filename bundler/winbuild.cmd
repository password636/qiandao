echo %PATH%
SET Path=%PATH%;%GOROOT%\bin\;
set GOPATH=%CD%
echo %GOPATH%
go get github.com/rakyll/statik
go get github.com/rakyll/statik/fs
bin\statik.exe -f -src=../build/ -dest=./src
go build -o build/qiandao.exe main-noui.go

