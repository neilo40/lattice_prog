all:
	GOOS=linux GOARCH=amd64 go build -o bin/prog -o bin/prog/prog-amd64 cmd/prog/main.go
	GOOS=linux GOARCH=arm GOARM=7 go build -o bin/prog -o bin/prog/prog-armv7l cmd/prog/main.go