.PHONY: release

default: release

release:
	GOOS=darwin GOARCH=amd64 go build -o release/badger cmd/main.go; cd release; zip badger-macOS.zip badger
	GOOS=linux GOARCH=amd64 go build -o release/badger cmd/main.go; cd release; zip badger-linux.zip badger
	GOOS=windows GOARCH=amd64 go build -o release/badger cmd/main.go; cd release; zip badger-windows.zip badger
