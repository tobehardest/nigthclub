set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
set etc=.\output\etc

go build -o ./output/nightclub.exe -ldflags "-s -w" ./nightclub/nightclub.go

if exist  %etc% (
   rmdir /s/q %etc%
)

mkdir %etc%

copy .\nightclub\etc .\output\etc
