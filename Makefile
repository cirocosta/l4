install:
	go install -v

fmt:
	go fmt
	cd ./lib && go fmt

image:
	docker build -t cirocosta/l4 .

test:
	cd ./lib && go test -v

.PHONY: install
