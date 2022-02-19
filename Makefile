foo:
	go run ./cmd/foo

bar:
	go run ./cmd/bar

baz:
	go run ./cmd/baz

call:
	curl http://localhost:8001/foo
