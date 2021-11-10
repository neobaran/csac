build:
	go build -ldflags="-s -w" -ldflags="-X 'main.Version=$(version)'" -o csac main.go
	$(if $(shell command -v upx), upx csac)

check:
	golangci-lint run