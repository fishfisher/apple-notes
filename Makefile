.PHONY: build install test clean

build:
	go build -o apple-notes .

install:
	go install .

test:
	go test -v ./...

clean:
	rm -f apple-notes
	rm -rf dist/

release:
	goreleaser release --clean

snapshot:
	goreleaser release --snapshot --clean
