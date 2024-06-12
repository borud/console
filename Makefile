all: test lint vet build

build: console

console:
	@cd cmd/$@ && go build -o ../../bin/$@

test:
	@go test ./...

race:
	@go test -race ./...

vet:
	@go vet ./...

lint:
	@revive ./...

clean:
	@rm -rf bin
