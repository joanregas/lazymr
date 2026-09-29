APP := lazymr

.PHONY: all build run test clean

all: build

build:
	go build -o $(APP) .

run:
	go run .

test:
	go test ./...

clean:
	rm -f $(APP)

