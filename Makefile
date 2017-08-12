VERSION := $(shell cat ./VERSION)

install:
	go install -v

fmt:
	go fmt
	cd ./lib && go fmt

image:
	docker build -t cirocosta/l4 .

test:
	cd ./lib && go test -v

release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)
	goreleaser --rm-dist


.PHONY: install fmt image test release
