build:
	go-bindata -debug assets/...
	go build

release:
	go-bindata assets/...
	go build -tags=release
