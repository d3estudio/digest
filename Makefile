.PHONY: build

VERSION := $(shell cat VERSION)

build: get_deps build_collector build_emoji_manager build_prefetcher build_processor build_tar build_sha

get_deps:
	go get -u github.com/jteeuwen/go-bindata/...
	go get github.com/tools/godep
	godep restore

build_collector:
	mkdir -p release
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o release/linux/amd64/collector 	github.com/d3estudio/digest/collector
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -o release/linux/arm64/collector 	github.com/d3estudio/digest/collector
	GOOS=linux   GOARCH=arm   CGO_ENABLED=0 go build -o release/linux/arm/collector 	github.com/d3estudio/digest/collector
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o release/windows/amd64/collector github.com/d3estudio/digest/collector
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -o release/darwin/amd64/collector 	github.com/d3estudio/digest/collector

build_prefetcher:
	mkdir -p release
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o release/linux/amd64/prefetcher 		github.com/d3estudio/digest/prefetcher
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -o release/linux/arm64/prefetcher 		github.com/d3estudio/digest/prefetcher
	GOOS=linux   GOARCH=arm   CGO_ENABLED=0 go build -o release/linux/arm/prefetcher 		github.com/d3estudio/digest/prefetcher
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o release/windows/amd64/prefetcher 	github.com/d3estudio/digest/prefetcher
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -o release/darwin/amd64/prefetcher 	github.com/d3estudio/digest/prefetcher

build_processor:
	mkdir -p release
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o release/linux/amd64/processor 		github.com/d3estudio/digest/processor
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -o release/linux/arm64/processor 		github.com/d3estudio/digest/processor
	GOOS=linux   GOARCH=arm   CGO_ENABLED=0 go build -o release/linux/arm/processor 		github.com/d3estudio/digest/processor
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o release/windows/amd64/processor 	github.com/d3estudio/digest/processor
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -o release/darwin/amd64/processor 		github.com/d3estudio/digest/processor

build_emoji_manager:
	mkdir -p release
	cd emoji-manager && ruby data/gen.rb && go-bindata -ignore .*\.rb data
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o release/linux/amd64/emoji-manager 		github.com/d3estudio/digest/emoji-manager
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -o release/linux/arm64/emoji-manager 		github.com/d3estudio/digest/emoji-manager
	GOOS=linux   GOARCH=arm   CGO_ENABLED=0 go build -o release/linux/arm/emoji-manager 		github.com/d3estudio/digest/emoji-manager
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o release/windows/amd64/emoji-manager 	github.com/d3estudio/digest/emoji-manager
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -o release/darwin/amd64/emoji-manager 		github.com/d3estudio/digest/emoji-manager

build_tar:
	mkdir -p release
	tar -cvzf release/linux/amd64/collector.tar.gz   -C release/linux/amd64   collector
	tar -cvzf release/linux/arm64/collector.tar.gz   -C release/linux/arm64   collector
	tar -cvzf release/linux/arm/collector.tar.gz     -C release/linux/arm     collector
	tar -cvzf release/windows/amd64/collector.tar.gz -C release/windows/amd64 collector
	tar -cvzf release/darwin/amd64/collector.tar.gz  -C release/darwin/amd64  collector

	tar -cvzf release/linux/amd64/prefetcher.tar.gz   -C release/linux/amd64   prefetcher
	tar -cvzf release/linux/arm64/prefetcher.tar.gz   -C release/linux/arm64   prefetcher
	tar -cvzf release/linux/arm/prefetcher.tar.gz     -C release/linux/arm     prefetcher
	tar -cvzf release/windows/amd64/prefetcher.tar.gz -C release/windows/amd64 prefetcher
	tar -cvzf release/darwin/amd64/prefetcher.tar.gz  -C release/darwin/amd64  prefetcher

	tar -cvzf release/linux/amd64/processor.tar.gz   -C release/linux/amd64   processor
	tar -cvzf release/linux/arm64/processor.tar.gz   -C release/linux/arm64   processor
	tar -cvzf release/linux/arm/processor.tar.gz     -C release/linux/arm     processor
	tar -cvzf release/windows/amd64/processor.tar.gz -C release/windows/amd64 processor
	tar -cvzf release/darwin/amd64/processor.tar.gz  -C release/darwin/amd64  processor

	tar -cvzf release/linux/amd64/emoji-manager.tar.gz   -C release/linux/amd64   emoji-manager
	tar -cvzf release/linux/arm64/emoji-manager.tar.gz   -C release/linux/arm64   emoji-manager
	tar -cvzf release/linux/arm/emoji-manager.tar.gz     -C release/linux/arm     emoji-manager
	tar -cvzf release/windows/amd64/emoji-manager.tar.gz -C release/windows/amd64 emoji-manager
	tar -cvzf release/darwin/amd64/emoji-manager.tar.gz  -C release/darwin/amd64  emoji-manager

build_sha:
	mkdir -p release
	shasum -a 256 release/linux/amd64/collector.tar.gz   > release/linux/amd64/collector.sha256
	shasum -a 256 release/linux/arm64/collector.tar.gz   > release/linux/arm64/collector.sha256
	shasum -a 256 release/linux/arm/collector.tar.gz     > release/linux/arm/collector.sha256
	shasum -a 256 release/windows/amd64/collector.tar.gz > release/windows/amd64/collector.sha256
	shasum -a 256 release/darwin/amd64/collector.tar.gz  > release/darwin/amd64/collector.sha256

	shasum -a 256 release/linux/amd64/prefetcher.tar.gz   > release/linux/amd64/prefetcher.sha256
	shasum -a 256 release/linux/arm64/prefetcher.tar.gz   > release/linux/arm64/prefetcher.sha256
	shasum -a 256 release/linux/arm/prefetcher.tar.gz     > release/linux/arm/prefetcher.sha256
	shasum -a 256 release/windows/amd64/prefetcher.tar.gz > release/windows/amd64/prefetcher.sha256
	shasum -a 256 release/darwin/amd64/prefetcher.tar.gz  > release/darwin/amd64/prefetcher.sha256

	shasum -a 256 release/linux/amd64/processor.tar.gz   > release/linux/amd64/processor.sha256
	shasum -a 256 release/linux/arm64/processor.tar.gz   > release/linux/arm64/processor.sha256
	shasum -a 256 release/linux/arm/processor.tar.gz     > release/linux/arm/processor.sha256
	shasum -a 256 release/windows/amd64/processor.tar.gz > release/windows/amd64/processor.sha256
	shasum -a 256 release/darwin/amd64/processor.tar.gz  > release/darwin/amd64/processor.sha256

	shasum -a 256 release/linux/amd64/emoji-manager.tar.gz   > release/linux/amd64/emoji-manager.sha256
	shasum -a 256 release/linux/arm64/emoji-manager.tar.gz   > release/linux/arm64/emoji-manager.sha256
	shasum -a 256 release/linux/arm/emoji-manager.tar.gz     > release/linux/arm/emoji-manager.sha256
	shasum -a 256 release/windows/amd64/emoji-manager.tar.gz > release/windows/amd64/emoji-manager.sha256
	shasum -a 256 release/darwin/amd64/emoji-manager.tar.gz  > release/darwin/amd64/emoji-manager.sha256

build_docker:
	docker build --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
               --build-arg VCS_REF=`git rev-parse --short HEAD` \
               --build-arg VERSION=`cat VERSION` \
							 -t "d3estudio/digest:$(VERSION)" \
							 .
	docker push "d3estudio/digest:$(VERSION)"
