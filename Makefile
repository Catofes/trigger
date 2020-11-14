build:
	mkdir -p bin
	CGO_ENABLED=0 go build -trimpath -o bin/web ./web
	CGO_ENABLED=0 go build -trimpath -o bin/trigger ./client
clean:
	rm -r bin/
.PHONY: build