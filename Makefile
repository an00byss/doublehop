Binary_Name=doublehop.exe

build:
	env GOOS=windows GOARCH=amd64 go build -ldflags "-w -s" -o ${BINARY_NAME} .
clean:
	go clean
