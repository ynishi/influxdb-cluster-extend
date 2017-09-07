.PHONY: build clean

build: clean
	cd cmd && go test
	cd cmd && go build

clean:
	rm -f "$(basename $(pwd))"
