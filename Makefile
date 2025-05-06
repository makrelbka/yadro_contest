build:
	go mod tidy
	go build -o generated/biathlon cmd/biathlon/main.go

run: build
	./generated/biathlon

test: build
	go test -v ./integration-test -count=1