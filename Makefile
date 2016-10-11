test:
	go test ./...

clean:
	go clean
	rm -rfv ./pkg/*

fmt:
	go fmt ./...
