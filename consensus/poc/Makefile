ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	cd ./shabal && cargo build --release
	#cp ./shabal/target/release/libshabal.dylib ./
	cp ./shabal/target/release/libshabal.so ./
	go build -ldflags="-r $(ROOT_DIR)" shabal.go

run: build
	./shabal
