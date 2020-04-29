.PHONY: build clean deploy gomodgen

build: 
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello github.com/gotha/splunk-lambda-example/hello
clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

