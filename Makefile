test:
	go test -cover ./...

armor:
	go build -o bin/armor ./cmd/armor

npm:
	cd frontend && npm start

local: armor npm