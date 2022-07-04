test:
	go test -cover ./...

armor:
	go build -o bin/armor ./cmd/armor