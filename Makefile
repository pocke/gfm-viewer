build:
	bower install
	go-bindata -debug assets/...
	go build

release:
	bower install
	go-bindata assets/...
	go build -tags=release
