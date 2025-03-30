.PHONY: test test-cover clean

test:
	go test ./... -v

test-cover:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -f coverage.out coverage.html 