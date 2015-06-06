build:
	bower install
	go-bindata -debug assets/...
	godep go build

release:
	bower install
	go-bindata assets/...
	godep go build -tags=release
