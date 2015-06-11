build:
	go-bindata -debug assets/...
	godep go build

release:
	go-bindata assets/...
	godep go build -tags=release

depends:
	go get github.com/tools/godep
	go get github.com/jteeuwen/go-bindata/...
	godep restore
	bower install

install:
	godep go install -tags=release

test:
	godep go test -v -race -tags=release
