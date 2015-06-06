build:
	go get github.com/tools/godep
	go get github.com/jteeuwen/go-bindata/...
	godep restore
	bower install
	go-bindata -debug assets/...
	godep go build

release:
	go get github.com/tools/godep
	go get github.com/jteeuwen/go-bindata/...
	godep restore
	bower install
	go-bindata assets/...
	godep go build -tags=release
