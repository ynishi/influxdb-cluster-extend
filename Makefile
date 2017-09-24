.PHONY: build clean

build: clean
	go test
	cd cmd && go test
	cd cmd && go build

clean:
	rm -f "$(basename $(pwd))"
