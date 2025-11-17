build:
	mkdir -p functions
	go get ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o functions/foosbot main.go
