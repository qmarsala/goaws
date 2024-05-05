$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -o bootstrap ./cmd/calculate-handicap/main.go
~\Go\Bin\build-lambda-zip.exe -o calculate-handicap.zip bootstrap

go build -o bootstrap ./cmd/enter-round/main.go
~\Go\Bin\build-lambda-zip.exe -o enter-round.zip bootstrap