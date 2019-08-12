default: clean

test:
	go test

check:
	go build
	go test -cover

clean:
	GO111MODULE=on go get -u
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor